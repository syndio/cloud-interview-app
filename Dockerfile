FROM golang:latest

RUN go install github.com/cespare/reflex@latest

WORKDIR /src

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
