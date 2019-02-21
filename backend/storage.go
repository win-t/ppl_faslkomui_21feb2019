package main

import (
	"context"
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	migration "github.com/payfazz/psql-migration"
)

type storage struct {
	errLog *log.Logger
	db     *sql.DB
}

func getStorage(errLog *log.Logger) *storage {
	db, err := sql.Open("postgres", getConf("DATABASE_URI"))
	if err != nil {
		errLog.Panicln(err)
	}

	if err := migration.Migrate(context.Background(), migration.MigrateParams{
		Database:      db,
		ErrorLog:      errLog,
		ApplicationID: applicationID,
		Statements:    migrationStatement,
	}); err != nil {
		errLog.Panicln(err)
	}

	return &storage{
		errLog: errLog,
		db:     db,
	}
}

func (s *storage) Close() error {
	return s.db.Close()
}

func (s *storage) resetCounter() error {
	if _, err := s.db.Exec(
		`update keyvalue set value=0 where key='counter';`,
	); err != nil {
		s.errLog.Println(err)
		return err
	}
	return nil
}

func (s *storage) getCounter() (int, error) {
	var res int

	if err := s.db.QueryRow(
		`select value from keyvalue where key='counter';`,
	).Scan(&res); err != nil {
		s.errLog.Println(err)
		return 0, err
	}

	return res, nil
}

func (s *storage) incCounter() error {
	tx, err := s.db.Begin()
	if err != nil {
		s.errLog.Println(err)
		return err
	}
	commited := false
	defer func() {
		if !commited {
			tx.Rollback()
		}
	}()

	var val int
	if err := tx.QueryRow(
		`select value from keyvalue where key='counter';`,
	).Scan(&val); err != nil {
		s.errLog.Println(err)
		return err
	}

	val++
	// val = val + 10

	if _, err := tx.Exec(
		`update keyvalue set value=$1 where key='counter';`,
		val,
	); err != nil {
		s.errLog.Println(err)
		return err
	}

	if err := tx.Commit(); err != nil {
		s.errLog.Println(err)
		return err
	}

	commited = true
	return nil
}
