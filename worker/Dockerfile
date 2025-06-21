# worker/Dockerfile
FROM golang:1.22-alpine

WORKDIR /app

# Sólo copiamos go.mod y generamos go.sum dentro de la imagen
COPY go.mod ./
RUN go mod tidy

# Copiamos todo el proyecto (incluye worker/)
COPY . .

# Entramos al subdirectorio worker y compilamos
WORKDIR /app/worker
RUN go build -o worker .

EXPOSE 8080
CMD ["./worker"]
