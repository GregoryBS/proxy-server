FROM golang:latest as build

WORKDIR /app

COPY main.go main.go
COPY go.mod go.mod

RUN go mod tidy
RUN go build -o main main.go

FROM ubuntu:latest
ENV TZ=Europe/Moscow
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

ENV DEBIAN_FRONTEND=noninteractive

COPY --from=build /app/main .

EXPOSE 8080
CMD ./main