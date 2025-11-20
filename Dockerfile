FROM golang:1.24 AS builder

WORKDIR /app

# кешируем зависимости
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o subtrack cmd/subtrack/main.go

FROM alpine:3.20

WORKDIR /app

RUN apk --no-cache add ca-certificates

COPY --from=builder /app/subtrack .

COPY migrations ./migrations

ENV HTTP_ADDRESS=:8080

CMD ["./subtrack"]
