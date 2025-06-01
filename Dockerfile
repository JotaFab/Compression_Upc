FROM golang:1.24-alpine

WORKDIR /app
COPY . .

# Crear directorio process si no existe y establecer permisos
RUN mkdir -p process && chmod 777 process

RUN go mod tidy
RUN go build -o /app/main .
RUN chmod +x /app/main

# Exponer el puerto 8080

EXPOSE 8080
CMD ["./main"]
