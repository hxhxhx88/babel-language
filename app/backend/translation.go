package backend

import (
	"context"

	"babel/openapi/gen/babelapi"

	"go.uber.org/zap"
)

func (s *mServer) ListTranslationBlocks(ctx context.Context, request babelapi.ListTranslationBlocksRequestObject) (babelapi.ListTranslationBlocksResponseObject, error) {
	tid := request.TranslationId
	f := request.Body.Filter
	p := normalizePagination(request.Body.Pagination)

	blocks, err := s.options.storageTranslation.ListBlocks(ctx, tid, f, p)
	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}

	var bs []babelapi.Block
	for _, b := range blocks {
		bs = append(bs, *b)
	}

	return &babelapi.ListTranslationBlocks200JSONResponse{
		Blocks: bs,
	}, nil
}

func (s *mServer) CountTranslationBlocks(ctx context.Context, request babelapi.CountTranslationBlocksRequestObject) (babelapi.CountTranslationBlocksResponseObject, error) {
	tid := request.TranslationId
	f := request.Body.Filter

	n, err := s.options.storageTranslation.CountBlocks(ctx, tid, f)
	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}

	return &babelapi.CountTranslationBlocks200JSONResponse{
		TotalCount: int(n),
	}, nil
}
