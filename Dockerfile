FROM golang:latest as builder

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
COPY . /go/src/backend

WORKDIR /go/src/backend

RUN go get github.com/line/line-bot-sdk-go/linebot
RUN go get github.com/joho/godotenv
RUN go get -u gorm.io/gorm
RUN go get gorm.io/driver/mysql
RUN go get github.com/golang-migrate/migrate
RUN go get github.com/golang-migrate/migrate/database/mysql
RUN go get github.com/golang-migrate/migrate/source/file
RUN go get github.com/pkg/errors

RUN GOOS=linux GOARCH=amd64 go build -o /main

# runtime image
FROM alpine
RUN apk update \
  && apk add --no-cache git curl make gcc g++
COPY --from=builder /main .

ENV PORT=${PORT}
ENTRYPOINT ["/main"]
