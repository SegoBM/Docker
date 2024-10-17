# Usa una imagen base de Go
FROM golang:1.21-alpine

# Establece el directorio de trabajo dentro del contenedor
WORKDIR /app

# Copia los archivos del proyecto al contenedor
COPY . .

# Descarga las dependencias del proyecto
RUN go mod download

# Compila la aplicación Go
RUN go build -o main .

# Expone el puerto en el que corre la API
EXPOSE 8080

# Comando para ejecutar la API
CMD ["./main"]