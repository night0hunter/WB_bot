package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
	db "wb_bot/db"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
)

const (
	TimeFormat         = "02.01.2006"
	MoscowLocationName = "Europe/Moscow"
)

var moscowLocation *time.Location

var prevCommands = map[int64]BotCommandNameType{}

// var prevCMutex sync.RWMutex

type BotCommandNameType uint8

const (
	BotCommandNameTypeUnknown = iota
	BotCommandNameTypeInputDate
	BotCommandNameTypeInputWarehouse
	BotCommandNameTypeInputCoeffLimit
	BotCommandNameTypeInputSupplyType
)

var botCommands = map[uint8]string{
	BotCommandNameTypeInputDate:       "Введите дату отслеживания в следующем формате: \"дд.мм.гггг-дд.мм.гггг\"",
	BotCommandNameTypeInputWarehouse:  "Выберите склад, который хотите отслеживать",
	BotCommandNameTypeInputCoeffLimit: "Введите лимит коэффициента",
	BotCommandNameTypeInputSupplyType: "Выберите тип поставки",
}

var trackings = map[int64]db.WarehouseData{}

// var users = map[int64]User{}

// var usersMutex sync.RWMutexs

// var connString = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=verify-ca", user, password, host, port, dbname)

var connString = fmt.Sprintf(

	"postgresql://%s:%s@%s:%s/%s",
	"postgres",
	"pass123",
	"localhost",
	"5432",
	"wb_bot_db",
)

