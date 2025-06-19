package main

// import (
// 	"context"
// 	"fmt"
// 	"log"
// 	"os"
// 	"os/signal"
// 	"syscall"
// 	"wb_bot/db"
// 	cronjob "wb_bot/internal/cronJob"
// 	"wb_bot/internal/handler"
// 	"wb_bot/internal/service"
// 	"wb_bot/internal/utils"

// 	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
// 	"github.com/robfig/cron/v3"

// 	"github.com/joho/godotenv"
// )

// func main() {
// 	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
// 	defer stop()

// 	if err := godotenv.Load(); err != nil {
// 		log.Fatalf("godotenv.Load: %s", err)
// 	}

// 	var (
// 		host     = os.Getenv("HOST")
// 		port     = os.Getenv("PORT")
// 		user     = os.Getenv("USER")
// 		password = os.Getenv("PASSWORD")
// 		dbname   = os.Getenv("DBNAME")
// 	)

// 	var connString = fmt.Sprintf(
// 		"postgresql://%s:%s@%s:%s/%s",
// 		user,
// 		password,
// 		host,
// 		port,
// 		dbname,
// 	)

// 	dbpool, err := db.NewPG(ctx, connString)
// 	if err != nil {
// 		log.Fatalf("db.NewPG: %s", err)
// 	}

// 	fmt.Printf("Base has been started on port %s\n", port)

// 	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))
// 	if err != nil {
// 		log.Fatalf("tgbotapi.NewBotAPI: %s", err)
// 	}

// 	fmt.Printf("Bot has been started\n")

// 	service := service.NewService(dbpool)
// 	handler := handler.NewHandler(bot, service)
// 	trackingCron := cronjob.NewSendTrackingCron(handler)

// 	c := cron.New(
// 		cron.WithLocation(utils.MoscowLocation),
// 		cron.WithParser(cron.NewParser(cron.Second|cron.Minute|cron.Hour|cron.Dom|cron.Month|cron.Dow)),
// 	)

// 	_, err = c.AddJob("0 * * * * *", trackingCron)
// 	if err != nil {
// 		fmt.Printf("c.AddJob: %s", err.Error())
// 	}

// 	c.Start()
// 	c.Entry(1).Job.Run()

// 	fmt.Println("Press Ctrl+C to exit...")

// 	<-ctx.Done()

// 	fmt.Println("\nShutdown signal received. Exiting...")
// }
