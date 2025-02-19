# URL-shortener

Сервис для сокращения ссылок на Go.

## Запуск приложения через docker

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