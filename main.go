package main

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"babel/backend"
	"babel/openapi"
)

var c struct {
	Port int
}

func main() {
	mustParseConfig()

	e := echo.New()
	e.Use(middleware.Logger())

	s := backend.New()
	openapi.RegisterHandlers(e, s)

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
