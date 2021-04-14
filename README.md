# heroku-go-linebot

## Build

```
docker build -t linebot-goserver-exp .
```

## Run

```
docker run -e "PORT=3000" -p 3000:3000 -t linebot-goserver-exp
```

## Use locally

```
./ngrok http 3000
```
