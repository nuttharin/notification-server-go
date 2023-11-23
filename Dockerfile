# syntax=docker/dockerfile:1

FROM golang:1.20.0

WORKDIR /app

COPY . .
RUN go mod tidy

RUN go build -o /docker-dz-notification-server

EXPOSE 8004

CMD [ "/docker-dz-notification-server" ]