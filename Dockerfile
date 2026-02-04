FROM golang:1.24.6-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git gcc musl-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o subtrack-service ./cmd/service
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o subtrack-cli ./cmd/cli

FROM alpine:latest

RUN apk --no-cache add ca-certificates sqlite-libs

WORKDIR /root/

COPY --from=builder /app/subtrack-service .
COPY --from=builder /app/subtrack-cli .

CMD ["./subtrack-service"]
