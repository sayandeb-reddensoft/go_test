FROM golang:1.24-alpine AS builder

WORKDIR /app/auth
COPY . .
RUN go mod tidy
RUN go build -o app-auth ./main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/auth/app-auth ./
COPY --from=builder /app/auth/templates ./templates

EXPOSE 8080
CMD ["./app-auth"]
