package main

// import (
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"log"
// 	"os"
// 	"strconv"
// 	"strings"
// 	"time"
// 	db "wb_bot/db"

// 	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
// 	"github.com/joho/godotenv"
// 	"github.com/pkg/errors"
// )

// const (
// 	TimeFormat         = "02.01.2006"
// 	MoscowLocationName = "Europe/Moscow"
// )

// var moscowLocation *time.Location

// var prevCommands = map[int64]BotCommandNameType{}

// // var prevCMutex sync.RWMutex

// type BotCommandNameType uint8

// const (
// 	BotCommandNameTypeUnknown BotCommandNameType = iota
// 	BotCommandNameTypeInputDate
// 	BotCommandNameTypeInputWarehouse
// 	BotCommandNameTypeInputCoeffLimit
// 	BotCommandNameTypeInputSupplyType
// )

// var botCommands = map[BotCommandNameType]string{
// 	BotCommandNameTypeInputDate:       "Введите дату отслеживания в следующем формате: \"дд.мм.гггг-дд.мм.гггг\"",
// 	BotCommandNameTypeInputWarehouse:  "Выберите склад, который хотите отслеживать",
// 	BotCommandNameTypeInputCoeffLimit: "Выберите лимит коэффициента или введите свой",
// 	BotCommandNameTypeInputSupplyType: "Выберите тип поставки",
// }

// var trackings = map[int64]db.WarehouseData{}

// type ButtonType uint8

// const (
// 	ButtonTypeCoeffLimit ButtonType = iota + 1
// 	ButtonTypeWarehouse
// 	ButtonTypeSupplyType
// 	ButtonTypeUserTrackings
// 	ButtonTypeUserTrackingStatus
// )

// type ButtonData struct {
// 	Type  ButtonType
// 	Value int
// }

// type Button struct {
// 	Data ButtonData
// 	Text string
// }

// func init() {
// 	var err error
// 	moscowLocation, err = time.LoadLocation(MoscowLocationName)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }

// func main() {
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

// 	dbpool, err := db.NewPG(context.Background(), connString)
// 	if err != nil {
// 		log.Fatalf("db.NewPG: %s", err)
// 	}

// 	defer dbpool.Close()

// 	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))
// 	if err != nil {
// 		log.Fatalf("tgbotapi.NewBotAPI: %s", err)
// 	}

// 	fmt.Printf("Bot has been started on port %s ...", port)

// 	updateConfig := tgbotapi.NewUpdate(0)

// 	updateConfig.Timeout = 30

// 	updates := bot.GetUpdatesChan(updateConfig)

// 	for update := range updates {
// 		if update.Message == nil && update.CallbackQuery == nil {
// 			continue
// 		}

// 		if update.CallbackQuery != nil {
// 			var buttonData ButtonData

// 			err := json.Unmarshal([]byte(update.CallbackQuery.Data), &buttonData)
// 			if err != nil {
// 				fmt.Printf("json.unMarshal: %s\n", err.Error())
// 			}

// 			switch buttonData.Type {
// 			case ButtonTypeWarehouse:
// 				prevCommands[update.CallbackQuery.Message.Chat.ID] = BotCommandNameTypeInputCoeffLimit

// 				tmpTracking := trackings[update.CallbackQuery.Message.Chat.ID]
// 				tmpTracking.Warehouse = buttonData.Value
// 				trackings[update.CallbackQuery.Message.Chat.ID] = tmpTracking

// 				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, fmt.Sprintf("Вы выбрали склад %d", buttonData.Value))
// 				if _, err := bot.Send(msg); err != nil {
// 					fmt.Printf("bot.Send: %s\n", err.Error())
// 				}

// 				msg = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, botCommands[BotCommandNameTypeInputCoeffLimit])
// 				if _, err := bot.Send(drawCoeffKeyboard(msg)); err != nil {
// 					fmt.Printf("bot.Send: %s\n", err.Error())
// 				}

