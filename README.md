# Simple GO CRM API
Простое приложение на GO, предназначенное для подтверждения навыков построения архитектуры REST API

## Краткое описание ТЗ
Реализовать приложение позволяющее пользователям бесплатно получать доступ к спортивным площадкам города

### Доступные операции
- Регистрация
- Авторизация JWT
- Работа с файлами
- Работа с адресами (подсказки в заполнении адреса)
- Работа со списками площадок и команд

### Роль пользователя
- Создание/Редактирование/Удаление карточки команды
- TODO Приглашение в команду других пользователей
- Бронирование площадки для игры

### Роль Администратора
- Создание/Редактирование/Удаление карточки площадки
- Создание/Редактирование/Удаление карточки команды
- TODO Редактирование/Удаление карточки пользователя

## Консольные команды

#### Получить все зависимости
```bash
go get .
```
#### Собрать приложение
```bash
go build -v .
```
#### Запустить приложение
```bash
go run .
```
#### Запустить контейнер Postgres
```bash
docker compose up -d go_db
```
#### Собрать контейнер API
```bash
docker compose build
```
#### Запустить контейнер API
```bash
docker compose up go-app
```
#### Собрать контейнер API и запустить
```bash
docker-compose --env-file .env.local up --build
```
#### Зайти в контейнер
```bash
docker exec -it goland_api bash
```
#### Список контейнеров
```bash
docker ps -a
```
#### Список образов
```bash
docker images
```

### Миграции

#### Создать файл миграции
```bash
migrate create -ext sql -dir db/migration -seq create_user_role
```
#### Выполнить
```bash
migrate -path db/migration -database "postgresql://postgres:postgres@localhost:5432/postgres?sslmode=disable" -verbose up
```

### Документация OpenAPI

#### Сгенирировать Swagger
```bash
swag init
```