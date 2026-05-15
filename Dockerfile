FROM golang:1.25.6 AS builder
WORKDIR /app
COPY . .
RUN go build -o david-workout-app ./cmd/api/

FROM debian:bookworm-slim
WORKDIR /app
COPY --from=builder /app/david-workout-app .
EXPOSE 8080
CMD ["./david-workout-app"]
