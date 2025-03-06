FROM golang:1.24.0-alpine AS build

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./

COPY . .

RUN go build -o /app/bin/app ./cmd/main.go

FROM alpine:latest AS runtime

WORKDIR /service

COPY --from=build /app/bin/app .
COPY --from=build /app/.config-mongodb.env ./.config-mongodb.env

CMD ["./app"]
