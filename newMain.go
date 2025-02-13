package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"wb_bot/db"
	"wb_bot/internal/handler"
	"wb_bot/internal/service"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("godotenv.Load: %s", err)
	}

	var (
		host     = os.Getenv("HOST")
		port     = os.Getenv("PORT")
		user     = os.Getenv("USER")
		password = os.Getenv("PASSWORD")
		dbname   = os.Getenv("DBNAME")
	)

	var connString = fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s",
		user,
		password,
		host,
		port,
		dbname,
	)

	dbpool, err := db.NewPG(context.Background(), connString)
	if err != nil {
		log.Fatalf("db.NewPG: %s", err)
	}

	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))
	if err != nil {
		log.Fatalf("tgbotapi.NewBotAPI: %s", err)
	}

	service := service.NewService(dbpool)

	handler := handler.NewHandler(bot, service)

	err = handler.Run(context.Background())
	if err != nil {
		fmt.Printf("handler.Run: %s", err.Error())
	}
}
