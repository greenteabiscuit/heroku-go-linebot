package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/line/line-bot-sdk-go/linebot"
)

// Geolocation ...
type Geolocation struct {
	Name string  `json:"name"`
	Lon  float64 `json:"lon"`
	Lat  float64 `json:"lat"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "what")
}

func lineHandler(w http.ResponseWriter, r *http.Request) {
	bot, err := linebot.New(
		os.Getenv("CHANNEL_SECRET"),
		os.Getenv("ACCESS_TOKEN"),
	)
	if err != nil {
		log.Fatal(err)
	}

	events, err := bot.ParseRequest(r)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			w.WriteHeader(400)
		} else {
			w.WriteHeader(500)
		}
		return
	}
	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				replyMessage := message.Text
				res := strings.Split(replyMessage, ",")
				replyMessage = res[0] + "から" + res[1] + "を移動しました"
				if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessage)).Do(); err != nil {
					log.Print(err)
				}
				key := os.Getenv("OPENWEATHER_API_KEY")
				url := "http://api.openweathermap.org/geo/1.0/direct?q=" + res[0] + "&limit=5&appid=" + key
				spaceClient := http.Client{
					Timeout: time.Second * 2, // Timeout after 2 seconds
				}
				req, err := http.NewRequest(http.MethodGet, url, nil)
				if err != nil {
					fmt.Fprintf(w, "bye\n")
				}
				req.Header.Set("User-Agent", "experiment")

				res, getErr := spaceClient.Do(req)
				if getErr != nil {
					log.Fatal(getErr)
				}

				if res.Body != nil {
					defer res.Body.Close()
				}

				body, readErr := ioutil.ReadAll(res.Body)
				if readErr != nil {
					log.Fatal(readErr)
				}
				var gl []Geolocation
				jsonErr := json.Unmarshal([]byte(body), &gl)
				if jsonErr != nil {
					log.Fatal(jsonErr)
				}
				replyMessageLon, replyMessageLat := gl[0].Lon, gl[0].Lat
				if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessageLon)).Do(); err != nil {
					log.Print(err)
				}
				if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessageLat)).Do(); err != nil {
					log.Print(err)
				}
			case *linebot.StickerMessage:
				replyMessage := fmt.Sprintf(
					"sticker id is %s, stickerResourceType is %s", message.StickerID, message.StickerResourceType)
				if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessage)).Do(); err != nil {
					log.Print(err)
				}
			}
		}
	}
}

func main() {
	port, _ := strconv.Atoi(os.Args[1])
	fmt.Printf("Starting server at Port %d", port)
	http.HandleFunc("/", handler)
	http.HandleFunc("/callback", lineHandler)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