// 				continue
// 			case ButtonTypeCoeffLimit:
// 				prevCommands[update.CallbackQuery.Message.Chat.ID] = BotCommandNameTypeInputSupplyType

// 				tmpTracking := trackings[update.CallbackQuery.Message.Chat.ID]
// 				tmpTracking.CoeffLimit = buttonData.Value
// 				trackings[update.CallbackQuery.Message.Chat.ID] = tmpTracking

// 				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, fmt.Sprintf("Вы выбрали %dx", buttonData.Value))
// 				if _, err := bot.Send(msg); err != nil {
// 					fmt.Printf("bot.Send: %s\n", err.Error())
// 				}

// 				msg = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, botCommands[BotCommandNameTypeInputSupplyType])
// 				if _, err := bot.Send(drawSupplyKeyboard(msg)); err != nil {
// 					fmt.Printf("bot.Send: %s\n", err.Error())
// 				}

// 				continue

// 			case ButtonTypeSupplyType:
// 				tmpTracking := trackings[update.CallbackQuery.Message.Chat.ID]
// 				tmpTracking.SupplyType = fmt.Sprint(buttonData.Value)
// 				trackings[update.CallbackQuery.Message.Chat.ID] = tmpTracking

// 				prevCommands[update.CallbackQuery.Message.Chat.ID] = BotCommandNameTypeInputDate

// 				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, fmt.Sprintf("Вы выбрали тип поставки %d", buttonData.Value))
// 				if _, err := bot.Send(msg); err != nil {
// 					fmt.Printf("bot.Send: %s\n", err.Error())
// 				}

// 				err = dbpool.InsertQuery(context.Background(), trackings[update.CallbackQuery.Message.Chat.ID])
// 				if err != nil {
// 					fmt.Printf("dbpool.InsertQuery: %s\n", err.Error())
// 				}

// 				msg = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Склад успешно добавлен!")
// 				if _, err := bot.Send(msg); err != nil {
// 					fmt.Printf("bot.Send: %s\n", err.Error())
// 				}

// 				continue
// 			case ButtonTypeUserTrackings:
// 				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Выберите действие")
// 				if _, err := bot.Send(drawTrackingStatusKeyboard(msg)); err != nil {
// 					fmt.Printf("bot.Send: %s\n", err.Error())
// 				}

// 				continue
// 			case ButtonTypeUserTrackingStatus:
// 				var msg tgbotapi.MessageConfig

// 				if buttonData.Value == 0 {
// 					msg = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Статус отслеживания успешно изменён!")
// 				}

// 				if buttonData.Value == 1 {
// 					msg = tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Отслеживание успешно удалено!")
// 				}

// 				if _, err := bot.Send(drawTrackingStatusKeyboard(msg)); err != nil {
// 					fmt.Printf("bot.Send: %v\n", err)

// 				}

// 				continue
// 			}
// 		}

// 		switch update.Message.Text {
// 		// case "test":
// 		// fmt.Println(warehouseList)
// 		case "/add":
// 			// prevCMutex.Lock()

// 			prevCommands[update.Message.Chat.ID] = BotCommandNameTypeInputDate
// 			// prevCMutex.Unlock()

// 			msg := tgbotapi.NewMessage(update.Message.Chat.ID, botCommands[BotCommandNameTypeInputDate])
// 			if _, err := bot.Send(msg); err != nil {
// 				fmt.Printf("bot.Send: %s\n", err.Error())
// 			}
// 		case "/check":
// 			var warehouses []db.WarehouseData
// 			var msg tgbotapi.MessageConfig

// 			warehouses, err = dbpool.SelectQuery(context.Background(), update.Message.Chat.ID)
// 			if err != nil {
// 				fmt.Printf("dbpool.SelectQuery: %s\n", err.Error())
// 			}

