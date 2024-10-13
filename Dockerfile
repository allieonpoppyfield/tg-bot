# Используем образ Go 1.20 на базе Alpine
FROM golang:1.22-alpine

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем go.mod и go.sum из корневой директории проекта
COPY go.mod go.sum ./

# Загружаем все зависимости
RUN go mod download

# Копируем весь исходный код в контейнер
COPY . .

# Собираем приложение из директории cmd, результат сборки помещаем в корень рабочей директории /app
RUN go build -o bot ./cmd/main.go

# Проверяем, что исполняемый файл создан
RUN ls -la /app

# Указываем команду по умолчанию для запуска контейнера
CMD ["./bot"]
