# Notification Platform

![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge\&logo=go\&logoColor=white)
![Gin](https://img.shields.io/badge/Gin-008ECF?style=for-the-badge\&logo=gin\&logoColor=white)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-4169E1?style=for-the-badge\&logo=postgresql\&logoColor=white)
![Apache Kafka](https://img.shields.io/badge/Apache%20Kafka-231F20?style=for-the-badge\&logo=apachekafka\&logoColor=white)
![Redis](https://img.shields.io/badge/Redis-DC382D?style=for-the-badge\&logo=redis\&logoColor=white)

An event-driven notification platform built with Go and microservices. Clients interact with a centralized API Gateway, while authentication, notification management, and email delivery are handled by independent services.

The platform uses **Kafka for asynchronous event-driven communication**, **PostgreSQL with database-per-service architecture**, **Redis for distributed rate limiting**, and **SMTP for email delivery**.

## Architecture

```mermaid
flowchart LR
    Client[Client]

    Gateway[API Gateway]
    Auth[Auth Service]
    Notification[Notification Service]
    Worker[Email Worker]

    AuthDB[(Auth PostgreSQL)]
    NotificationDB[(Notification PostgreSQL)]
    Kafka[(Apache Kafka)]
    Redis[(Redis)]
    SMTP[SMTP Server]

    Client --> Gateway

    Gateway -->|/auth/*| Auth
    Gateway -->|/notifications/*| Notification
    Gateway --> Redis

    Auth --> AuthDB
    Notification --> NotificationDB

    Notification -->|notification.created| Kafka
    Kafka --> Worker

    Worker --> SMTP

    Worker -->|notification.sent| Kafka
    Worker -->|notification.failed| Kafka

    Kafka --> Notification
```

## Features

* Microservice-based architecture
* Centralized API Gateway
* JWT authentication
* Secure password hashing with bcrypt
* PostgreSQL database per service
* Apache Kafka event-driven communication
* Asynchronous email delivery
* SMTP integration
* Automatic email retry
* Notification delivery status tracking
* Redis-based per-user rate limiting
* Reverse proxy routing
* Request validation
* End-to-end API integration tests

## Tech Stack

| Technology   | Purpose                    |
| ------------ | -------------------------- |
| Go           | Backend services           |
| Gin          | HTTP framework             |
| PostgreSQL   | Persistent storage         |
| pgx          | PostgreSQL driver          |
| Apache Kafka | Event-driven communication |
| kafka-go     | Kafka client               |
| Redis        | Rate limiting              |
| JWT          | Authentication             |
| bcrypt       | Password hashing           |
| SMTP         | Email delivery             |
| Docker       | Local infrastructure       |

---

# Services

## API Gateway

The API Gateway is the **single public entry point** for clients.

It is responsible for:

* Routing requests to internal services
* JWT authentication
* Extracting the authenticated user's ID
* Redis-based rate limiting
* Reverse proxying requests

Clients do not directly communicate with the internal services.

```text
Client
   │
   ▼
API Gateway
   │
   ├── /auth ──────────────► Auth Service
   │
   └── /notifications ────► Notification Service
```

Protected requests pass through:

```text
JWT Authentication
        │
        ▼
Redis Rate Limiter
        │
        ▼
  Reverse Proxy
        │
        ▼
Notification Service
```

### Gateway Port

```text
localhost:8084
```

---

# Auth Service

The Auth Service is responsible for user accounts and authentication.

### Responsibilities

* User registration
* User login
* Email uniqueness validation
* Password hashing
* Password verification
* JWT generation

Passwords are hashed with bcrypt before being stored.

The service retrieves users by email during login and compares the submitted password against the stored bcrypt hash.

### Authentication Flow

```text
POST /auth/login
       │
       ▼
 Auth Service
       │
       ├── Find user by email
       │
       ├── Compare password
       │
       └── Generate JWT
               │
               ▼
          Return JWT
```

The client then sends the token with subsequent protected requests:

```http
Authorization: Bearer <token>
```

---

# Notification Service

The Notification Service manages notification records and notification delivery state.

### Responsibilities

* Create notifications
* Retrieve notifications
* Validate notification requests
* Persist notifications
* Publish notification events
* Consume delivery result events
* Update notification status

The Notification Service owns **its own PostgreSQL database**.

It does not directly access the Auth Service's database.

---

# Email Worker

The Email Worker is an asynchronous Kafka consumer responsible for delivering email notifications.

It:

1. Consumes `notification.created`
2. Attempts to send the email through SMTP
3. Retries failed deliveries
4. Publishes `notification.sent` after successful delivery
5. Publishes `notification.failed` when all attempts fail

This keeps email delivery outside of the synchronous HTTP request.

```text
Notification Service
        │
        │ notification.created
        ▼
      Kafka
        │
        ▼
   Email Worker
        │
        │ SMTP
        ▼
    Email Server
```

---

# Database Architecture

The project follows a **database-per-service** approach.

Each service owns its own data.

```text
Auth Service
     │
     ▼
 Auth PostgreSQL
     │
     └── users


Notification Service
     │
     ▼
Notification PostgreSQL
     │
     └── notifications
```

The services do not directly access each other's databases.

This keeps the services loosely coupled and allows each service to evolve its data model independently.

For local development, the databases can still run on the same PostgreSQL server:

```text
             PostgreSQL
             /        \
            ▼          ▼
        auth_db    notification_db
```

The important architectural boundary is **data ownership**, not having physically separate PostgreSQL servers.

---

# Event-Driven Architecture

Kafka is used to communicate between the Notification Service and Email Worker.

## Notification Creation

When a notification is created:

```text
POST /notifications
        │
        ▼
Notification Service
        │
        ├── Save notification
        │
        └── Publish notification.created
                         │
                         ▼
                       Kafka
                         │
                         ▼
                    Email Worker
```

The API does not wait for the email to be sent.

This allows email delivery to happen asynchronously.

---

# Kafka Topics

| Topic                  | Producer             | Consumer             | Purpose                |
| ---------------------- | -------------------- | -------------------- | ---------------------- |
| `notification.created` | Notification Service | Email Worker         | Request email delivery |
| `notification.sent`    | Email Worker         | Notification Service | Successful delivery    |
| `notification.failed`  | Email Worker         | Notification Service | Failed delivery        |

The event flow is:

```text
notification.created
          │
          ▼
     Email Worker
      /        \
 success       failure
    │             │
    ▼             │
notification.sent │
                  │
                  ▼
         notification.failed
```

---

# Email Retry

Email delivery failures are retried before the notification is marked as failed.

The worker currently makes up to **3 delivery attempts**.

```text
Attempt 1
   │
   ├── Success ──► notification.sent
   │
   └── Failure
         │
         ▼
      Attempt 2
         │
         ├── Success ──► notification.sent
         │
         └── Failure
               │
               ▼
            Attempt 3
               │
               ├── Success ──► notification.sent
               │
               └── Failure ──► notification.failed
```

This protects against temporary SMTP failures.

---

# Notification Status

The Notification Service tracks the result of email delivery.

The basic lifecycle is:

```text
Created
   │
   ▼
Processing
   │
   ├────────► Sent
   │
   └────────► Failed
```

The Email Worker does not directly modify the Notification Service database.

Instead, it publishes an event:

```text
Email Worker
     │
     ▼
notification.sent
     │
     ▼
Kafka
     │
     ▼
Notification Service
     │
     ▼
Update notification status
```

This keeps the services decoupled.

---

# Rate Limiting

The API Gateway uses Redis to limit requests on a per-user basis.

The current limit is:

```text
100 requests per minute per authenticated user
```

Each user receives an independent Redis counter.

Example:

```text
rate_limit:user-A = 42
rate_limit:user-B = 17
```

The Redis key is generated from the authenticated user's ID:

```text
rate_limit:<user-id>
```

Redis `INCR` is used to atomically increment the request counter.

The key expires after the configured time window, allowing the counter to reset automatically.

Requests exceeding the limit receive:

```http
429 Too Many Requests
```

The gateway also exposes:

```http
X-RateLimit-Limit
X-RateLimit-Remaining
```

in successful responses.

### Rate Limiting Flow

```text
Request
   │
   ▼
JWT Middleware
   │
   │ Extract user ID
   ▼
Rate Limiter
   │
   ▼
Redis INCR
   │
   ├── Within limit ──► Continue
   │
   └── Exceeded ──────► 429
```

---

# Authentication

The API Gateway validates JWTs before allowing access to protected endpoints.

```text
Client
  │
  │ Authorization: Bearer <JWT>
  ▼
API Gateway
  │
  ▼
JWT Middleware
  │
  ├── Invalid ──► 401 Unauthorized
  │
  └── Valid
       │
       ▼
 Extract User ID
       │
       ▼
 Rate Limiter
       │
       ▼
 Notification Service
```

The JWT contains the authenticated user's identity and expiration information.

Invalid credentials result in:

```http
401 Unauthorized
```

---

# API

The API Gateway exposes the public API.

## Authentication

### Register

```http
POST /auth/register
Content-Type: application/json
```

Example:

```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

Response:

```http
201 Created
```

### Login

```http
POST /auth/login
Content-Type: application/json
```

Example:

```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

Successful response:

```http
200 OK
```

The response contains the JWT access token.

---

# Notifications

All notification endpoints require a valid JWT.

## Create Notification

```http
POST /notifications
Authorization: Bearer <JWT>
Content-Type: application/json
```

Example:

```json
{
  "recipient": "user@example.com",
  "channel": "email",
  "subject": "Welcome",
  "body": "Welcome to the platform!"
}
```

Response:

```http
201 Created
```

The returned notification ID can be used to retrieve the notification.

## Get Notification

```http
GET /notifications/:id
Authorization: Bearer <JWT>
```

Response:

```http
200 OK
```

Possible errors include:

```text
400 Bad Request
401 Unauthorized
404 Not Found
429 Too Many Requests
```

---

# Project Structure

```text
Notification-Platform/
│
├── api-gateway/
│   ├── cmd/
│   │   └── server/
│   ├── internal/
│   │   ├── config/
│   │   ├── middleware/
│   │   └── proxy/
│   └── test/
│
├── auth-service/
│   ├── cmd/
│   │   └── server/
│   ├── internal/
│   │   ├── config/
│   │   ├── db/
│   │   ├── handler/
│   │   ├── jwt/
│   │   ├── model/
│   │   ├── password/
│   │   ├── repository/
│   │   ├── service/
│   │   └── validation/
│   └── migrations/
│
├── notification-service/
│   ├── cmd/
│   │   └── server/
│   ├── internal/
│   │   ├── config/
│   │   ├── db/
│   │   ├── handler/
│   │   ├── kafka/
│   │   ├── middleware/
│   │   ├── model/
│   │   ├── repository/
│   │   ├── service/ 
│   │   └── validation/
│   └── migrations/
│
└── email-worker/
    ├── cmd/
    │   └── server/
    └── internal/
        ├── config/
        ├── kafka/
        ├── model/
        └── service/
```

Each service has its own Go module and can be developed independently.

---

# Running Locally

## Prerequisites

* Go
* PostgreSQL
* Docker
* SMTP account/provider

Kafka and Redis can be run using Docker.

## Infrastructure

The default local configuration expects:

```text
Kafka    localhost:9092
Redis    localhost:6379
```

PostgreSQL connection details are configured through environment variables.

---

## Kafka

Kafka can be run in KRaft mode without ZooKeeper.

### Linux / macOS

```bash
docker run -d \
  --name kafka \
  -p 9092:9092 \
  -e KAFKA_NODE_ID=1 \
  -e KAFKA_PROCESS_ROLES=broker,controller \
  -e KAFKA_LISTENER_SECURITY_PROTOCOL_MAP=CONTROLLER:PLAINTEXT,PLAINTEXT:PLAINTEXT \
  -e KAFKA_LISTENERS=PLAINTEXT://:9092,CONTROLLER://:9093 \
  -e KAFKA_ADVERTISED_LISTENERS=PLAINTEXT://localhost:9092 \
  -e KAFKA_CONTROLLER_QUORUM_VOTERS=1@localhost:9093 \
  -e KAFKA_CONTROLLER_LISTENER_NAMES=CONTROLLER \
  -e KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR=1 \
  -e KAFKA_TRANSACTION_STATE_LOG_REPLICATION_FACTOR=1 \
  -e KAFKA_TRANSACTION_STATE_LOG_MIN_ISR=1 \
  -e KAFKA_GROUP_INITIAL_REBALANCE_DELAY_MS=0 \
  apache/kafka:latest
```

### Windows PowerShell

```powershell
docker run -d `
  --name kafka `
  -p 9092:9092 `
  -e KAFKA_NODE_ID=1 `
  -e KAFKA_PROCESS_ROLES=broker,controller `
  -e KAFKA_LISTENER_SECURITY_PROTOCOL_MAP=CONTROLLER:PLAINTEXT,PLAINTEXT:PLAINTEXT `
  -e KAFKA_LISTENERS=PLAINTEXT://:9092,CONTROLLER://:9093 `
  -e KAFKA_ADVERTISED_LISTENERS=PLAINTEXT://localhost:9092 `
  -e KAFKA_CONTROLLER_QUORUM_VOTERS=1@localhost:9093 `
  -e KAFKA_CONTROLLER_LISTENER_NAMES=CONTROLLER `
  -e KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR=1 `
  -e KAFKA_TRANSACTION_STATE_LOG_REPLICATION_FACTOR=1 `
  -e KAFKA_TRANSACTION_STATE_LOG_MIN_ISR=1 `
  -e KAFKA_GROUP_INITIAL_REBALANCE_DELAY_MS=0 `
  apache/kafka:latest
```

Check that Kafka is running:

```bash
docker ps
```

### Create Kafka Topics

Create the topics required by the platform.

#### Linux / macOS

```bash
docker exec kafka /opt/kafka/bin/kafka-topics.sh \
  --create \
  --topic notification.created \
  --bootstrap-server localhost:9092

docker exec kafka /opt/kafka/bin/kafka-topics.sh \
  --create \
  --topic notification.sent \
  --bootstrap-server localhost:9092

docker exec kafka /opt/kafka/bin/kafka-topics.sh \
  --create \
  --topic notification.failed \
  --bootstrap-server localhost:9092
```

#### Windows PowerShell

```powershell
docker exec kafka /opt/kafka/bin/kafka-topics.sh --create --topic notification.created --bootstrap-server localhost:9092

docker exec kafka /opt/kafka/bin/kafka-topics.sh --create --topic notification.sent --bootstrap-server localhost:9092

docker exec kafka /opt/kafka/bin/kafka-topics.sh --create --topic notification.failed --bootstrap-server localhost:9092
```

List the topics to verify:

```bash
docker exec kafka /opt/kafka/bin/kafka-topics.sh \
  --list \
  --bootstrap-server localhost:9092
```

PowerShell:

```powershell
docker exec kafka /opt/kafka/bin/kafka-topics.sh --list --bootstrap-server localhost:9092
```

---

## Redis

Redis is used by the API Gateway for rate limiting.

### Linux / macOS

```bash
docker run -d \
  --name redis \
  -p 6379:6379 \
  redis:latest
```

### Windows PowerShell

```powershell
docker run -d `
  --name redis `
  -p 6379:6379 `
  redis:latest
```

Check Redis:

```bash
docker exec redis redis-cli ping
```

Expected:

```text
PONG
```

PowerShell uses the same command:

```powershell
docker exec redis redis-cli ping
```

---

## Environment Variables

Each service has its own environment configuration. Copy the corresponding `.env.example` file to `.env` and replace the placeholder values.

### Auth Service

Create `auth-service/.env`:

```env
DATABASE_URL=postgres://postgres:password@localhost:5432/auth_db
JWT_SECRET=your-super-secret-jwt-key
```

### Notification Service

Create `notification-service/.env`:

```env
DATABASE_URL=postgres://postgres:password@localhost:5432/notification_db
```

### Email Worker

Create `email-worker/.env`:

```env
SMTP_HOST=smtp.example.com
SMTP_PORT=587
SMTP_USERNAME=your-email@example.com
SMTP_PASSWORD=your-smtp-password
```

### API Gateway

Create `api-gateway/.env`:

```env
JWT_SECRET=your-super-secret-jwt-key
```

> **Note:** The `JWT_SECRET` must be the same in the Auth Service and API Gateway.

> **Security:** Never commit `.env` files containing real credentials, passwords, JWT secrets, or SMTP credentials.

---

## Start Services

Run each service separately.

### Auth Service

```bash
cd auth-service
go run ./cmd/server
```

### Notification Service

```bash
cd notification-service
go run ./cmd/server
```

### Email Worker

```bash
cd email-worker
go run ./cmd/server
```

### API Gateway

```bash
cd api-gateway
go run ./cmd/server
```

The API Gateway is the public entry point:

```text
http://localhost:8084
```

Clients should communicate with the API Gateway rather than directly accessing the internal services.

---

# Testing

The project includes API Gateway integration tests covering both successful and failure scenarios.

Run:

```bash
go test ./test -v
```

The test suite covers:

* User registration
* User login
* Invalid login credentials
* Unknown email
* Invalid request bodies
* Missing authorization headers
* Invalid JWTs
* Malformed Bearer tokens
* Empty Bearer tokens
* Invalid notification IDs
* Notification creation
* Notification retrieval
* Complete notification flow
* Rate limiting

The main happy path is:

```text
Register
   │
   ▼
Login
   │
   ▼
Receive JWT
   │
   ▼
Create Notification
   │
   ▼
Notification ID
   │
   ▼
Retrieve Notification
```

The notification delivery path is asynchronous:

```text
Create Notification
       │
       ▼
notification.created
       │
       ▼
  Email Worker
       │
       ▼
     SMTP
       │
       ▼
notification.sent / notification.failed
```

---

# Design Decisions

## API Gateway

Clients interact with a single public endpoint rather than needing to know the location of individual services.

This centralizes:

* Authentication
* Rate limiting
* Routing

## Database Per Service

Each service owns its own database.

This prevents direct database coupling and allows services to evolve their schemas independently.

## Kafka

Kafka is used for asynchronous communication because email delivery should not block the notification creation request.

## Redis

Redis provides shared, atomic request counters for rate limiting.

This allows the rate limiter to work even if multiple API Gateway instances are running.

## Email Worker

Email delivery is isolated into a separate worker so the Notification Service does not need to handle SMTP operations itself.

## Retries

Temporary email delivery failures are retried before publishing a final failure event.

---

# Security

The project implements several security measures:

* JWT-based authentication
* JWT expiration
* bcrypt password hashing
* Generic invalid-login responses
* Request validation
* Protected notification endpoints
* Redis-based rate limiting
* Parameterized database queries
* Environment-based secrets
* Authentication at the API Gateway

The Notification Service is not intended to be accessed directly by external clients.

---

# Future Improvements

Possible extensions include:

* Delayed Kafka retry topics
* Dead-letter queue
* Exponential backoff for retries
* Kafka consumer group scaling
* Additional notification channels
* Notification templates
* User notification preferences
* Delivery analytics
* Structured logging
* Distributed tracing
* Docker Compose for the complete local environment
* CI/CD pipeline
* OpenAPI documentation

---

# License

This project is licensed under the MIT License.
