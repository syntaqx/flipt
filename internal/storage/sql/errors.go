package sql

import (
	"database/sql"
	"errors"

	"github.com/go-sql-driver/mysql"
	"github.com/lib/pq"
	"github.com/mattn/go-sqlite3"
	errs "go.flipt.io/flipt/errors"
)

var (
	errNotFound           = errs.ErrNotFound("resource")
	errConstraintViolated = errs.ErrInvalid("contraint violated")
	errNotUnique          = errs.ErrInvalid("not unique")
	errForeignKeyNotFound = errs.ErrNotFound("associated resource not found")
	errCanceled           = errs.ErrCanceled("query canceled")
)

// AdaptError converts specific known-driver errors into wrapped storage errors.
func (d Driver) AdaptError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, sql.ErrNoRows) {
		return errNotFound
	}

	switch d {
	case SQLite:
		return adaptSQLiteError(err)
	case CockroachDB, Postgres:
		return adaptPostgresError(err)
	case MySQL:
		return adaptMySQLError(err)
	}

	return err
}

func adaptSQLiteError(err error) error {
	var serr sqlite3.Error

	if errors.As(err, &serr) {
		if serr.Code == sqlite3.ErrConstraint {
			switch serr.ExtendedCode {
			case sqlite3.ErrConstraintForeignKey:
				return errForeignKeyNotFound
			case sqlite3.ErrConstraintUnique:
				return errNotUnique
			}

			return errConstraintViolated
		}
	}

	return err
}

func adaptPostgresError(err error) error {
	const (
		constraintForeignKeyErr = "foreign_key_violation"
		constraintUniqueErr     = "unique_violation"
		queryCanceled           = "query_canceled"
	)

	var perr *pq.Error

	if errors.As(err, &perr) {
		switch perr.Code.Name() {
		case constraintUniqueErr:
			return errNotUnique
		case constraintForeignKeyErr:
			return errForeignKeyNotFound
		case queryCanceled:
			return errCanceled
		}
	}

	return err
}

func adaptMySQLError(err error) error {
	const (
		constraintForeignKeyErrCode uint16 = 1452
		constraintUniqueErrCode     uint16 = 1062
	)

	var merr *mysql.MySQLError

	if errors.As(err, &merr) {
		switch merr.Number {
		case constraintForeignKeyErrCode:
			return errForeignKeyNotFound
		case constraintUniqueErrCode:
			return errNotUnique
		}
	}

	return err
}
