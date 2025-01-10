# Приложение для хранения данных о музыке 
* Реализует REST API для добавления, получения и редактирования данных о музыке
* Поддерживает пагинацию текста по куплетам
* Хранит данные в PostgreSQL

## Запуск 
Клонировать репозиторий 
```bash
git clone https://github.com/madrat19/Music.git
```
Запустить контейнеры
```bash
docker compose up --build -d
```

## Конфигурация
* Данные для конфигурации хранятся в файле .env 
* Среди них данные для подключения к БД, хост приложения, уровень логирования и адрес API для получаения данных о музыке
* Всё переменные, кроме MUSICINFO можно оставить неизменными

## Music info API
* Приложение реализует mock версию music info API.
* Она может выдать данные только об одной песни: "Roads" группы "Portishead".
* Чтобы задействовать её, нужно в файле .env для переменной MUSICINFO выставить значение "mock".
* В противном случае нужно выставить адрес реального API в той же переменной.

## Документация
* В папке api содержится swagger, описывающий весь api приложения
* Swagger UI будет доступен после запуска приложения по адресу:
```bash
http://localhost:8080/swagger/index.html
```


## Логирование
Приложение поддеживает 3 уровня логирования
* fatal: только критические ошибки
* error: любые ошибки
* info: наиболее полная информация

Настроить уровень можно выставив соответсвующее значение для переменной LOGLEVEL в файле .env

Все логи сохраняются в файле app.log, получить к нему доступ можно следующим образом:
```bash
docker exec -it --user=root app /bin/sh
```
```bash
cat app.log
```
## Тестирование 
Тестировать функционал можно с помощью Swagger UI или набора curl-запросов ниже


При тестировании изнутри контейнера необходимо будет установить в него curl:
```bash
docker exec -it --user=root app /bin/sh
```
```bash
apk --no-cache add curl
```

Добавяем новую песню:
```bash
curl -X POST http://localhost:8080/songs --url-query song=Roads --url-query group=Portishead
```

Получаем список песен:
```bash
curl -X GET http://localhost:8080/songs 
```

Получаем конкретную страницу из спсика песен (по 10 на странице):
```bash
curl --url-query page=1 http://localhost:8080/songs
```

Получаем список песен с фильтрацией по полям:
```bash
curl --url-query releasedate=22.08.1994 --url-query group=Portishead http://localhost:8080/songs
```

Получаем текст песни:
```bash
curl --url-query song=Roads --url-query group=Portishead http://localhost:8080/text
```

Получаем конкретный куплет из текста песни:
```bash
curl --url-query song=Roads --url-query group=Portishead --url-query verse=2 http://localhost:8080/text
```
Изменяем данные о песни:
```bash
curl -X PATCH http://localhost:8080/songs -H "Content-Type: application/json; ; charset=utf-8" -d '{"song": "Roads", "group": "Portishead", "releasedate": "01.01.2024"}'
```

Удаляем песню:
```bash
curl -X DELETE --url-query song=Roads --url-query group=Portishead http://localhost:8080/songs
```

meow



