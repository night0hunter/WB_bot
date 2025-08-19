package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"wb_bot/db"
	wbadp "wb_bot/internal/adapter/wb-adp"
	"wb_bot/internal/handler"
	"wb_bot/internal/service"
	logger "wb_bot/pkg/log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := godotenv.Load(); err != nil {
		logger.FatalKV(ctx, "godotenv.Load", zap.Error(err))
	}

	logger.SetLevel(getLogLevel())

	var (
		host     = os.Getenv("HOST")
		port     = os.Getenv("PORT")
		user     = os.Getenv("POSTGRES_USER")
		password = os.Getenv("POSTGRES_PASSWORD")
		dbname   = os.Getenv("POSTGRES_DB")
	)

	var connString = fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s",
		user,
		password,
		host,
		port,
		dbname,
	)

	dbpool, err := db.NewPG(ctx, connString)
	if err != nil {
		logger.FatalKV(ctx, "db.NewPG", zap.Error(err))
	}

	fmt.Printf("Database has been started on port %s\n", port)

	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))
	if err != nil {
		logger.FatalKV(ctx, "tgbotapi.NewBotAPI", zap.Error(err))
	}

	fmt.Printf("Bot has been started\n")

	service := service.New(dbpool, &wbadp.Adapter{})

	h := handler.New(bot, service)
	err = h.Run(ctx)
	if err != nil {
		logger.FatalKV(ctx, "handler.Run", zap.Error(err))
	}
}

func getLogLevel() zapcore.Level {
	v := os.Getenv("LOG_LEVEL")
	if v == "" {
		return zap.WarnLevel
	}

	parsed, _ := strconv.Atoi(v)

	return zapcore.Level(parsed)
}
