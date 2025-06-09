# Banking Service REST API in Go

This project is a REST API for a banking service implemented in Go. It follows a clean architecture with separate layers for models, repositories, services, handlers, and middleware. The API supports user registration, authentication, bank account management, card operations, transfers, credit operations, financial analytics, and integrations with external services (Central Bank of Russia for the key rate and SMTP for email notifications).

## Prerequisites

- Go 1.23+
- PostgreSQL 17 with the `pgcrypto` extension enabled
- PGP key pair for card data encryption
- SMTP email account for sending notifications

## Directory Structure

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
│   └── auth_handler.go
│   └── account_handler.go
│   └── card_handler.go
│   └── transfer_handler.go
│   └── analytics_handler.go
│   └── credit_handler.go
├── middleware/
│   └── auth.go
└── utils/
    ├── luhn.go
    ├── crypto.go
    └── soap_client.go (TODO)
```

## Environment Variables

Create a `.env` file or export the following environment variables:

```bash
export DATABASE_URL="postgres://username:password@localhost:5432/banking_db?sslmode=disable"
export JWT_SECRET="your_jwt_secret"
export PGP_PRIVATE_KEY_PATH="/path/to/pgp_private_key.asc"
export PGP_PUBLIC_KEY_PATH="/path/to/pgp_public_key.asc"
export SMTP_HOST="smtp.example.com"
export SMTP_PORT="587"
export SMTP_USER="your_email@example.com"
export SMTP_PASS="your_email_password"
export PORT="8080"
```

## Database Setup

1. Create a PostgreSQL database (e.g., `banking_db`).
2. Enable the `pgcrypto` extension:

```sql
CREATE EXTENSION IF NOT EXISTS pgcrypto;
```

3. Create the required tables:

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

## Installation and Running

1. **Clone or download the project** and navigate into the directory:

    ```bash
    unzip banking_service_project.zip
    cd banking_service_project
    ```

2. **Install dependencies**:

    ```bash
    go mod tidy
    ```

3. **Set environment variables** as described above (or use a .env loader).

4. **Run the server**:

    ```bash
    go run main.go
    ```

   The server will start on `http://localhost:8080` (or the port you specified).

## API Endpoints

### Public

- `POST /register` — Register a new user. JSON body: `{"username": "...", "email": "...", "password": "..."}`.
- `POST /login` — Authenticate. JSON body: `{"email": "...", "password": "..."}`. Returns a JWT token.

### Protected (Require `Authorization: Bearer <token>` header)

- `POST /accounts` — Create a new bank account.
- `GET /accounts` — Retrieve all accounts for the authenticated user.
- `POST /cards?account_id={account_id}` — Generate a new virtual card for the given account.
- `GET /cards?account_id={account_id}` — Retrieve all cards for the given account.
- `POST /transfer` — Transfer funds. JSON body: `{"from_account_id": ..., "to_account_id": ..., "amount": ...}`.
- `GET /analytics` — Get monthly income/expense analytics.
- `GET /credits/{creditId}/schedule` — Get payment schedule for a specific credit.
- `GET /accounts/{accountId}/predict?days={n}` — Predict account balance after N days.
- `POST /credits/apply` — Apply for a credit. JSON body: `{"account_id": ..., "principal": ..., "annual_rate": ..., "term_months": ...}`.

## Security

- **Passwords** are hashed using `bcrypt`.
- **JWT** is used for authentication (`golang-jwt/jwt/v5`).
- **Card numbers** and **CVV** are intended to be encrypted with PGP; placeholder functions exist in `utils`.
- **HMAC** is used for data integrity checks.

## Notes

- The current implementation provides function stubs and TODO comments for SOAP integration with the Central Bank of Russia (CBR) and PGP encryption/decryption. You will need to fill these in with appropriate libraries.
- Error handling can be improved further, especially around transactions and rollbacks.
- Scheduling of credit payments can be implemented using a cron job or any scheduler library.
- Logging can be enhanced with `logrus` fields for request tracing.
- Be sure to secure your environment variables and production settings when deploying.

---

By following these instructions, you should be able to set up, run, and test the banking service API locally. Good luck!
