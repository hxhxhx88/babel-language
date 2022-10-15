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

func (s *mServer) CreateCorpusTranslation(ctx context.Context, request babelapi.CreateCorpusTranslationRequestObject) (babelapi.CreateCorpusTranslationResponseObject, error) {
	id, err := s.options.storageCorpus.CreateTranslation(ctx, request.CorpusId, &request.Body.Translation)
	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}
	return &babelapi.CreateCorpusTranslation200JSONResponse{
		Id: id,
	}, nil
}

func (s *mServer) ListCorpuses(ctx context.Context, request babelapi.ListCorpusesRequestObject) (babelapi.ListCorpusesResponseObject, error) {
	corpuses, err := s.options.storageCorpus.List(ctx)
	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}
	return &babelapi.ListCorpuses200JSONResponse{
		Corpuses: corpuses,
	}, nil
}
