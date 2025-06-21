# dispatcher/Dockerfile
FROM golang:1.22-alpine

WORKDIR /app

# Sólo copiamos go.mod y generamos go.sum dentro de la imagen
COPY go.mod ./
RUN go mod tidy

# Copiamos el resto del proyecto
COPY . .

# Entramos a la carpeta dispatcher y compilamos
WORKDIR /app/dispatcher
RUN go build -o dispatcher main.go

EXPOSE 8000
CMD ["./dispatcher"]