// 			if len(warehouses) == 0 {
// 				msg = tgbotapi.NewMessage(update.Message.Chat.ID, "На данный момент вы не отслеживаете ни одного склада.\nЧтобы добавить, введите /add")
// 				if _, err := bot.Send(msg); err != nil {
// 					fmt.Printf("bot.Send: %s\n", err.Error())
// 				}

// 				continue
// 			}

// 			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Список отслеживаемых складов:")
// 			if _, err := bot.Send(msg); err != nil {
// 				fmt.Printf("bot.Send: %s\n", err.Error())
// 			}

// 			for _, wh := range warehouses {
// 				isActive := "Активно"
// 				if !wh.IsActive {
// 					isActive = "Неактивно"
// 				}

// 				msg = tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Склад: %d\nДата отслеживания: %s-%s\nЛимит коэффициента: x%d и меньше\nТип поставки: %s\nАктивно/Неактивно: %s", wh.Warehouse, wh.FromDate.Format(TimeFormat), wh.ToDate.Format(TimeFormat), wh.CoeffLimit, wh.SupplyType, isActive))
// 				if _, err := bot.Send(msg); err != nil {
// 					fmt.Printf("bot.Send: %s\n", err.Error())
// 				}
// 			}
// 		case "/help":
// 			var msg tgbotapi.MessageConfig

// 			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Hello world")
// 			if _, err := bot.Send(msg); err != nil {
// 				fmt.Printf("bot.Send: %s\n", err.Error())
// 			}
// 			// tgbotapi.NewDeleteMessage(update.Message.Chat.ID, msg.ReplyToMessageID)
// 		case "/stop":
// 			var warehouses []db.WarehouseData
// 			var msg tgbotapi.MessageConfig

// 			warehouses, err = dbpool.SelectQuery(context.Background(), update.Message.Chat.ID)
// 			if err != nil {
// 				fmt.Printf("dbpool.SelectQuery: %s\n", err.Error())
// 			}

// 			if len(warehouses) == 0 {
// 				msg = tgbotapi.NewMessage(update.Message.Chat.ID, "На данный момент вы не отслеживаете ни одного склада.\nЧтобы добавить, введите /add")
// 				if _, err := bot.Send(msg); err != nil {
// 					fmt.Printf("bot.Send: %s\n", err.Error())
// 				}

// 				continue
// 			}

// 			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Выберите отслеживание из списка ниже, чтобы приостановить/удалить")
// 			if _, err := bot.Send(drawTrackingsKeyboard(msg, warehouses)); err != nil {
// 				fmt.Printf("bot.Send: %s\n", err.Error())
// 			}

// 		default:
// 			var msg tgbotapi.MessageConfig

// 			// prevCMutex.RLock()
// 			prevCommand, ok := prevCommands[update.Message.Chat.ID]
// 			// prevCMutex.RUnlock()

// 			if !ok {
// 				msg = tgbotapi.NewMessage(update.Message.Chat.ID, "I don't know this command")
// 				if _, err := bot.Send(msg); err != nil {
// 					fmt.Printf("bot.Send: %s\n", err.Error())
// 				}

// 				break
// 			}

// 			switch prevCommand {
// 			case BotCommandNameTypeInputDate:
// 				dateFrom, dateTo, err := parseDate(update.Message.Text)
// 				if err != nil {
// 					msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Неверный формат даты, попробуйте ещё раз")
// 					if _, err := bot.Send(msg); err != nil {
// 						fmt.Printf("bot.Send: %s\n", err.Error())
// 					}

// 					continue
// 				}

// 				// prevCMutex.Lock()
// 				prevCommands[update.Message.Chat.ID] = BotCommandNameTypeInputWarehouse
// 				// prevCMutex.Unlock()

// 				// usersMutex.Lock()
// 				trackings[update.Message.Chat.ID] = db.WarehouseData{ChatID: update.Message.Chat.ID}
// 				// usersMutex.Unlock()

