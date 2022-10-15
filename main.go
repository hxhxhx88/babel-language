package main

import (
	"embed"
	"fmt"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"babel/app/backend"
	"babel/app/storage/postgres"
	"babel/openapi/gen/babelapi"
)

//go:embed app/frontend/build/*
var frontend embed.FS
var frontendFS = echo.MustSubFS(frontend, "app/frontend/build")

var c struct {
	Port int

	Postgres postgres.Auth
}

func main() {
	mustParseConfig()

	// logger
	mustSetupLogger()

	// server
	e := echo.New()
	e.Use(middleware.Logger())

	// backend
	s, teardown := mustCreateServer()
	defer teardown()
	babelapi.RegisterHandlers(e.Group("/api"), s)

	// frontend
	e.StaticFS("/app", frontendFS)
	e.StaticFS("/", frontendFS)

	// start
	lisAddr := fmt.Sprintf(":%d", c.Port)
	e.Logger.Fatal(e.Start(lisAddr))
}

func mustParseConfig() {
	withFlagSetFork := func(root *pflag.FlagSet, prefix string, do func(*pflag.FlagSet)) {
		sub := pflag.NewFlagSet(prefix, pflag.ContinueOnError)
		sub.SetNormalizeFunc(func(f *pflag.FlagSet, name string) pflag.NormalizedName {
			// In order to correctly map to struct, nested fields should be separated by dots.
			return pflag.NormalizedName(prefix + "." + name)
		})
		do(sub)
		root.AddFlagSet(sub)
	}

	pflag.IntVarP(&c.Port, "port", "p", 12346, "port to listen")
	withFlagSetFork(pflag.CommandLine, "postgres", func(fs *pflag.FlagSet) {
		fs.StringVar(&c.Postgres.Username, "username", "", "postgres username")
		fs.StringVar(&c.Postgres.Password, "password", "", "postgres password")
		fs.StringVar(&c.Postgres.Host, "host", "", "postgres host")
		fs.IntVar(&c.Postgres.Port, "port", 0, "postgres port")
		fs.StringVar(&c.Postgres.Database, "database", "", "postgres database")
	})
	pflag.Parse()

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
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

func mustCreateServer() (babelapi.ServerInterface, func()) {
	var opts []backend.Option

	// postgres
	pg, err := postgres.Connect(c.Postgres)
	if err != nil {
		panic(err)
	}

	opts = append(opts,
		backend.WithCorpusStorage(postgres.NewCorpus(pg)),
		backend.WithTranslationStorage(postgres.NewTranslation(pg)),
	)

	// backedn
	s, err := backend.New(opts...)
	if err != nil {
		panic(err)
	}

	return babelapi.NewStrictHandler(s, nil), func() { pg.Close() }
}

func mustSetupLogger() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(logger)
}
