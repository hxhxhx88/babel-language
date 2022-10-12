package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"babel/backend"
	"babel/openapi"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())

	s := backend.New()
	openapi.RegisterHandlers(e, s)

	e.Logger.Fatal(e.Start(":12346"))
}
