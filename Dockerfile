# Build stage
FROM golang:1.22 AS builder

WORKDIR /go/src/zomato-backend-assignment

COPY go.mod go.sum ./
RUN go mod download

COPY . /go/src/zomato-backend-assignment

RUN go build -o main ./cmd/server

# Run stage
FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /go/src/zomato-backend-assignment/main .
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*


EXPOSE 8080

CMD ["./main"]