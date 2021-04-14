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
RUN go build main.go

# runtime image
FROM alpine
RUN apk update \
  && apk add --no-cache git curl make gcc g++
COPY --from=builder /go/src/backend /app
WORKDIR /app

CMD /app/main $PORT
