
FROM golang:1.25-alpine AS builder

WORKDIR /app


RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .


RUN go build -o main ./cmd

RUN go build -o migrate ./cmd/migrate

FROM alpine:3.22

WORKDIR /root/


COPY --from=builder /app/main .
COPY --from=builder /app/migrate .
COPY --from=builder /app/cmd/migrate/migrations ./cmd/migrate/migrations/
COPY --from=builder /app/log/ ./log

EXPOSE 8080


CMD ["./main"]