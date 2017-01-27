FROM golang:latest
MAINTAINER Bassam Ismail <skippednote@gmail.com>

WORKDIR /app

RUN go build main.go

CMD ["./main.go"]
