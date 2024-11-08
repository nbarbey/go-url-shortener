# syntax=docker/dockerfile:1

FROM golang:1.23 AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY urlshortener/*.go ./urlshortener/
COPY main.go ./

RUN GOOS=linux go build -o /url-shortener

FROM frolvlad/alpine-glibc AS run
COPY --from=build /url-shortener /url-shortener

EXPOSE 8080

CMD ["/url-shortener"]