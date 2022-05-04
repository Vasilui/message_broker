# message_broker

## Тестовое задание на позицию Go developer

## Описание задачи
Stack: go, fiber, postgresql, insomnia
Нужно реализовать систему транзакций.
Как происходит транзакция:
Идет запрос на сервер от клиента, по клиенту выстраивается очередь на вывод.
Важно : 
1) у каждого клиента есть своя очередь; 
2) при нехватке денег, нужно блокировать запрос

Что нужно реализовавть :
бд на postgresql, где будет схема с клиентами и их балансами
сервер, которые проверяет все условия(хватает ли денег, если сервер упадет, то история, которая идет на вывод не должна пропасть) и делает изменение баланса(на + или -)


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

