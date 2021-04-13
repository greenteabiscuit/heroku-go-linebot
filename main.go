package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
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
	// localであればlocalの.envを読み込む
	if os.Getenv("USE_HEROKU") != "1" {
		err := godotenv.Load()
		if err != nil {
			panic(err)
		}
	}
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

				distStr := strconv.Itoa(int(distance(lat1, lon1, lat2, lon2, "K")))
				replyMessage = res[0] + "から" + res[1] + "を移動しました。" + res[0] + "の緯度経度は" + lonStr1 + ", " + latStr1 + "です。" + res[1] + "の緯度経度は" + lonStr2 + ", " + latStr2 + "です。" + "距離は" + distStr + "kmです。"
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

func distance(lat1 float64, lng1 float64, lat2 float64, lng2 float64, unit ...string) float64 {
	const PI float64 = 3.141592653589793

	radlat1 := float64(PI * lat1 / 180)
	radlat2 := float64(PI * lat2 / 180)

	theta := float64(lng1 - lng2)
	radtheta := float64(PI * theta / 180)

	dist := math.Sin(radlat1)*math.Sin(radlat2) + math.Cos(radlat1)*math.Cos(radlat2)*math.Cos(radtheta)

	if dist > 1 {
		dist = 1
	}

	dist = math.Acos(dist)
	dist = dist * 180 / PI
	dist = dist * 60 * 1.1515

	if len(unit) > 0 {
		if unit[0] == "K" {
			dist = dist * 1.609344
		} else if unit[0] == "N" {
			dist = dist * 0.8684
		}
	}

	return dist
}

func main() {
	port, _ := strconv.Atoi(os.Args[1])
	fmt.Printf("Starting server at Port %d", port)
	http.HandleFunc("/", handler)
	http.HandleFunc("/callback", lineHandler)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
