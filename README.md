# Сервис для поздравлений с днем рождения
## Описание сервиса
Сервис предназначен для управления поздравлениями с днем рождения. Пользователи могут регистрироваться, добавлять сотрудников, подписываться на уведомления о днях рождения, а также удалять подписки. Сервис реализован с использованием фреймворка Fiber и базы данных PostgreSQL.

## Запуск проекта
Для запуска проекта вам понадобится Docker. Создайте файл docker-compose.yml с настройками базы данных и запустите контейнеры:

Поднитие контейнера:
```bash
docker-compose up -d
```
После поднятия постгреса нужно запустить миграции испольщуя migrate:
```bash
migrate -database "postgres://test:test@localhost:1111/birthdays?sslmode=disable" -path "./db/migrations" up  
```
После запуска контейнера, запустите сервер Go-приложения:
```bash
go run main.go
```
Теперь сервер должен быть доступен по адресу http://localhost:1234.

Регистрация:
```bash
URL: /register
Метод: POST
Описание: Регистрация нового пользователя.

curl -X POST http://localhost:1234/register \
-H "Content-Type: application/json" \
-d '{
  "name": "John Doe",
  "email": "john.doe@example.com",
  "password": "password",
  "birthday": "01.01.1990"
}'
```
Авторизация:
```bash
URL: /login
Метод: POST
Описание: Вход в систему.

curl -X POST http://localhost:1234/login \
-H "Content-Type: application/json" \
-d '{
  "email": "john.doe@example.com",
  "password": "password"
}'
```
Выход:
```bash
URL: /logout
Метод: POST
Описание: Выход из системы.

curl -X POST http://localhost:1234/logout
```

### Эти эндпоинты требуют авторизации. Не забудьте добавить заголовок Authorization с JWT-токеном
Получить всех сотрудников:
```bash
URL: /api/employees
Метод: GET
Описание: Получение списка всех сотрудников.

curl -X GET http://localhost:1234/api/employees \
--cookie "token=<your-jwt-token>"
```
Удалить сотрудника:
```bash
URL: /api/employees
Метод: DELETE
Описание: Удаление сотрудника по email.

curl -X DELETE http://localhost:1234/api/employees \
--cookie "token=<your-jwt-token>" \
-H "Content-Type: application/json" \
-d '{
  "email": "john.doe@example.com"
}'
```
Получить все подписки пользователя:
``` bash
URL: /api/subs
Метод: GET
Описание: Получение всех подписок текущего пользователя.

curl -X GET http://localhost:1234/api/subs \
--cookie "token=<your-jwt-token>"
```
Создать подписку:
```bash
URL: /api/subs
Метод: POST
Описание: Создание новой подписки.

curl -X POST http://localhost:1234/api/subs \
--cookie "token=<your-jwt-token>" \
-H "Content-Type: application/json" \
-d '{
  "id": 2,
  "time": "12:00"
}'
```
Удалить подписку:
```bash
URL: /api/subs
Метод: DELETE
Описание: Удаление подписки.

curl -X DELETE http://localhost:1234/api/subs \
--cookie "token=<your-jwt-token>" \
-H "Content-Type: application/json" \
-d '{
  "id": 2
}'
```
