package postgres

import (
	"github.com/jackc/pgconn"
)

func IsErrorUniqueViolation(err error) *pgconn.PgError {
	if err, ok := err.(*pgconn.PgError); ok {
		// 23505 is the PostgreSQl error code for unique violation
		// https://www.postgresql.org/docs/14/errcodes-appendix.html
		if err.Code == "23505" {
			return err
		}
	}
	return nil
}
