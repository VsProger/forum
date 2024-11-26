# Используем официальный образ Golang как базовый
FROM golang:1.22-alpine

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем go.mod и go.sum файлы в рабочую директорию
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod tidy

# Копируем весь проект в рабочую директорию
COPY . .

# Собираем приложение
RUN go build -o main ./cmd/web

# Открываем порт, на котором будет работать приложение
EXPOSE 8080

# Запускаем приложение
CMD ["./main"]
