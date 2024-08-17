FROM golang:1.22 AS builder
WORKDIR /go/src/github.com/OwO-Network/deeplx-bot
COPY bot.go ./
COPY go.mod ./
COPY go.sum ./
RUN go get -d -v ./
RUN CGO_ENABLED=0 go build -a -installsuffix cgo -o deeplx-bot .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /go/src/github.com/OwO-Network/deeplx-bot/deeplx-bot /app/deeplx-bot
CMD ["/app/deeplx-bot"]