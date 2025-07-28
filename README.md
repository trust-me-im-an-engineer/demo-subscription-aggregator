# Сервис агрегации подписок

REST API сервис для управления подписками пользователей с возможностью CRUDL-операций и расчета стоимости подписок.

## Возможности

- **Управление подписками**:
    - Создание, просмотр, обновление и удаление подписок
    - Получение списка всех подписок
- **Расчет стоимости**:
    - Расчет общей стоимости подписок с фильтрацией по:
        - ID пользователя
        - Названию сервиса (частичное совпадение)
        - Диапазону дат
- **База данных**:
    - PostgreSQL в качестве хранилища
    - Автоматическое применение миграций при запуске
- **Документация**:
    - Полная Swagger/OpenAPI документация
- **Развертывание**:
    - Поддержка Docker и Docker Compose

## API Эндпоинты

| Метод | Эндпоинт                     | Описание                          |
|--------|------------------------------|--------------------------------------|
| POST   | /subscriptions               | Создать новую подписку            |
| GET    | /subscriptions               | Получить все подписки                |
| GET    | /subscriptions/{id}          | Получить подписку по ID               |
| PUT    | /subscriptions/{id}          | Обновить подписку                |
| DELETE | /subscriptions/{id}          | Удалить подписку                |
| GET    | /subscriptions/total-cost    | Рассчитать общую стоимость подписок   |
| GET    | /swagger/                    | Просмотр Swagger документации           |


## Начало работы

### Требования
- Docker и Docker Compose
- Go 1.24+ (для локальной разработки)

### Запуск через Docker Compose

1. Создайте файл `.env` в корне проекта:
   ```
   DB_PASSWORD=ваш_пароль_бд
   POSTGRES_PASSWORD=ваш_postgres_пароль
   ```

2. Запустите сервисы:
   ```bash
   docker-compose up
   ```

Сервис будет доступен по адресу `http://localhost:8080`

### Локальная разработка

1. Установите зависимости:
   ```bash
   go mod download
   ```

2. Настройте переменные окружения (или создайте файл `.env`):
   ```bash
   export APP_ADDRESS=localhost:8000
   export APP_LOG_LEVEL=INFO
   export DB_HOST=localhost
   export DB_NAME=postgres
   export DB_PASSWORD=secret-password
   export DB_PORT=5432
   export DB_USER=postgres
   export POSTGRES_PASSWORD=secret-password
   ```

3. Запустите сервис:
   ```bash
   go run cmd/server/main.go
   ```

## Тестирование

Проект включает файл `test/test.http` с простейшими тестами API-запросов, которые можно использовать с HTTP-клиентами в IDE (например VS Code или JetBrains).

## Конфигурация

Настройки управляются через переменные окружения:

| Переменная           | Описание                   | По умолчанию |
|----------------------|----------------------------|--------------|
| APP_ADDRESS          | Адрес сервера              | 0.0.0.0:8080 |
| APP_SHUTDOWN_TIMEOUT | Таймаут graceful shutdown  | 10s          |
| DB_HOST              | Хост PostgreSQL            | -            |
| DB_PORT              | Порт PostgreSQL            | -            |
| DB_USER              | Пользователь PostgreSQL    | -            |
| DB_PASSWORD          | Пароль PostgreSQL          | -            |
| DB_NAME              | Имя базы данных PostgreSQL | -            |

## Документация

API документация доступна по адресу `http://localhost:8080/swagger/` при запущенном сервисе.

## Структура проекта

```
.
├── cmd
│   └── server          # Точка входа приложения
├── internal
│   ├── config          # Загрузка конфигурации
│   ├── models          # Модели данных и DTO
│   ├── repository      # Слой работы с БД
│   │   └── postgres    # Реализация для PostgreSQL
│   ├── service         # Бизнес-логика
│   ├── validation      # Валидация запросов
│   └── web             # HTTP обработчики и роутинг
├── migrations          # Миграции базы данных
├── pkg
│   ├── logger          # Настройка логгирования
│   └── monthyear       # Кастомный тип даты
└── test                # Тестовые запросы
```