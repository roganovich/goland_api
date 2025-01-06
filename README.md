# Получить все зависимости
`go get .`

# Запустиить приложение
`go run .`

# Запустиить контейнер Postgres
`docker compose up -d go_db`

# Собрать контейнер API
`docker compose build`

# Запустиить контейнер API
`docker compose up go-app`

# Список контейнеров
`docker ps -a`

# Список образов
`docker images`

# Миграции
## Создать файл миграции
migrate create -ext sql -dir db/migration -seq create_user

## Выполнить
migrate -path db/migration -database "postgresql://postgres:postgres@localhost:5432/postgres?sslmode=disable" -verbose up