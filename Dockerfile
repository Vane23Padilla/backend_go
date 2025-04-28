# Usa la imagen oficial de Go
FROM golang:1.21

# Crea el directorio de trabajo en el contenedor
WORKDIR /app

# Copia los módulos y descarga dependencias
COPY go.mod go.sum ./
RUN go mod download

# Copia todo el código fuente
COPY . .

# Compila la aplicación
RUN go build -o server_estudiantes

# Expone el puerto en el contenedor
EXPOSE 8080

# Comando que corre tu aplicación
CMD ["./server_estudiantes"]

