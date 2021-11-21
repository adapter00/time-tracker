package main

import (
	"log"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
	sqldblogger "github.com/simukti/sqldb-logger"
	"github.com/simukti/sqldb-logger/logadapter/zerologadapter"
)

func setupDBTest() *sqlx.DB {
	dsn := "host=127.0.0.1 user=user password=password dbname=slack_test sslmode=disable"
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatalf("fail to connect database,err:%v", err)
	}
	logger := zerolog.New(
		zerolog.ConsoleWriter{Out: os.Stdout, NoColor: false},
	)
	// populate log pre-fields here before set to
	db.DB = sqldblogger.OpenDriver(
		dsn,
		db.Driver(),
		zerologadapter.New(logger),
		sqldblogger.WithSQLQueryAsMessage(true),
		sqldblogger.WithSQLArgsFieldname("sql_args"),
		sqldblogger.WithQueryerLevel(sqldblogger.LevelDebug),
		sqldblogger.WithPreparerLevel(sqldblogger.LevelDebug),
		sqldblogger.WithExecerLevel(sqldblogger.LevelDebug),
	)

	db.MustExec(schema)
	return db
}

var testDB *sqlx.DB

func TestMain(m *testing.M) {
	testDB = setupDBTest()
	defer func() {
		testDB.Exec("truncate table tracks;")
	}()
	os.Exit(m.Run())
}
