FROM golang:1.16-alpine

WORKDIR /app

COPY . ./
RUN apk add poppler-utils
RUN go mod download

RUN go build -o /main

EXPOSE 1323

CMD [ "/main" ]