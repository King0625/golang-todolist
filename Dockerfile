FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o ./todolist ./cmd/api

FROM alpine:latest
COPY --from=builder /app/migration /app/migration
COPY --from=builder /app/todolist /app/
WORKDIR /app

EXPOSE 11451

CMD ["./todolist"]