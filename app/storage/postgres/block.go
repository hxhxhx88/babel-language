package postgres

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"babel/app/storage"
	"babel/openapi/gen/babelapi"
)

func NewBlock(pg *sqlx.DB) storage.Block {
	return &mBlock{
		pg: pg,
	}
}

type mBlock struct {
	pg *sqlx.DB
}

func (s *mBlock) Translate(ctx context.Context, bid storage.BlockId, tids []storage.TranslationId) ([]*babelapi.Block, error) {
	query, args, err := sqlx.In(`
		SELECT
			r.id,
			r.translation_id,
			r.rank,
			r.uuid,
			r.content
		FROM blocks AS b
		JOIN translations AS bt ON b.translation_id = bt.id
		JOIN blocks AS r ON r.uuid = b.uuid
		JOIN translations AS rt ON r.translation_id = rt.id
		WHERE b.id = ? AND bt.corpus_id = rt.corpus_id AND rt.id IN (?)
	`, bid, tids)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	query = sqlx.Rebind(sqlx.DOLLAR, query)

	var records []mBlockRecord
	if err := s.pg.SelectContext(ctx, &records, query, args...); err != nil {
		return nil, errors.WithStack(err)
	}

	var blocks []*babelapi.Block
	for _, rec := range records {
		blocks = append(blocks, rec.ToOpenApi())
	}

	return blocks, nil
}
