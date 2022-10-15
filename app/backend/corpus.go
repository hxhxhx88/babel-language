package backend

import (
	"context"

	"go.uber.org/zap"

	"babel/openapi/gen/babelapi"
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

func (s *mServer) ListCorpuses(ctx context.Context, request babelapi.ListCorpusesRequestObject) (babelapi.ListCorpusesResponseObject, error) {
	corpuses, err := s.options.storageCorpus.List(ctx)
	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}
	var cs []babelapi.Corpus
	for _, c := range corpuses {
		cs = append(cs, *c)
	}
	return &babelapi.ListCorpuses200JSONResponse{
		Corpuses: cs,
	}, nil
}

func (s *mServer) GetCorpus(ctx context.Context, request babelapi.GetCorpusRequestObject) (babelapi.GetCorpusResponseObject, error) {
	corpus, err := s.options.storageCorpus.Get(ctx, request.CorpusId)
	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}
	return &babelapi.GetCorpus200JSONResponse{
		Corpus: *corpus,
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

func (s *mServer) ListCorpusTranslations(ctx context.Context, request babelapi.ListCorpusTranslationsRequestObject) (babelapi.ListCorpusTranslationsResponseObject, error) {
	translations, err := s.options.storageCorpus.ListTranslations(ctx, request.CorpusId)
	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}
	var ts []babelapi.Translation
	for _, t := range translations {
		ts = append(ts, *t)
	}
	return &babelapi.ListCorpusTranslations200JSONResponse{
		Translations: ts,
	}, nil
}
