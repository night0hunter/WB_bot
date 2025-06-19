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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dbpool, err := db.NewPG(ctx, connString)
	if err != nil {
		log.Fatalf("db.NewPG: %s", err)
	}

	fmt.Printf("Database has been started on port %s\n", port)

	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))
	if err != nil {
		log.Fatalf("tgbotapi.NewBotAPI: %s", err)
	}

	fmt.Printf("Bot has been started\n")

	service := service.NewService(dbpool)

	h := handler.New(bot, service)
	// handler := handler.NewHandler(bot, service)

	err = h.Run(ctx)
	// err = handler.Run(ctx)
	if err != nil {
		fmt.Printf("handler.Run: %s", err.Error())
	}
}
