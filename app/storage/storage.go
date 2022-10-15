package storage

import (
	"context"

	"babel/openapi/gen/babelapi"
)

type IdType = string

type Corpus interface {
	Create(ctx context.Context, corpus *babelapi.CorpusDraft) (IdType, error)
	CreateTranslation(ctx context.Context, corpusId IdType, tranlation *babelapi.TranslationDraft) (IdType, error)
	List(ctx context.Context) ([]*babelapi.Corpus, error)
	Get(ctx context.Context, corpusId IdType) (*babelapi.Corpus, error)
}
