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
				key := os.Getenv("OPENWEATHER_API_KEY")
				url1 := "http://api.openweathermap.org/geo/1.0/direct?q=" + res[0] + "&limit=5&appid=" + key
				url2 := "http://api.openweathermap.org/geo/1.0/direct?q=" + res[1] + "&limit=5&appid=" + key
				client1 := http.Client{
					Timeout: time.Second * 2, // Timeout after 2 seconds
				}
				client2 := http.Client{
					Timeout: time.Second * 2, // Timeout after 2 seconds
				}
				req1, _ := http.NewRequest(http.MethodGet, url1, nil)
				req2, _ := http.NewRequest(http.MethodGet, url2, nil)

				req1.Header.Set("User-Agent", "experiment")
				req2.Header.Set("User-Agent", "experiment")

				httpResponse1, _ := client1.Do(req1)
				httpResponse2, _ := client2.Do(req2)

				if httpResponse1.Body != nil {
					defer httpResponse1.Body.Close()
				}

				if httpResponse2.Body != nil {
					defer httpResponse2.Body.Close()
				}

				body1, _ := ioutil.ReadAll(httpResponse1.Body)
				body2, _ := ioutil.ReadAll(httpResponse2.Body)

				var gl1 []Geolocation
				jsonErr1 := json.Unmarshal([]byte(body1), &gl1)
				if jsonErr1 != nil {
					log.Fatal(jsonErr1)
				}
				var gl2 []Geolocation
				jsonErr2 := json.Unmarshal([]byte(body2), &gl2)
				if jsonErr2 != nil {
					log.Fatal(jsonErr2)
				}
				lon1, lat1 := gl1[0].Lon, gl1[0].Lat
				lonStr1, latStr1 := strconv.FormatFloat(lon1, 'f', 2, 64), strconv.FormatFloat(lat1, 'f', 2, 64)
				lon2, lat2 := gl2[0].Lon, gl2[0].Lat
				lonStr2, latStr2 := strconv.FormatFloat(lon2, 'f', 2, 64), strconv.FormatFloat(lat2, 'f', 2, 64)

				replyMessage = res[0] + "から" + res[1] + "を移動しました。" + res[0] + "の緯度経度は" + lonStr1 + ", " + latStr1 + "です。" + res[1] + "の緯度経度は" + lonStr2 + ", " + latStr2 + "です。"
				if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessage)).Do(); err != nil {
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
