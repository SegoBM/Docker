FROM golang:1.21-alpine

WORKDIR /app

COPY . .

RUN go mod download

# Compilar aplicación
RUN go build -o main .

# Expone el puerto 8080
EXPOSE 8080

CMD ["./main"]