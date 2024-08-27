FROM golang:alpine3.20 AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download 

COPY . ./

RUN go build -o /traderbot

FROM alpine:3.20 AS final

WORKDIR /

COPY --from=build /traderbot /traderbot

ENTRYPOINT ["/traderbot"]

