

# Демо проект todo-api

Этот репозиторий является демонстрацией моих навыков работы с Go и другими технологиями.

Технологии которые используются в проекте:
- Go
- Docker
- sqlite3
- migrate
- task

## TODO-лист

REST API для управления задачами в TODO-листе

## Возможности

- Базовые CRUD-операции с задачами:
    - получить все задачи
    - получить задачу по id
    - создать новую задачу

## Быстрый старт

Для запуска понадобится [Taskfile](https://taskfile.dev)

### Команды для разработки

Подробнее с командами можно ознакомиться в Taskfile.yaml

```bash
# Миграции
task migrate-create 
task migrate-up    
task migrate-down

# Сборка
task copy-db  
task build

# Docker
task docker-build
task docker-run
```

### Запуск проекта

```bash
# Обновляем базу данных до последней версии
task migrate-up

# Создаем ссылку на базу данных в папке с исполняемым файлом
task copy-db

# Собираем сервис в исполняемый файл
task build
    
# Запускаем сервис
build/api
```

## API

### Задачи

#### Получить все задачи

```bash
curl http://localhost:8080/tasks
```

**Ответ**: `200 OK` с массивом задач в JSON

#### Получить задачу по id

```bash
curl http://localhost:8080/tasks/1
```

**Ответ**: `200 OK` с задачей в JSON

#### Создать задачу

```bash
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{"title":"Новая задача", "status":"статус задачи"}'
```

**Ответ**: `201 Created` с созданной задачей в JSON
