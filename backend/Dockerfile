FROM golang:1.20-alpine

WORKDIR /app

COPY . .

RUN go mod init backend
RUN go mod tidy
RUN go build -o backend .

CMD ["./backend"]

