## Требования

* Go 1.23+
* PostgreSQL 17 с включённым расширением `pgcrypto`
* Пара PGP-ключей для шифрования данных карт
* SMTP-аккаунт для отправки уведомлений по электронной почте

## Структура проекта

```
banking_service_project/
├── go.mod
├── go.sum
├── main.go
├── README.md
├── models/
│   ├── user.go
│   ├── account.go
│   ├── card.go
│   ├── transaction.go
│   ├── credit.go
│   └── payment_schedule.go
├── repositories/
│   ├── user_repository.go
│   ├── account_repository.go
│   ├── card_repository.go
│   ├── transaction_repository.go
│   ├── credit_repository.go
│   └── payment_schedule_repository.go
├── services/
│   ├── auth_service.go
│   ├── account_service.go
│   ├── card_service.go
│   ├── transfer_service.go
│   ├── credit_service.go
│   ├── analytics_service.go
│   └── external_service.go
├── handlers/
│   ├── auth_handler.go
│   ├── account_handler.go
│   ├── card_handler.go
│   ├── transfer_handler.go
│   ├── analytics_handler.go
│   └── credit_handler.go
├── middleware/
│   └── auth.go
└── utils/
    ├── luhn.go
    ├── crypto.go
    └── soap_client.go (TODO)
```

* **go.mod** и **go.sum** — файлы зависимостей проекта.
* **main.go** — точка входа: подключение к базе, инициализация репозиториев, сервисов, обработчиков и запуск HTTP-сервера.
* **models/** — структуры данных (Users, Accounts, Cards, Transactions, Credits, PaymentSchedules) с JSON-тегами.
* **repositories/** — слой доступа к PostgreSQL: параметризованные SQL-запросы для создания, получения и обновления сущностей.
* **services/** — бизнес-логика: регистрация/логин (bcrypt + JWT), управление счетами, переводы, генерация карт по алгоритму Луна, расчёт аннуитета для кредитов, заготовки для интеграции с ЦБ РФ (SOAP) и SMTP (Gomail), а также аналитика.
* **handlers/** — HTTP-обработчики: парсинг JSON из запросов, валидация, вызов сервисов и возвращение JSON-ответов с корректными статусами.
* **middleware/** — JWT-аутентификация: проверка токена в заголовке `Authorization`, извлечение `userID` в контекст запроса.
* **utils/** — вспомогательные функции: генерация номера карты по алгоритму Луна, генерация CVV, заготовки для PGP-шифрования и SOAP-клиента (пока помечены как `TODO`).

## Переменные окружения


```bash
export DATABASE_URL="postgres://username:password@localhost:5432/banking_db?sslmode=disable"
export JWT_SECRET="ваш_секрет_для_JWT"
export PGP_PRIVATE_KEY_PATH="/path/to/pgp_private_key.asc"
export PGP_PUBLIC_KEY_PATH="/path/to/pgp_public_key.asc"
export SMTP_HOST="smtp.example.com"
export SMTP_PORT="587"
export SMTP_USER="your_email@example.com"
export SMTP_PASS="your_email_password"
export PORT="8080"
```

* **DATABASE\_URL** — строка подключения к базе PostgreSQL.
* **JWT\_SECRET** — секрет для подписи JWT-токенов.
* **PGP\_PRIVATE\_KEY\_PATH** и **PGP\_PUBLIC\_KEY\_PATH** — пути до PGP-ключей (используются для шифрования/дешифрования данных карт).
* **SMTP\_HOST**, **SMTP\_PORT**, **SMTP\_USER**, **SMTP\_PASS** — настройки SMTP-сервера для отправки email-уведомлений.
* **PORT** — порт, на котором будет запущен HTTP-сервер (по умолчанию 8080).

## Настройка базы данных


   ```sql
   CREATE DATABASE banking_db;
   \c banking_db
   ```

   ```sql
   CREATE EXTENSION IF NOT EXISTS pgcrypto;
   ```

   ```sql
   CREATE TABLE users (
       id SERIAL PRIMARY KEY,
       username VARCHAR(50) UNIQUE NOT NULL,
       email VARCHAR(100) UNIQUE NOT NULL,
       password TEXT NOT NULL,
       created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW()
   );

   CREATE TABLE accounts (
       id SERIAL PRIMARY KEY,
       user_id INTEGER REFERENCES users(id),
       balance NUMERIC(20,2) NOT NULL DEFAULT 0,
       created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW()
   );

   CREATE TABLE cards (
       id SERIAL PRIMARY KEY,
       account_id INTEGER REFERENCES accounts(id),
       encrypted_number TEXT NOT NULL,
       encrypted_cvv TEXT NOT NULL,
       expires_at TIMESTAMP WITHOUT TIME ZONE NOT NULL,
       created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW()
   );

   CREATE TABLE transactions (
       id SERIAL PRIMARY KEY,
       from_account_id INTEGER REFERENCES accounts(id),
       to_account_id INTEGER REFERENCES accounts(id),
       amount NUMERIC(20,2) NOT NULL,
       type VARCHAR(20) NOT NULL,
       created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW()
   );

   CREATE TABLE credits (
       id SERIAL PRIMARY KEY,
       account_id INTEGER REFERENCES accounts(id),
       principal NUMERIC(20,2) NOT NULL,
       interest_rate NUMERIC(5,2) NOT NULL,
       term_months INTEGER NOT NULL,
       created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW()
   );

   CREATE TABLE payment_schedules (
       id SERIAL PRIMARY KEY,
       credit_id INTEGER REFERENCES credits(id),
       due_date TIMESTAMP WITHOUT TIME ZONE NOT NULL,
       amount NUMERIC(20,2) NOT NULL,
       paid BOOLEAN NOT NULL DEFAULT FALSE
   );
   ```

## Доступные эндпоинты

### Публичные

* `POST /register` — регистрация нового пользователя.
  **Тело запроса (JSON):**

  ```json
  {
    "username": "ivan_petrov",
    "email": "ivan@example.com",
    "password": "пароль123"
  }
  ```

  **Возвращает:** созданного пользователя (без поля `password`) и статус `201 Created`.

* `POST /login` — аутентификация.
  **Тело запроса (JSON):**

  ```json
  {
    "email": "ivan@example.com",
    "password": "пароль123"
  }
  ```

  **Возвращает:** `{ "token": "<JWT-токен>" }`. Срок жизни токена — 24 часа.

### Защищённые (требуют заголовок `Authorization: Bearer <token>`)

* `POST /accounts` — создать новый банковский счёт.
* `GET /accounts` — получить все счета аутентифицированного пользователя.
* `POST /cards?account_id={account_id}` — сгенерировать виртуальную карту для указанного счёта.
* `GET /cards?account_id={account_id}` — получить все карты по указанному счёту.
* `POST /transfer` — совершить перевод.
  **Тело запроса (JSON):**

  ```json
  {
    "from_account_id": 1,
    "to_account_id": 2,
    "amount": 500.00
  }
  ```
* `GET /analytics` — получить аналитику за текущий месяц (доходы/расходы).
* `GET /credits/{creditId}/schedule` — получить график платежей по кредиту с `creditId`.
* `GET /accounts/{accountId}/predict?days={n}` — прогноз баланса на `n` дней вперёд.
* `POST /credits/apply` — подать заявку на кредит.
  **Тело запроса (JSON):**

  ```json
  {
    "account_id": 1,
    "principal": 10000,
    "annual_rate": 12,
    "term_months": 12
  }
  ```

---