// 				tmpTracking := trackings[update.Message.Chat.ID]

// 				tmpTracking.FromDate = dateFrom
// 				tmpTracking.ToDate = dateTo
// 				trackings[update.Message.Chat.ID] = tmpTracking

// 				msg = tgbotapi.NewMessage(update.Message.Chat.ID, botCommands[BotCommandNameTypeInputWarehouse])
// 				if _, err := bot.Send((msg)); err != nil {
// 					fmt.Printf("bot.Send: %s\n", err.Error())
// 				}

// 				msg = tgbotapi.NewMessage(update.Message.Chat.ID, botCommands[BotCommandNameTypeInputWarehouse])
// 				if _, err := bot.Send(drawWarehouseKeyboard(msg)); err != nil {
// 					fmt.Printf("bot.Send: %s\n", err.Error())
// 				}

// 			case BotCommandNameTypeInputWarehouse:
// 				msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Пожалуйста, выберите нужный склад из предложенного списка")
// 				if _, err := bot.Send((msg)); err != nil {
// 					fmt.Printf("bot.Send: %s\n", err.Error())
// 				}

// 				continue
// 			case BotCommandNameTypeInputCoeffLimit:
// 				err := parseCoeffLimit(update.Message.Text)
// 				if err != nil {
// 					msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Неверный формат коэффициента, попробуйте ещё раз")
// 					if _, err := bot.Send(msg); err != nil {
// 						fmt.Printf("bot.Send: %s\n", err.Error())
// 					}

// 					continue
// 				}

// 				tmpTracking := trackings[update.Message.Chat.ID]
// 				tmp, err := strconv.Atoi(update.Message.Text)
// 				if err != nil {
// 					fmt.Printf("strconv.Atoi(coeff input): %s\n", err.Error())
// 				}

// 				tmpTracking.CoeffLimit = tmp
// 				trackings[update.Message.Chat.ID] = tmpTracking

// 				prevCommands[update.Message.Chat.ID] = BotCommandNameTypeInputSupplyType

// 				msg = tgbotapi.NewMessage(update.Message.Chat.ID, botCommands[BotCommandNameTypeInputSupplyType])
// 				if _, err := bot.Send(drawSupplyKeyboard(msg)); err != nil {
// 					fmt.Printf("bot.Send: %s\n", err.Error())
// 				}

// 			case BotCommandNameTypeInputSupplyType:
// 				msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Пожалуйста, выберите нужный тип поставки из предложенного списка")
// 				if _, err := bot.Send(drawSupplyKeyboard(msg)); err != nil {
// 					fmt.Printf("bot.Send: %s\n", err.Error())
// 				}
// 			default:
// 				msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Нет такого закона")
// 				bot.Send(msg)

// 			}
// 		}

// 	}

// 	bot.Debug = true
// }

// func parseDate(dateString string) (time.Time, time.Time, error) {
// 	datesRaw := strings.Split(dateString, "-")
// 	if len(datesRaw) != 2 {
// 		return time.Time{}, time.Time{}, errors.New("There must be 2 dates")
// 	}

// 	dateFrom, err := time.ParseInLocation(TimeFormat, datesRaw[0], moscowLocation)
// 	if err != nil {
// 		return time.Time{}, time.Time{}, errors.Wrap(err, "dateFrom: time.ParseInLocation")
// 	}

// 	dateTo, err := time.ParseInLocation(TimeFormat, datesRaw[1], moscowLocation)
// 	if err != nil {
// 		return time.Time{}, time.Time{}, errors.Wrap(err, "dateTo: time.ParseInLocation")
// 	}

// 	return dateFrom, dateTo, nil
// }

// func parseCoeffLimit(coeff string) error {
// 	_, err := strconv.Atoi(coeff)
// 	if err != nil {
// 		return errors.Wrap(err, "coeffLimit: strconv.Atoi")
// 	}

