package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
	"wb_bot/db"
	wbadp "wb_bot/internal/adapter/wb-adp"
	cronjob "wb_bot/internal/cronJob"
	"wb_bot/internal/handler"
	"wb_bot/internal/middleware"
	"wb_bot/internal/service"
	"wb_bot/internal/utils"
	logger "wb_bot/pkg/log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/net/publicsuffix"

	"github.com/joho/godotenv"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := godotenv.Load(); err != nil {
		logger.FatalKV(ctx, "godotenv.Load", zap.Error(err))
	}

	logger.SetLevel(getLogLevel())

	var bearer = os.Getenv("BEARER_TOKEN")

	jar, err := cookiejar.New(&cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	})
	if err != nil {
		logger.FatalKV(ctx, "cookiejar.New", zap.Error(err))
	}

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

	fmt.Printf("Base has been started on port %s\n", port)

	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))
	if err != nil {
		logger.FatalKV(ctx, "tgbotapi.NewBotAPI", zap.Error(err))
	}

	fmt.Printf("Bot has been started\n")

	httpClient := &http.Client{
		Timeout: 10 * time.Second,
		Transport: middleware.Chain(
			nil,
			middleware.SetHeader(middleware.AuthorizationHeader, "Bearer "+bearer),
		),
		Jar: jar,
	}
	// mockClient := &mock.MockClient{}

	wbSupplyAdp := wbadp.New(
		httpClient,
		"https://supplies-api.wildberries.ru",
	)

	service := service.New(dbpool, wbSupplyAdp)
	handler := handler.New(bot, service)
	trackingCron := cronjob.NewSendTrackingCron(handler)

	c := cron.New(
		cron.WithLocation(utils.MoscowLocation),
		cron.WithParser(cron.NewParser(cron.Second|cron.Minute|cron.Hour|cron.Dom|cron.Month|cron.Dow)),
	)

	_, err = c.AddJob("0 * * * * *", trackingCron)
	if err != nil {
		logger.FatalKV(ctx, "trackingCron: AddJob", zap.Error(err))
	}

	c.Start()
	c.Entry(1).Job.Run()

	fmt.Println("Press Ctrl+C to exit...")

	<-ctx.Done()

	fmt.Println("\nShutdown signal received. Exiting...")
}

func getLogLevel() zapcore.Level {
	v := os.Getenv("LOG_LEVEL")
	if v == "" {
		return zap.WarnLevel
	}

	parsed, _ := strconv.Atoi(v)

	return zapcore.Level(parsed)
}
