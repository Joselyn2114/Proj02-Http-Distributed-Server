FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o worker-pi main.go

FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/worker-pi .
ENV PORT=8081
EXPOSE 8081
CMD ["./worker-pi"]
