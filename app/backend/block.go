package backend

import (
	"context"

	"babel/openapi/gen/babelapi"

	"go.uber.org/zap"
)

func (s *mServer) TranslateBlock(ctx context.Context, request babelapi.TranslateBlockRequestObject) (babelapi.TranslateBlockResponseObject, error) {
	tids := request.Body.TranslationIds
	if len(tids) == 0 {
		return &babelapi.TranslateBlock200JSONResponse{
			Blocks: make([]babelapi.Block, 0),
		}, nil
	}

	bid := request.BlockId
	blocks, err := s.options.storageBlock.Translate(ctx, bid, tids)
	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}

	var bs []babelapi.Block
	for _, b := range blocks {
		bs = append(bs, *b)
	}

	return &babelapi.TranslateBlock200JSONResponse{
		Blocks: bs,
	}, nil
}
