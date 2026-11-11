

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

Для работы с миграциями [migrate](https://github.com/golang-migrate/migrate)

### Команды для разработки

```bash
# Установка migrate
task install-migrate

# Работа с миграциями
task migrate-create 
task migrate-up    
task migrate-down

# Сборка проекта
task build
task run
task e2e-test
```
#### task migrate-create

Команда выполняет создание `.up.sql` и `.down.sql` файлов миграции в `./db/migrations`.

В команду передается имя файла миграции.

```bash
# Создает:
#  ./db/migrations/YYYYMMDDHHmmss_example_migration.up.sql
#  ./db/migrations/YYYYMMDDHHmmss_example_migration.up.sql
task migrate-create -- example_migration
```

#### task migrate-up

Команда выполняет `.up.sql` миграции.

В команду можно передать количество миграций, которое необходимо выполнить. По умолчанию выполняются все миграции.

```bash
# Выполнит все миграции
task migrate-up

# Выполнит 1 миграцию
task migrate-up -- 1
```


#### task migrate-down

Команда выполняет `.down.sql` миграции (откатывает изменения соответствующих `.up.sql` миграций).

В команду можно передать количество миграций, которое необходимо выполнить. По умолчанию выполняются все миграции, но требуется подтверждение.

```bash
# Выполнит все миграции
task migrate-down

# Выполнит 1 миграцию
task migrate-down -- 1
```

#### task run

Команда выполняет запуск исполняемого файла todolist в `.build`.

В команду можно передать параметр `env` - окружение. Допустимые значения `default`, `prod`, `test`, `docker`.

В зависимости от окружения настройки берутся из соответствующего конфигурационного файла: `config.yaml`, `config.prod.yaml`, `config.test.yaml`, `config.docker.yaml`. По умолчанию установлено значение `default`.

```bash
# По умолчанию env=default
task run

# Тоже самое что task run
task run env=default 
```

### Запуск проекта

```bash
task deploy
cd build
./server
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
