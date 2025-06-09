FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o worker-matrix main.go

FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/worker-matrix .
ENV PORT=8082
EXPOSE 8082
CMD ["./worker-matrix"]
