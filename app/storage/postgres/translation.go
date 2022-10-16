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

func NewTranslation(pg *sqlx.DB) storage.Translation {
	return &mTranslation{
		pg: pg,
	}
}

type mTranslation struct {
	pg *sqlx.DB
}

func (s *mTranslation) SearchBlocks(ctx context.Context, tid storage.TranslationId, f *babelapi.BlockFilter, p *babelapi.Pagination) ([]*babelapi.Block, error) {
	tid_, err := strconv.Atoi(tid)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	limit := p.PageSize
	offset := p.PageSize * p.Page

	filterClause, filterArgs := makeBlockFilterStatement(tid_, f)
	query := fmt.Sprintf(`
		SELECT
			b.id,
			b.translation_id,
			b.rank,
			b.uuid,
			b.content
		FROM blocks AS b
		%s
		LIMIT %d
		OFFSET %d
	`, filterClause, limit, offset)

	stmt, err := s.pg.PrepareNamedContext(ctx, query)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var records []mBlockRecord
	if err := stmt.SelectContext(ctx, &records, filterArgs); err != nil {
		return nil, errors.WithStack(err)
	}

	var blocks []*babelapi.Block
	for _, rec := range records {
		blocks = append(blocks, rec.ToOpenApi())
	}

	return blocks, nil
}

func (s *mTranslation) CountBlocks(ctx context.Context, tid storage.TranslationId, f *babelapi.BlockFilter) (uint64, error) {
	tid_, err := strconv.Atoi(tid)
	if err != nil {
		return 0, errors.WithStack(err)
	}

	filterClause, filterArgs := makeBlockFilterStatement(tid_, f)
	query := fmt.Sprintf(`SELECT COUNT(*) FROM blocks AS b %s`, filterClause)

	stmt, err := s.pg.PrepareNamedContext(ctx, query)
	if err != nil {
		return 0, errors.WithStack(err)
	}

	var n uint64
	if err := stmt.GetContext(ctx, &n, filterArgs); err != nil {
		return 0, errors.WithStack(err)
	}

	return n, nil
}

func makeBlockFilterStatement(translationId int, f *babelapi.BlockFilter) (string, map[string]interface{}) {
	var conds []string
	var joins []string
	args := make(map[string]interface{})

	conds = append(conds, "b.translation_id = :translation_id")
	args["translation_id"] = translationId

	if pid := f.ParentBlockId; pid != nil {
		joins = append(joins, `JOIN blocks AS pb ON b.uuid LIKE pb.uuid || '%'`)
		conds = append(conds,
			"pb.id = :parent_block_id",
			"b.rank = pb.rank + 1",
		)
		args["parent_block_id"] = *pid
	} else {
		conds = append(conds, "b.rank = 1")
	}

	q := strings.Join(joins, " ")
	q += ` WHERE `
	q += strings.Join(conds, " AND ")

	return q, args
}
