FROM golang:1.24.0-alpine3.21 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o /app/cmd/server/bin /app/cmd/server/main.go

FROM alpine:3.21.3 AS server

WORKDIR /app

COPY --from=builder /app/cmd/server/bin /app/.env.production ./
COPY --from=builder /app/migration /app/migration

EXPOSE 8000

ENTRYPOINT [ "./server/bin" ]