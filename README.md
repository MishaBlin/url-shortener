# URL-shortener

Сервис для сокращения ссылок на Go.

### Запуск приложения через docker

Есть два варианта запуска приложения: 
1. Запуск одного контейнера, который использует свой внутренний in-memory storage в качестве базы данных.
2. Запуск двух контейнеров через docker compose. В этом случае будет использован postgres. 

### Запуск с in-memory storage

Запускаем приложение в контейнере и передаем тип хранилища `storage-type`:
```bash
docker run -p 8080:8080 -e ENV=local mishablin/url-shortener --storage-type=memory
```

### Запуск с хранилищем postgres

Запускаем через `docker-compose`:
```bash
docker-compose up --build
```

*Для работы с базой данных нужен файл* `.env`. Пример файла можно посмотреть в `example.env`.

### Описание запросов:

- `POST` на url `localhost:8080/url`. Сохраняет переданный url в базе данных и возвращает сокращенный url.

Тело запроса:
```json
{
  "url": "https://google.com"
}
```

Ответ:
```json
{
  "alias": "abcd1234"
}
```

- `GET` на url `localhost:8080/url/{alias}`. Редирект на оригинальный url по сокращенному.