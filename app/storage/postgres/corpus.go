package postgres

import (
	"context"
	"fmt"
	"strconv"
	"strings"

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

func txCreate(ctx context.Context, tx *sqlx.Tx, corpus *babelapi.CorpusDraft) (int, error) {
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

	var translationIds []int
	{
		var vs []any
		var ps []string
		for _, t := range *corpus.Translations {
			p := fmt.Sprintf("($%d, $%d)", len(vs)+1, len(vs)+2)
			ps = append(ps, p)
			vs = append(vs, corpusId, t.LanguageIso6393)
		}
		if len(vs) > 0 {
			query := `INSERT INTO translations (corpus_id, language_iso_639_3) VALUES ` + strings.Join(ps, ",")
			query += " RETURNING id "

			rows, err := tx.QueryxContext(ctx, query, vs...)
			if err != nil {
				return 0, errors.WithStack(err)
			}
			for rows.Next() {
				var tid int
				if err := rows.Scan(&tid); err != nil {
					return 0, errors.WithStack(err)
				}
				translationIds = append(translationIds, tid)
			}
		}
	}

	if len(translationIds) == 0 {
		return corpusId, nil
	}

	{
		var recs []map[string]any
		for i, t := range *corpus.Translations {
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

		if len(recs) > 0 {
			if _, err := tx.NamedExecContext(ctx, `
				INSERT INTO blocks
					(translation_id, content, rank, uuid)
				VALUES
					(:translation_id, :content, :rank, :uuid)
			`, recs); err != nil {
				return 0, errors.WithStack(err)
			}
		}
	}

	return corpusId, nil
}
