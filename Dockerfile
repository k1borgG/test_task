FROM golang:latest as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
RUN go mod tidy

COPY . .

WORKDIR /app/cmd/service

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o service .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/cmd/service/service .

EXPOSE 8080

CMD ["./service"]
