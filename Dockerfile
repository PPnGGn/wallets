FROM golang:1.25 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o wallets-app ./cmd/main.go
FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/wallets-app .
EXPOSE 8080
CMD ["/app/wallets-app"]