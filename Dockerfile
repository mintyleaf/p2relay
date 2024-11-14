FROM golang:latest
LABEL maintainer="mintyleaf <mintyleafdev@gmail.com>"

RUN apt-get update

WORKDIR /app

COPY go.mod go.sum ./
COPY . .
RUN go mod download
RUN go build -o p2relay main.go

CMD ["./p2relay"]
