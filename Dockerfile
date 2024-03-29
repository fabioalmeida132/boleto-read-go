FROM ubuntu:latest


RUN apt-get update
RUN apt-get install -y wget git gcc g++

RUN wget -P /tmp "https://dl.google.com/go/go1.19.linux-amd64.tar.gz"

RUN tar -C /usr/local -xzf "/tmp/go1.19.linux-amd64.tar.gz"
RUN rm "/tmp/go1.19.linux-amd64.tar.gz"

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH
RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"

WORKDIR /app

COPY . ./
RUN apt-get install -y poppler-utils
RUN apt-get -y install libzbar-dev
RUN apt-get install -y tesseract-ocr libtesseract-dev libleptonica-dev
RUN go mod tidy

RUN go build -o /main

CMD [ "/main" ]

#FROM golang:alpine
#
#WORKDIR /app
#
#COPY . ./
#RUN apk add --no-cache git gcc musl-dev
#RUN apk add poppler-utils
#RUN apk add zbar-dev
#RUN go mod download
#RUN go build -tags musl -o /main
#
#
#EXPOSE 1323
#
#CMD [ "/main" ]