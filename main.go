package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"github.com/willis7/golossary/slack"
)

func init() {
	viper.SetConfigName("app_config")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}

func main() {
	db := InitDb()

	token := viper.GetString("slack.token")
	client := slack.NewClient(token, db)
	client.Connect()
	defer client.Close()
	go client.Dispatch()

	sigterm := make(chan os.Signal)
	signal.Notify(sigterm, os.Interrupt)
	for {
		select {
		case <-sigterm:
			log.Println("terminate signal recvd")
			err := client.Shutdown()
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-time.After(time.Second):
			}
			client.Close()
			return
		}
	}
}

// InitDb
func InitDb() *sql.DB {
	dbInfo := fmt.Sprintf(
		"user=%s dbname=%s password=%s host=%s port=%s sslmode=disable",
		viper.GetString("database.username"),
		viper.GetString("database.name"),
		viper.GetString("database.password"),
		viper.GetString("database.hostname"),
		viper.GetString("database.port"),
	)
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	_, err = db.Exec(
		`CREATE TABLE if NOT EXISTS "words" (
			    "id" serial,
			    "word" varchar(56) NOT NULL UNIQUE,
			    "definition" varchar(500) NOT NULL,
			    CONSTRAINT words_pk PRIMARY KEY ("id")
		    )`)

	if err != nil {
		log.Fatal(err)
	}
	return db
}
