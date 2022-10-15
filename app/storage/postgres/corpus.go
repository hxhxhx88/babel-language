package postgres

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	iso639_3 "github.com/barbashov/iso639-3"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"babel/app/storage"
	"babel/openapi/gen/babelapi"
)

func NewCorpus(pg *sqlx.DB) storage.Corpus {
	return &mCorpus{
		pg: pg,
	}
}

type mCorpus struct {
	pg *sqlx.DB
}

func (s *mCorpus) Create(ctx context.Context, corpus *babelapi.CorpusDraft) (storage.IdType, error) {
	tx, err := s.pg.BeginTxx(ctx, nil)
	if err != nil {
		return "", errors.WithStack(err)
	}

	id, err := txCreate(ctx, tx, corpus)
	if err != nil {
		tx.Rollback()
		return "", err
	}

	if err = tx.Commit(); err != nil {
		return "", err
	}

	return strconv.Itoa(id), nil
}

func (s *mCorpus) CreateTranslation(ctx context.Context, corpusId storage.IdType, translation *babelapi.TranslationDraft) (storage.IdType, error) {
	cid, err := strconv.Atoi(corpusId)
	if err != nil {
		return "", errors.WithStack(err)
	}

	tx, err := s.pg.BeginTxx(ctx, nil)
	if err != nil {
		return "", errors.WithStack(err)
	}

	ids, err := txCreateTranslation(ctx, tx, cid, *translation)
	if err != nil {
		tx.Rollback()
		return "", err
	}

	if err = tx.Commit(); err != nil {
		return "", err
	}

	return strconv.Itoa(ids[0]), nil
}

func (s *mCorpus) List(ctx context.Context) ([]babelapi.Corpus, error) {
	var records []struct {
		Id                      int    `db:"id"`
		Title                   string `db:"title"`
		OriginalLanguageIso6393 string `db:"original_language_iso_639_3"`
	}
	if err := s.pg.Select(&records, "SELECT id, title, original_language_iso_639_3 FROM corpuses ORDER BY title ASC"); err != nil {
		return nil, errors.WithStack(err)
	}

	var corpuses []babelapi.Corpus
	for _, rec := range records {
		corpuses = append(corpuses, babelapi.Corpus{
			Id:                      strconv.Itoa(rec.Id),
			Title:                   rec.Title,
			OriginalLanguageIso6393: rec.OriginalLanguageIso6393,
		})
	}

	return corpuses, nil
}

const mParamLimit = 65535

func txCreate(ctx context.Context, tx *sqlx.Tx, corpus *babelapi.CorpusDraft) (int, error) {
	if iso639_3.FromPart3Code(corpus.OriginalLanguageIso6393) == nil {
		return 0, errors.Errorf("invalid ISO-639-3 code [%s] for corpus", corpus.OriginalLanguageIso6393)
	}

	var corpusId int
	{
		stmt, err := tx.PrepareNamed(`
			INSERT INTO corpuses
				(title, original_language_iso_639_3)
			VALUES
				(:title, :original_language_iso_639_3)
			RETURNING id
		`)
		if err != nil {
			return 0, errors.WithStack(err)
		}

		if err := stmt.GetContext(ctx, &corpusId, map[string]any{
			"title":                       corpus.Title,
			"original_language_iso_639_3": corpus.OriginalLanguageIso6393,
		}); err != nil {
			return 0, errors.WithStack(err)
		}
	}

	if corpus.Translations == nil || len(*corpus.Translations) == 0 {
		return corpusId, nil
	}

	if _, err := txCreateTranslation(ctx, tx, corpusId, *corpus.Translations...); err != nil {
		return 0, err
	}

	return corpusId, nil
}

func txCreateTranslation(ctx context.Context, tx *sqlx.Tx, corpusId int, translations ...babelapi.TranslationDraft) ([]int, error) {
	for idx, t := range translations {
		if iso639_3.FromPart3Code(t.LanguageIso6393) == nil {
			return nil, errors.Errorf("invalid ISO-639-3 code [%s] for translation at index %d", t.LanguageIso6393, idx)
		}
	}

	var translationIds []int
	{
		var vs []any
		var ps []string
		for _, t := range translations {
			p := fmt.Sprintf("($%d, $%d, $%d)", len(vs)+1, len(vs)+2, len(vs)+3)
			ps = append(ps, p)
			vs = append(vs, corpusId, t.Title, t.LanguageIso6393)
		}
		if len(vs) > 0 {
			query := `INSERT INTO translations (corpus_id, title, language_iso_639_3) VALUES ` + strings.Join(ps, ",") + " RETURNING id "

			rows, err := tx.QueryxContext(ctx, query, vs...)
			if err != nil {
				return nil, errors.WithStack(err)
			}
			for rows.Next() {
				var tid int
				if err := rows.Scan(&tid); err != nil {
					return nil, errors.WithStack(err)
				}
				translationIds = append(translationIds, tid)
			}
		}
	}

	var recs []map[string]any
	for i, t := range translations {
		if t.Blocks == nil {
			continue
		}

		tid := translationIds[i]
		for _, b := range *t.Blocks {
			recs = append(recs, map[string]any{
				"translation_id": tid,
				"content":        b.Content,
				"rank":           b.Rank,
				"uuid":           b.Uuid,
			})
		}
	}

	if len(recs) == 0 {
		return translationIds, nil
	}

	numParamPerRec := len(recs[0])
	batchSize := mParamLimit / numParamPerRec
	numBatch := len(recs) / batchSize
	for i := 0; i < numBatch+1; i++ {
		var batch []map[string]any
		if i == numBatch {
			if len(recs)%batchSize == 0 {
				break
			}
			batch = recs[numBatch*batchSize:]
		} else {
			batch = recs[i*batchSize : (i+1)*batchSize]
		}

		if _, err := tx.NamedExecContext(ctx, `
			INSERT INTO blocks
				(translation_id, content, rank, uuid)
			VALUES
				(:translation_id, :content, :rank, :uuid)
		`, batch); err != nil {
			return nil, errors.WithStack(err)
		}
	}

	return translationIds, nil
}
