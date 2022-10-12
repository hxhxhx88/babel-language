package backend

import (
	"context"

	"babel/app/buildtime"
	"babel/openapi/gen/babelapi"
)

func New() babelapi.StrictServerInterface {
	return &mServer{}
}

type mServer struct {
}

func (s *mServer) GetMetadata(ctx context.Context, request babelapi.GetMetadataRequestObject) (babelapi.GetMetadataResponseObject, error) {
	return &babelapi.GetMetadata200JSONResponse{
		CommitIdentifier: buildtime.CommitIdentifier,
		CommitTime:       buildtime.CommitTime,
		Version:          buildtime.Version,
	}, nil
}
