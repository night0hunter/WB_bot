package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"os"
	"os/signal"
	"syscall"
	"time"
	"wb_bot/db"
	wbadp "wb_bot/internal/adapter/wb-adp"
	cronjob "wb_bot/internal/cronJob"
	"wb_bot/internal/handler"
	"wb_bot/internal/middleware"
	"wb_bot/internal/service"
	"wb_bot/internal/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"golang.org/x/net/publicsuffix"

	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	jar, err := cookiejar.New(&cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	})
	if err != nil {
		log.Fatal(ctx, "cookiejar.New")
	}

	if err := godotenv.Load(); err != nil {
		log.Fatalf("godotenv.Load: %s", err)
	}

	var (
		cookie      = os.Getenv("COOKIE")
		authorizeV3 = os.Getenv("AUTHORIZE_V3")
	)

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

	dbpool, err := db.NewPG(ctx, connString)
	if err != nil {
		log.Fatalf("db.NewPG: %s", err)
	}

	fmt.Printf("Base has been started on port %s\n", port)

	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))
	if err != nil {
		log.Fatalf("tgbotapi.NewBotAPI: %s", err)
	}

	fmt.Printf("Bot has been started\n")

	httpClient := &http.Client{
		Timeout: 10 * time.Second,
		Transport: middleware.Chain(
			nil,
			middleware.SetAccCookieMiddleware(cookie),
			middleware.SetHeader(middleware.Authorizev3Header, authorizeV3),
		),
		Jar: jar,
	}
	// mockClient := &mock.MockClient{}

	wbSupplyAdp := wbadp.New(
		httpClient,
		"https://seller-supply.wildberries.ru",
	)

	service := service.New(dbpool, wbSupplyAdp)
	handler := handler.New(bot, service)
	bookDraftCron := cronjob.NewBookDraftCron(handler)

	c := cron.New(
		cron.WithLocation(utils.MoscowLocation),
		cron.WithParser(cron.NewParser(cron.Second|cron.Minute|cron.Hour|cron.Dom|cron.Month|cron.Dow)),
	)

	_, err = c.AddJob("0 * * * * *", bookDraftCron)
	if err != nil {
		fmt.Printf("c.AddJob: %s", err.Error())
	}

	c.Start()
	c.Entry(1).Job.Run()

	fmt.Println("Press Ctrl+C to exit...")

	<-ctx.Done()

	fmt.Println("\nShutdown signal received. Exiting...")
}
