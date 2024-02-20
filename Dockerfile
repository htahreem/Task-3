FROM golang:latest

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o t3 .

EXPOSE 3000

CMD ["./t3"]
