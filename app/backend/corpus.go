package backend

import (
	"context"

	"babel/openapi/gen/babelapi"

	"go.uber.org/zap"
)

func (s *mServer) CreateCorpus(ctx context.Context, request babelapi.CreateCorpusRequestObject) (babelapi.CreateCorpusResponseObject, error) {
	id, err := s.options.storageCorpus.Create(ctx, &request.Body.Corpus)
	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}
	return &babelapi.CreateCorpus200JSONResponse{
		Id: id,
	}, nil
}
