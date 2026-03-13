# Build stage
FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o main cmd/api/main.go

# Final stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .

# Default Environment Variables
ENV GIN_MODE=debug
ENV PORT=4000
ENV DB_HOST=db
ENV DB_USER=postgres
ENV DB_PASSWORD=password_kamu
ENV DB_NAME=postgres
ENV DB_PORT=5432
ENV DB_SSLMODE=disable
ENV JWT_SECRET=YOUR-JWT-SECRET
ENV JWT_EXPIRES=enable
ENV JWT_EXPIRES_IN=24

EXPOSE 4000
CMD ["./main"]