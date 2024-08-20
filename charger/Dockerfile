FROM golang:1.20-alpine

WORKDIR /app

COPY . .

RUN go mod init charger
RUN go mod tidy
RUN go build -o charger .

CMD ["./charger"]

