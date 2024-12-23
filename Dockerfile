FROM golang:1.22.1-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o app .

EXPOSE 6000

CMD ["./app"]
