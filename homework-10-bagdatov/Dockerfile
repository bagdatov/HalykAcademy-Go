FROM golang:1.16-alpine

RUN mkdir /app

ADD . /app

WORKDIR /app

RUN go mod download

RUN go build -o /web

EXPOSE 8080

ENV HTTP_PORT=8080

CMD ["/web"]