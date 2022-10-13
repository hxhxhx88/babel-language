package postgres

import (
	"fmt"

	"github.com/cenkalti/backoff"
	_ "github.com/jackc/pgx/v5/stdlib" // pgx driver
	"github.com/jmoiron/sqlx"
	"github.com/labstack/gommon/log"
	"github.com/pkg/errors"
)

type Auth struct {
	Username string
	Password string
	Database string
	Host     string
	Port     int
}

func Connect(auth Auth) (*sqlx.DB, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		auth.Username,
		auth.Password,
		auth.Host,
		auth.Port,
		auth.Database,
	)
	pg, err := sqlx.Open("pgx", connStr)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	ping := func() error {
		err := pg.Ping()
		if err != nil {
			log.Info("postgres is not ready...backing off...")
			return err
		}
		log.Info("postgres is ready!")
		return nil
	}

	err = backoff.Retry(ping, backoff.NewExponentialBackOff())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return pg, nil
}