func init() {
	var err error
	moscowLocation, err = time.LoadLocation(MoscowLocationName)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	fmt.Println(os.Getenv("TELEGRAM_APITOKEN"))
	if err := godotenv.Load(); err != nil {
		log.Fatalf("godotenv.Load: %s", err)
	}

	// var (
	// 	host     = os.Getenv("HOST")
	// 	port     = os.Getenv("PORT")
	// 	user     = os.Getenv("USER")
	// 	password = os.Getenv("PASSWORD")
	// 	dbname   = os.Getenv("DBNAME")
	// )

	// var connString = fmt.Sprintf(
	// 	"postgresql://%s:%s@%s:%s/%s",
	// 	user,
	// 	password,
	// 	host,
	// 	port,
	// 	dbname,
	// )

	dbpool, err := db.NewPG(context.Background(), connString)
	if err != nil {
		log.Fatalf("db.NewPG: %s", err)
	}

	defer dbpool.Close()

	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))
	if err != nil {
		log.Fatalf("tgbotapi.NewBotAPI: %s", err)
	}

	updateConfig := tgbotapi.NewUpdate(0)

	updateConfig.Timeout = 30

	updates := bot.GetUpdatesChan(updateConfig)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		switch update.Message.Text {
		case "/start":
			// prevCMutex.Lock()

			prevCommands[update.Message.Chat.ID] = BotCommandNameTypeInputDate
			// prevCMutex.Unlock()

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, botCommands[BotCommandNameTypeInputDate])
			if _, err := bot.Send(msg); err != nil {
				fmt.Printf("bot.Send: %s\n", err.Error())
			}
		default:
			var msg tgbotapi.MessageConfig

			// prevCMutex.RLock()
			prevCommand, ok := prevCommands[update.Message.Chat.ID]
			// prevCMutex.RUnlock()

			if !ok {
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, "I don't know this command")
				if _, err := bot.Send(msg); err != nil {
					fmt.Printf("bot.Send: %s\n", err.Error())
				}

				break
			}

			switch prevCommand {
			case BotCommandNameTypeInputDate:
				dateFrom, dateTo, err := parseDate(update.Message.Text)
				if err != nil {
					msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Неверный формат даты")
					if _, err := bot.Send(msg); err != nil {
						fmt.Printf("bot.Send: %s\n", err.Error())
					}

					continue
				}

				// prevCMutex.Lock()
				prevCommands[update.Message.Chat.ID] = BotCommandNameTypeInputWarehouse
				// prevCMutex.Unlock()

				// usersMutex.Lock()
				trackings[update.Message.Chat.ID] = db.WarehouseData{ChatID: update.Message.Chat.ID}
				// usersMutex.Unlock()

				tmpTracking := trackings[update.Message.Chat.ID]

				tmpTracking.FromDate = dateFrom
				tmpTracking.ToDate = dateTo
				trackings[update.Message.Chat.ID] = tmpTracking

				msg = tgbotapi.NewMessage(update.Message.Chat.ID, botCommands[BotCommandNameTypeInputWarehouse])
				if _, err := bot.Send(msg); err != nil {
					fmt.Printf("bot.Send: %s\n", err.Error())
				}

			case BotCommandNameTypeInputWarehouse:
				// usersMutex.Lock()
				// tmpUser := users[update.Message.Chat.ID]
				// tmpUser.Surname = update.Message.Text
				// users[update.Message.Chat.ID] = tmpUser

				tmpTracking := trackings[update.Message.Chat.ID]
				tmpTracking.Warehouse = update.Message.Text
				trackings[update.Message.Chat.ID] = tmpTracking
				// usersMutex.Unlock()

				// prevCMutex.Lock()
				prevCommands[update.Message.Chat.ID] = BotCommandNameTypeInputCoeffLimit
				// prevCMutex.Unlock()

				msg = tgbotapi.NewMessage(update.Message.Chat.ID, botCommands[BotCommandNameTypeInputCoeffLimit])
				if _, err := bot.Send(msg); err != nil {
					fmt.Printf("bot.Send: %s\n", err.Error())
				}

			case BotCommandNameTypeInputCoeffLimit:
				tmpTracking := trackings[update.Message.Chat.ID]
				tmpTracking.CoeffLimit = update.Message.Text
				trackings[update.Message.Chat.ID] = tmpTracking

				prevCommands[update.Message.Chat.ID] = BotCommandNameTypeInputSupplyType

				msg = tgbotapi.NewMessage(update.Message.Chat.ID, botCommands[BotCommandNameTypeInputSupplyType])
				if _, err := bot.Send(msg); err != nil {
					fmt.Printf("bot.Send: %s\n", err.Error())
				}
			case BotCommandNameTypeInputSupplyType:
				tmpTracking := trackings[update.Message.Chat.ID]
				tmpTracking.SupplyType = update.Message.Text
				trackings[update.Message.Chat.ID] = tmpTracking

				err = dbpool.InsertQuery(context.Background(), trackings[update.Message.Chat.ID])
				if err != nil {
					fmt.Printf("dbpool.InsertQuery: %s\n", err.Error())
				}

				msg = tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("ChatID: %d\nDateFrom: %s\nDateTo: %s\nWarehouse: %s\nCoeffLim: %s\nSupplyType: %s\nIsActive: %t", trackings[update.Message.Chat.ID].ChatID, trackings[update.Message.Chat.ID].FromDate, trackings[update.Message.Chat.ID].ToDate, trackings[update.Message.Chat.ID].Warehouse, trackings[update.Message.Chat.ID].CoeffLimit, trackings[update.Message.Chat.ID].SupplyType, trackings[update.Message.Chat.ID].IsActive))
				if _, err := bot.Send(msg); err != nil {
					fmt.Printf("bot.Send: %s\n", err.Error())
				}
			default:
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Нет такого закона")
				bot.Send(msg)

			}
		}

	}

	bot.Debug = true
}

func parseDate(dateString string) (time.Time, time.Time, error) {
	datesRaw := strings.Split(dateString, "-")
	if len(datesRaw) != 2 {
		return time.Time{}, time.Time{}, errors.New("There must be 2 dates")
	}

	dateFrom, err := time.ParseInLocation(TimeFormat, datesRaw[0], moscowLocation)
	if err != nil {
		return time.Time{}, time.Time{}, errors.Wrap(err, "dateFrom: time.Parse")
	}

	dateTo, err := time.ParseInLocation(TimeFormat, datesRaw[1], moscowLocation)
	if err != nil {
		return time.Time{}, time.Time{}, errors.Wrap(err, "dateTo: time.Parse")
	}

	return dateFrom, dateTo, nil
}
