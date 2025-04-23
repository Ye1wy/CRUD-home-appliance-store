FROM golang:1.24.4-alpine AS build

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./

COPY . .

RUN go build -o /app/crud-service ./cmd/crud/main.go

FROM alpine:latest AS runtime

WORKDIR /service

COPY --from=build /app/crud-service .
COPY --from=build /app/.env ./.env

CMD ["./crud-service"]
