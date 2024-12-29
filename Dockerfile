# Используем официальный образ Go
FROM golang:1.20

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем файлы go.mod и go.sum
COPY go.mod go.sum ./

# Скачиваем зависимости
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение
RUN go build -o main .

# Устанавливаем переменную окружения для JSON-ключа
ENV GOOGLE_APPLICATION_CREDENTIALS /app/credentials.json

# Копируем JSON-ключ в образ (замени credentials.json на имя твоего файла)
COPY credentials.json /app/credentials.json

# Запускаем приложение
CMD ["./main"]