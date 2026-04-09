# URL Shortener Service

HTTP-сервис на Go для создания и получения сокращённых ссылок.

Поддерживает два типа хранилища:

* `memory` — потокобезопасное in-memory хранилище
* `postgres` — постоянное хранение в PostgreSQL

Тип хранилища выбирается через переменную окружения при запуске.

---

## Возможности

* HTTP API
* `POST /shorten`
* `GET /{short}`
* поддержка `memory` и `postgres`
* выбор storage через ENV
* thread-safe in-memory repository
* PostgreSQL persistence
* unit и concurrent tests
* Docker + docker-compose
* SQL migrations

---

### Компоненты

* **storage**

  * интерфейс `Repository`
  * `MemoryRepository`
  * `PostgresRepository`
  * поиск в обе стороны: `short -> long`, `long -> short`
  * потокобезопасность через `sync.RWMutex` для in-memory и `pgx.Pool` для postrges
  
* **generator**

  * интерфейс `Generator`
  * `MemoryGenerator` на `atomic.Uint64`
  * `PostgresGenerator` на sequence `short_url_seq`
  
* **codec**

  * кодирование `uint64 -> string`
  * fixed length = `10`
  * base-63 alphabet
  
* **service**

  * нормализация URL
  * возврат существующего short URL для уже сохранённого original URL
  * retry логика генерации short url
  * получение original URL
  
* **handlers**

  * `POST /shorten`
  * `GET /{short}`
  * JSON request/response
  * redirect через `302 Found`
  
* **app**

  * выбор storage по `STORAGE`
  
  ---
  
  ## Архитектура

Client -> HTTP Handlers -> Service -> Repository / Generator -> Memory | PostgreSQL

---

## HTTP API

Сервис принимает HTTP-запросы.

### POST `/shorten`

Создаёт сокращённую ссылку.

#### Request

```json
{"url":"https://finance.ozon.ru"}
```

#### Response

```json
{"shortUrl":"http://localhost:8080/aaaaaaaaab"}
```

#### Status codes

* `200 OK`
* `400 Bad Request`
* `500 Internal Server Error`

---

### GET `/{short}`

Возвращает redirect на original URL.

#### Status codes

* `302 Found`
* `400 Bad Request`
* `404 Not Found`
* `500 Internal Server Error`

---

## Генерация short URL

Генерация ссылки состоит из двух шагов:

1. получение уникального `uint64 id`
2. кодирование `id` в строку длиной 10 символов

### Генерация ID

* `memory` → `atomic.Uint64`
* `postgres` → `SELECT nextval('short_url_seq')`

### Кодирование

Используемый алфавит:

```text
aabcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_
```

`id` переводится в систему счисления по основанию `63` и записывается в буфер фиксированной длины `10`.

Гарантируется:

* длина всегда `10`
* одинаковый `id` даёт одинаковый short url
* разные `id` дают разные short url
* используются только разрешённые символы

---

## Запуск

Приложение поддерживает запуск:

* локально
* через Docker Compose
* в режиме `memory`
* в режиме `postgres`

---

### 1) Подготовка `.env`

### Memory mode

```bash
STORAGE=memory
SERV_PORT
```

### Postgres mode

```bash
STORAGE=postgres
SERV_PORT

POSTGRES_DB
POSTGRES_USER
POSTGRES_PASSWORD
DB_PORT
DATABASE_URL
```

### 2) Сборка и запуск через Docker Compose

```bash
make docker-build
make docker-up
```

После запуска сервис доступен на:

```text
http://localhost:8080
```

### 3) Проверка API

```bash
curl -X POST http://localhost:8080/shorten \
  -H "Content-Type: application/json" \
  -d '{"url":"https://finance.ozon.ru"}'
```

---

## Тестирование

### Запустить тесты

```bash
test-integration
```

Поднимается тестовая бд, после тестов - удаляется 

Покрываются:

* codec
* generator
* repository
* service
* handlers
* concurrent scenarios

---

## Поведение под нагрузкой

* В режиме `memory` генерация ID использует `atomic.Uint64`, что обеспечивает потокобезопасное получение уникальных ID при конкурентных запросах.
* `MemoryRepository` использует `sync.RWMutex`: чтения выполняются параллельно, запись эксклюзивна.
* В режиме `postgres` уникальность ID обеспечивается sequence `short_url_seq`, которая корректно работает при параллельном доступе из нескольких goroutine и даже нескольких инстансов приложения.
* Для short url используется детерминированное кодирование уникального `id`, поэтому при уникальном ID коллизий short url не возникает в пределах диапазона кодируемых значений.
* Конкурентные сценарии покрыты тестами для `generator` и `repository`.

---

## Возможные улучшения

* вынести base URL в конфиг
* добавить metrics
* graceful shutdown
* Swagger