// 	return nil
// }

// func generateKeyboard(buttons ...Button) (tgbotapi.InlineKeyboardMarkup, error) {
// 	keyboardButtons := make([]tgbotapi.InlineKeyboardButton, 0, len(buttons))

// 	for _, button := range buttons {
// 		jsonData, err := json.Marshal(button.Data)
// 		if err != nil {
// 			return tgbotapi.InlineKeyboardMarkup{}, errors.Wrap(err, "json.Marshal")
// 		}

// 		keyboardButtons = append(keyboardButtons, tgbotapi.NewInlineKeyboardButtonData(button.Text, string(jsonData)))
// 	}

// 	return tgbotapi.NewInlineKeyboardMarkup(
// 		tgbotapi.NewInlineKeyboardRow(keyboardButtons...),
// 	), nil
// }

// func drawWarehouseKeyboard(msg tgbotapi.MessageConfig) tgbotapi.MessageConfig {
// 	tmpMarkup, err := generateKeyboard([]Button{
// 		{
// 			Data: ButtonData{
// 				Type:  ButtonTypeWarehouse,
// 				Value: 507, // wareID
// 			},
// 			Text: "Коледино", // wareName
// 		},
// 		{
// 			Data: ButtonData{
// 				Type:  ButtonTypeWarehouse,
// 				Value: 316646,
// 			},
// 			Text: "Шушары СГТ",
// 		},
// 		{
// 			Data: ButtonData{
// 				Type:  ButtonTypeWarehouse,
// 				Value: 301229,
// 			},
// 			Text: "Подольск 4",
// 		},
// 		{
// 			Data: ButtonData{
// 				Type:  ButtonTypeWarehouse,
// 				Value: 120762,
// 			},
// 			Text: "Электросталь",
// 		},
// 		{
// 			Data: ButtonData{
// 				Type:  ButtonTypeWarehouse,
// 				Value: 206348,
// 			},
// 			Text: "Тула",
// 		},
// 		{
// 			Data: ButtonData{
// 				Type:  ButtonTypeWarehouse,
// 				Value: 130744,
// 			},
// 			Text: "Краснодар",
// 		},
// 		{
// 			Data: ButtonData{
// 				Type:  ButtonTypeWarehouse,
// 				Value: 208277,
// 			},
// 			Text: "Невинномысск",
// 		},
// 		{
// 			Data: ButtonData{
// 				Type:  ButtonTypeWarehouse,
// 				Value: 117986,
// 			},
// 			Text: "Казань",
// 		},
// 		{
// 			Data: ButtonData{
// 				Type:  ButtonTypeWarehouse,
// 				Value: 117501,
// 			},
// 			Text: "Подольск",
// 		},
// 		{
// 			Data: ButtonData{
// 				Type:  ButtonTypeWarehouse,
// 				Value: 1733,
// 			},
// 			Text: "Екатеринбург - Испытателей 14г",
// 		},
// 		{
// 			Data: ButtonData{
// 				Type:  ButtonTypeWarehouse,
// 				Value: 218644,
// 			},
// 			Text: "СЦ Хабаровск", // add "Хабаровск"
// 		},
// 		{
// 			Data: ButtonData{
// 				Type:  ButtonTypeWarehouse,
// 				Value: 218644, // find id
// 			},
// 			Text: "Санкт-Петербург Уткина Заводь",
// 		},
// 		{
// 			Data: ButtonData{
// 				Type:  ButtonTypeWarehouse,
// 				Value: 206236,
// 			},
// 			Text: "Белые Столбы",
// 		},
// 		{
// 			Data: ButtonData{
// 				Type:  ButtonTypeWarehouse,
// 				Value: 686,
// 			},
// 			Text: "Новосибирск",
// 		},
// 	}...)
// 	if err != nil {
// 		fmt.Printf("generateKeyboard: %s\n", err.Error())
// 	}

