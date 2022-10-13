package storage

import (
	"context"

	"babel/openapi/gen/babelapi"
)

type IdType = string

type Corpus interface {
	Create(context.Context, *babelapi.CorpusDraft) (IdType, error)
}
