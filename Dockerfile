# ---------- BUILD STAGE ----------
FROM golang:1.22 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app ./cmd/server/main.go


# ---------- RUN STAGE ----------
FROM debian:bookworm-slim

WORKDIR /root/

# install required tools
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

# copy binary
COPY --from=builder /app/app .

# copy migrations
COPY --from=builder /app/migrations ./migrations

EXPOSE 8080

CMD ["./app"]