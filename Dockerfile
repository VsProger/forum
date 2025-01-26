# Используем официальный образ Golang как базовый
FROM golang:1.22-alpine

# Устанавливаем необходимые зависимости для CGO
RUN apk add --no-cache gcc musl-dev

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем go.mod и go.sum
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod tidy

# Копируем весь проект
COPY . .

# Собираем приложение с CGO_ENABLED=1
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o main ./cmd/web

# Открываем порт приложения
EXPOSE 8081

# Запускаем приложение
CMD ["./main"]
