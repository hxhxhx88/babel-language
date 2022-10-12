package backend

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"babel/openapi"
)

func New() openapi.ServerInterface {
	return &mServer{}
}

type mServer struct {
}

func (s *mServer) Ping(ctx echo.Context) error {
	return ctx.NoContent(http.StatusOK)
}
