# message_broker

## Тестовое задание на позицию Go разработчика

## Создание базы данных
### Все описание базы данных находится в файле database/psql.db.sql
### После создания базы можно запускать проект

## Запуск прокета:
```go
go run main.go
```

После запуска проекта  произойдет обращение к базе данных, из токорой забираются все пользователи.
Для каждого пользователя запускается goroutine в которой в бесконечном цикле из базы вытаскиваются таски
конкретного пользователя. Таски делятся на два типа: пополнение баланса и его списание. Если баланс 
недостаточен для списания, то таска отменяется.

## Примеры обращения к API

Все примеры описаны в формате OpenApi в файле openapi.yaml

