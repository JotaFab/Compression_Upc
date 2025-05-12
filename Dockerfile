FROM golang:1.24-alpine

WORKDIR /app
COPY . .

# Crear directorio process si no existe y establecer permisos
RUN mkdir -p process && chmod 777 process

EXPOSE 8080
CMD ["go", "run", "."]
