# Получить все зависимости
`go get ./cmd/server/main.go`

# Собрать приложение
`go build -v ./cmd/server/main.go`

# Запустить приложение
`go run ./cmd/server/main.go`

# Запустить контейнер Postgres
`docker compose up -d go_db`

# Собрать контейнер API
`docker compose build`

# Запустить контейнер API
`docker compose up go-app`

# Собрать контейнер API и запустить
`docker compose up go-app --build`

# Список контейнеров
`docker ps -a`

# Список образов
`docker images`

# Миграции
## Создать файл миграции
migrate create -ext sql -dir db/migration -seq create_media

## Выполнить
migrate -path db/migration -database "postgresql://postgres:postgres@localhost:5432/postgres?sslmode=disable" -verbose up