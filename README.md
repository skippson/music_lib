# Music library

 В данном проектк я реализовал онлайн библиотеку песен, позволяющаю добавлять, удалять или манипулировать данными об песнях.
 Вся информациия хранится в базе данных PostgreSQL

## Установка и запуск

### Шаг 1: Настройка конфигурации

В папке `config` находится файл `config.env`. Вы можете использовать начальную конфигурацию или установить свои собственные настройки в этом файле.

### Шаг 2: Поднятие базы данных

Если вы решили использовать начальную конфигурацию, перед началом работы необходимо поднять базу данных. Для этого выполните следующие команды:

```bash
cd docker
docker-compose up -d
```

### Шаг 3: Запуск сервера

Для старта сервера выполните следующие команды:

```bash
cd src
go run cmd/server/main.go
```

## Тестирование

Формат структуры:
```json
{
    "group": "Group name",
    "song": "Song name",
    "release_date": "years-moths-days",
    "text": "Lyrics"
}
```

---
### Добавление новой песни

```bash
curl -X POST "http://localhost:8888/music" \
     -H "Content-Type: application/json" \
     -d '{"group": "Group Name", "song": "Song Name"}'
```
---
### Получение данных всей библиотеки

```bash
curl -X GET "http://localhost:8888/music"
```
---
### Получение данных библиотеки с пагинацией 

```bash
curl -X GET "http://localhost:8888/music/1/1"
```
---
### Удаление песни

```bash
curl -X  DELETE "http://localhost:8888/music/Group%20Name/Song%20Name"
```
---
### Обновление данных о песне

```bash
curl -X PUT "http://localhost:8888/music/Group%20Name/Song%20Name" \
     -H "Content-Type: application/json" \
     -d '{"group": "New Group Name", "song": "New Song Name", "release_date": "2022-3-3", "text": "I wanna rock"}'
```
---
### Запрос с фильтром

```bash
curl -X GET "http://localhost:8888/music/filter?group=Group%20Name" 
```

### Получение текста песни
---
```bash
curl -X GET "http://localhost:8888/music/Group%20Name/Song%20Name/lyrics"
```
---
### Запрос с фильтром и пагинацией

```bash
curl -X GET "http://localhost:8888/music/filter/1/10?group=Group%20Name" 
```
---
### Получение текста песни с пагинацией

```bash
curl -X GET "http://localhost:8888/music/Group%20Name/Song%20Name/lyrics/1/10"
```

