package storage

import (
	"context"

	"babel/openapi/gen/babelapi"
)

type idType = string
type CorpusId = idType
type TranslationId = idType
type BlockId = idType

type Corpus interface {
	Create(context.Context, *babelapi.CorpusDraft) (CorpusId, error)
	List(context.Context) ([]*babelapi.Corpus, error)
	Get(context.Context, CorpusId) (*babelapi.Corpus, error)

	CreateTranslation(context.Context, CorpusId, *babelapi.TranslationDraft) (TranslationId, error)
	ListTranslations(context.Context, CorpusId) ([]*babelapi.Translation, error)
}

type Translation interface {
	SearchBlocks(context.Context, TranslationId, *babelapi.BlockFilter, *babelapi.Pagination) ([]*babelapi.Block, error)
	CountBlocks(context.Context, TranslationId, *babelapi.BlockFilter) (uint64, error)
}

type Block interface {
	Translate(context.Context, BlockId, []TranslationId) ([]*babelapi.Block, error)
}
