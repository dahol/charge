
FROM golang:1.20-alpine

WORKDIR /app

COPY . .

RUN go mod init generator
RUN go mod tidy
RUN go build -o generator .

CMD ["./generator"]

