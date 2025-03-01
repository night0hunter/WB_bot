// package api

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
	"wb_bot/db"
	"wb_bot/internal/dto"

	"github.com/davecgh/go-spew/spew"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

type Response struct {
	Date            time.Time `json:"date"`
	Coefficient     int       `json:"coefficient"`
	WarehouseID     int       `json:"warehouseID"`
	WarehouseName   string    `json:"warehouseName"`
	BoxTypeName     string    `json:"boxTypeName"`
	BoxTypeID       int       `json:"boxTypeId"`
	IsSortingCenter bool      `json:"isSortingCenter"`
}

type mergedResp struct {
	UserID          int64
	Date            time.Time
	Coefficient     int
	WarehouseID     int
	WarehouseName   string
	BoxTypeName     string
	BoxTypeID       int
	IsSortingCenter bool
	IsAvtive        bool
}

func GetWarehouseList() []Response {
	// url := "https://supplies-api-sandbox.wildberries.ru/api/v1/acceptance/coefficients"
	reqURL := os.Getenv("REQ_URL")

	// var bearer = "Bearer " + "eyJhbGciOiJFUzI1NiIsImtpZCI6IjIwMjQwOTA0djEiLCJ0eXAiOiJKV1QifQ.eyJlbnQiOjEsImV4cCI6MTc0MzIzMTM3MCwiaWQiOiIwMTkyMzRkNy0zMGM5LTc2ZDQtYTUxYy03MDlhZDViZjI0ZGIiLCJpaWQiOjQyNDA0MzM0LCJvaWQiOjEzOTA5MjIsInMiOjAsInNpZCI6IjYyYWMzZjVlLTZkODUtNDNkNy1iNTg0LTlmNjhmNzAwZjk0ZSIsInQiOnRydWUsInVpZCI6NDI0MDQzMzR9.Q9iFktjcoWCPGveWRH2zOwxwYW0tdQShZfVBgP0RzOoar2DiD1sLZU8i8WNf2JcZtt7sNHEbULc0QKfQ-hIs8Q"
	var bearer = "Bearer " + os.Getenv("BEARER_TOKEN")

	client := http.Client{
		Timeout: time.Second * 2,
	}

	req, err := http.NewRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		fmt.Printf("http.NewRequest: %s", err.Error())
	}

	req.Header.Add("Authorization", bearer)

	res, err := client.Do(req)
	if err != nil {
		fmt.Printf("client.Do: %s", err.Error())
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		fmt.Printf("io.ReadAll: %s", readErr.Error())
	}

	resp := []Response{}

	jsonErr := json.Unmarshal(body, &resp)
	if jsonErr != nil {
		fmt.Printf("json.Unmarshall: %s", jsonErr.Error())
	}

	// fmt.Printf("HTTP: %s\n", res.Status)

	sortedResp := make([]Response, 0, 20)
	for _, wh := range resp {
		if wh.BoxTypeID == 2 {
			sortedResp = append(sortedResp, wh)
		}
	}

	return sortedResp
}

func checkApprops(warehouses []Response, userTrackings []dto.WarehouseData) []mergedResp {
	var appropriate []mergedResp
	var tmp mergedResp

	for _, userTr := range userTrackings {
		for _, wh := range warehouses {
			if wh.Coefficient == -1 { // warehouse does not accept anything
				continue
			}

			if userTr.Warehouse != int64(wh.WarehouseID) {
				continue
			}

			if wh.Date.Before(userTr.FromDate) || wh.Date.After(userTr.ToDate) {
				continue
			}

			if userTr.CoeffLimit < wh.Coefficient {
				continue
			}

			if userTr.SupplyType != strconv.Itoa(wh.BoxTypeID) {
				continue
			}

			tmp.UserID = userTr.ChatID
			tmp.Date = wh.Date
			tmp.Coefficient = wh.Coefficient
			tmp.WarehouseID = wh.WarehouseID
			tmp.WarehouseName = wh.WarehouseName
			tmp.BoxTypeName = wh.BoxTypeName
			tmp.BoxTypeID = wh.BoxTypeID
			tmp.IsSortingCenter = wh.IsSortingCenter
			tmp.IsAvtive = userTr.IsActive

			appropriate = append(appropriate, tmp)
		}
	}

	return appropriate
}

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

	defer dbpool.Close()

	spew.Dump(GetWarehouseList())

	userTrackings, err := dbpool.SelectQuery(context.Background(), 1120114786)
	if err != nil {
		fmt.Printf("dbpool.SelectQuery: %s\n", err.Error())
	}

	// spew.Dump(checkApprops(GetWarehouseList(), userTrackings))

	if err := godotenv.Load(); err != nil {
		log.Fatalf("godotenv.Load: %s", err)
	}

	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))
	if err != nil {
		log.Fatalf("tgbotapi.NewBotAPI: %s", err)
	}

	fmt.Printf("Bot has been started on port %d ...\n", 5432)

	// c := cron.New()
	// c.AddFunc("* * * * * *", func() { sendMessage(bot, checkApprops(GetWarehouseList(), userTrackings)) })
	// c.AddFunc("0 * * * * *", func() { fmt.Println("Hello") })
	// c.Start()

	for {
		sendMessage(bot, checkApprops(GetWarehouseList(), userTrackings))

		time.Sleep(time.Minute)
	}
}

func sendMessage(bot *tgbotapi.BotAPI, warehouses []mergedResp) {
	if warehouses == nil {
		return
	}

	for i := 0; i < len(warehouses); i++ {
		if !warehouses[i].IsAvtive {
			continue
		}

		msg := tgbotapi.NewMessage(warehouses[i].UserID, fmt.Sprintf("Склад: %s\nДата: %s\nКоэффициент: x%d\nТип поставки: %s", warehouses[i].WarehouseName, warehouses[i].Date.Format("02.01.2006"), warehouses[i].Coefficient, warehouses[i].BoxTypeName))
		if _, err := bot.Send(msg); err != nil {
			fmt.Printf("bot.Send: %s\n", err.Error())
		}
	}
}
