package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	sqldblogger "github.com/simukti/sqldb-logger"
	"github.com/simukti/sqldb-logger/logadapter/zerologadapter"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("failed to env file,err:%v", err)
	}
	appToken := os.Getenv("APP_TOKEN")
	botToken := os.Getenv("BOT_TOKEN")
	api := slack.New(botToken,
		slack.OptionDebug(false),
		slack.OptionLog(log.New(os.Stdout, "slack connect: ", log.Lshortfile|log.LstdFlags)),
		slack.OptionAppLevelToken(appToken),
	)
	client := socketmode.New(
		api,
		socketmode.OptionDebug(false),
		socketmode.OptionLog(log.New(os.Stdout, "socket connect: ", log.Lshortfile|log.LstdFlags)),
	)
	db := setupDB()
	tt := NewTimeTrackerCommand(client, db)
	go func(tt *TimeTrackerCommand) {
		for evt := range client.Events {
			switch evt.Type {
			case socketmode.EventTypeConnecting:
				fmt.Println("connecting to Slack with Socket Mode...")
			case socketmode.EventTypeConnectionError:
				fmt.Printf("connection failed. Retrying later...,err:%v", evt.Request.Reason)
			case socketmode.EventTypeConnected:
				fmt.Println("connected to Slack with Socket Mode.")
			case socketmode.EventTypeSlashCommand:
				cmd, ok := evt.Data.(slack.SlashCommand)
				if !ok {
					fmt.Printf("Ignored %v", evt)
					continue
				}
				tt.Do(evt, cmd)
			}
		}
	}(tt)
	client.Run()
}


func setupDB() *sqlx.DB {
	dsn := "host=postgres user=user password=password dbname=slack sslmode=disable"
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatalf("fail to connect database,err:%v", err)
	}
	logger := zerolog.New(
		zerolog.ConsoleWriter{Out: os.Stdout, NoColor: false},
	)
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
	return db
}