// 	msg.ReplyMarkup = tmpMarkup

// 	return msg

// }

// func drawCoeffKeyboard(msg tgbotapi.MessageConfig) tgbotapi.MessageConfig {
// 	tmpMarkup, err := generateKeyboard([]Button{
// 		{
// 			Data: ButtonData{
// 				Type:  ButtonTypeCoeffLimit,
// 				Value: 0,
// 			},
// 			Text: "Бесплатно",
// 		},
// 		{
// 			Data: ButtonData{
// 				Type:  ButtonTypeCoeffLimit,
// 				Value: 1,
// 			},
// 			Text: "1x",
// 		},
// 		{
// 			Data: ButtonData{
// 				Type:  ButtonTypeCoeffLimit,
// 				Value: 2,
// 			},
// 			Text: "2x",
// 		},
// 		{
// 			Data: ButtonData{
// 				Type:  ButtonTypeCoeffLimit,
// 				Value: 3,
// 			},
// 			Text: "3x",
// 		},
// 	}...)
// 	if err != nil {
// 		fmt.Printf("generateKeyboard: %s\n", err.Error())
// 	}

// 	msg.ReplyMarkup = tmpMarkup

// 	return msg
// }

// func drawSupplyKeyboard(msg tgbotapi.MessageConfig) tgbotapi.MessageConfig {
// 	tmpMarkup, err := generateKeyboard([]Button{
// 		{
// 			Data: ButtonData{
// 				Type:  ButtonTypeSupplyType,
// 				Value: 2,
// 			},
// 			Text: "Короб",
// 		},
// 		{
// 			Data: ButtonData{
// 				Type:  ButtonTypeSupplyType,
// 				Value: 5,
// 			},
// 			Text: "Монопалеты",
// 		},
// 		{
// 			Data: ButtonData{
// 				Type:  ButtonTypeSupplyType,
// 				Value: 6,
// 			},
// 			Text: "Суперсейф",
// 		},
// 	}...)
// 	if err != nil {
// 		fmt.Printf("generateKeyboard: %s\n", err.Error())
// 	}

// 	msg.ReplyMarkup = tmpMarkup

// 	return msg
// }

// func drawTrackingsKeyboard(msg tgbotapi.MessageConfig, warehouses []db.WarehouseData) tgbotapi.MessageConfig {
// 	var buttons []Button
// 	var button Button

// 	for _, wh := range warehouses {
// 		button.Text = strconv.Itoa(wh.Warehouse) + " " + wh.FromDate.Format(TimeFormat) + "-" + wh.ToDate.Format(TimeFormat)
// 		button.Data.Type = ButtonTypeUserTrackings
// 		button.Data.Value = int(wh.TrackingID)

// 		buttons = append(buttons, button)
// 	}

// 	tmpMarkup, err := generateKeyboard(buttons...)
// 	if err != nil {
// 		fmt.Printf("generateKeyboard: %s\n", err.Error())
// 	}

// 	msg.ReplyMarkup = tmpMarkup

// 	return msg
// }

// // todo: add adaptive message depended on IsActive field
// func drawTrackingStatusKeyboard(msg tgbotapi.MessageConfig) tgbotapi.MessageConfig {
// 	tmpMarkup, err := generateKeyboard([]Button{
// 		{
// 			Data: ButtonData{
// 				Type:  ButtonTypeUserTrackingStatus,
// 				Value: 0,
// 			},
// 			Text: "Приостановить/Возобновить",
// 		},
// 		{
// 			Data: ButtonData{
// 				Type:  ButtonTypeUserTrackingStatus,
// 				Value: 1,
// 			},
// 			Text: "Удалить",
// 		},
// 	}...)
// 	if err != nil {
// 		fmt.Printf("generateKeyboard: %s\n", err.Error())
// 	}

// 	msg.ReplyMarkup = tmpMarkup

// 	return msg
// }
