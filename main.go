package main

import (
	"embed"
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"babel/app/backend"
	"babel/openapi/gen/babelapi"
)

//go:embed app/frontend/build/*
var frontend embed.FS
var frontendFS = echo.MustSubFS(frontend, "app/frontend/build")

var c struct {
	Port int
}

func main() {
	mustParseConfig()

	e := echo.New()
	e.Use(middleware.Logger())

	// backend
	s := backend.New()
	h := babelapi.NewStrictHandler(s, nil)
	babelapi.RegisterHandlers(e.Group("/api"), h)

	// frontend
	e.StaticFS("/view", frontendFS)
	e.StaticFS("/", frontendFS)

	// start
	lisAddr := fmt.Sprintf(":%d", c.Port)
	e.Logger.Fatal(e.Start(lisAddr))
}

func mustParseConfig() {
	pflag.IntVarP(&c.Port, "port", "p", 12346, "port to listen")
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	viper.SetEnvPrefix("babel")
	viper.AutomaticEnv()

	viper.SetConfigName("conf.yaml")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			panic(err)
		}
	}
	if err := viper.Unmarshal(&c); err != nil {
		panic(err)
	}
}
