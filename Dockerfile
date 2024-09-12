# syntax=docker/dockerfile:1

FROM golang:1.23.1

WORKDIR /app

COPY go.mod /app

RUN go mod download

COPY . ./

RUN go build -o ./yumfoods ./cmd/.

EXPOSE 4000

CMD ["./yumfoods"]
