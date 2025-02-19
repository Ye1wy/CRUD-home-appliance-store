FROM golang:1.23.4 AS build
WORKDIR /go/src
COPY

ENV CGO_ENABLED=0

RUN go build -o openapi .

FROM scratch AS runtime
ENV GIN_MODE=release
COPY --from=build /go/src/openapi ./
EXPOSE 8080/tcp
ENTRYPOINT ["./openapi"]
