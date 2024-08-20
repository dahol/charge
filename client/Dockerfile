FROM golang:1.20-alpine

WORKDIR /app

COPY . .

RUN go mod init client && go mod tidy && go build -o client .
RUN go mod tidy
RUN go build -o client .

CMD ["./client"]

