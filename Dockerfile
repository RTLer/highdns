FROM golang:1.18-alpine3.16 AS build
RUN apk add upx

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY ./cmd      ./cmd
COPY ./dns      ./dns
COPY ./main.go  ./main.go

RUN go build -o /highdns -ldflags "-s -w"
RUN upx /highdns


FROM alpine:3.16.1 as app
WORKDIR /
COPY --from=build /highdns /highdns
COPY ./config.yaml /config.yaml

EXPOSE 5354

ENTRYPOINT ["/highdns", "serv"]
