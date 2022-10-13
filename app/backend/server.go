package backend

import (
	"context"

	"babel/app/buildtime"
	"babel/openapi/gen/babelapi"
)

func New(opts ...Option) (babelapi.StrictServerInterface, error) {
	o := &Options{}
	for _, opt := range opts {
		opt(o)
	}
	if err := o.Validate(); err != nil {
		return nil, err
	}

	s := &mServer{
		options: o,
	}

	return s, nil
}

type mServer struct {
	options *Options
}

func (s *mServer) GetMetadata(ctx context.Context, request babelapi.GetMetadataRequestObject) (babelapi.GetMetadataResponseObject, error) {
	return &babelapi.GetMetadata200JSONResponse{
		CommitIdentifier: buildtime.CommitIdentifier,
		CommitTime:       buildtime.CommitTime,
		Version:          buildtime.Version,
	}, nil
}
