# Получить все зависимости
`go get .`

# Собрать приложение
`go build -v .`

# Запустить приложение
`go run .`

# Запустить контейнер Postgres
`docker compose up -d go_db`

# Собрать контейнер API
`docker compose build`

# Запустить контейнер API
`docker compose up go-app`

# Собрать контейнер API и запустить
`docker compose up go-app --build`

# Зайти в контейнер
docker exec -it goland_api bash

# Список контейнеров
`docker ps -a`

# Список образов
`docker images`

# Миграции
## Создать файл миграции
migrate create -ext sql -dir db/migration -seq create_user_role

## Выполнить
migrate -path db/migration -database "postgresql://postgres:postgres@localhost:5432/postgres?sslmode=disable" -verbose up

## Сгенирировать Swagger
swag init