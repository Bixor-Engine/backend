# 🔁 CEX Backend — Go + Rust + PostgreSQL

A high-performance centralized crypto exchange backend built with **Go** (API & services), **Rust** (matching engine), and **PostgreSQL** (data storage). Designed for scalability, modularity, and performance.

---

## 🏗️ Architecture Overview

```text
+--------------------+       +-------------------+
|     Client UI      | <-->  |   Go API Gateway  |
+--------------------+       +-------------------+
                                     |
                                     v
+--------------------+       +-------------------+
|   PostgreSQL DB    | <-->  |  Go Settlement    |
+--------------------+       +-------------------+
                                     ^
                                     |
                            +-------------------+
                            |  Rust Match Engine|
                            +-------------------+
```

---

## 📦 Components

### 1. API Gateway (Go)

* REST/gRPC endpoints
* JWT auth, session handling
* Rate limiting
* WebSocket price/trade feed

### 2. Matching Engine (Rust)

* Pure Rust high-performance engine
* Central Limit Order Book (CLOB)
* Fast trade matching & queuing
* Emits `TradeExecuted`, `OrderUpdated`

### 3. Settlement Service (Go)

* Balance updates, freezing
* Fee logic, PnL, trade logging
* Wallet management & DB updates

### 4. PostgreSQL Database

* Persistent data store for:

  * Users
  * Wallets
  * Orders & trades
  * Markets & fees
  * Deposits/withdrawals

---

## 🛠️ Tech Stack

| Layer           | Tech                            |
| --------------- | ------------------------------- |
| Language (API)  | Go (1.22+)                      |
| Matching Engine | Rust (2021 Edition)             |
| Database        | PostgreSQL 15+                  |
| Message Queue   | Redis / NATS / Kafka (optional) |
| WebSocket       | gorilla/websocket               |
| Auth            | JWT / OAuth2                    |
| ORM             | GORM / pgx (Go)                 |

---

## 🧰 Folder Structure

```text
exchange-backend/
├── cmd/
│   ├── api-gateway/       # Go REST + WS API
│   ├── settlement/        # Go settlement logic
│   └── market-engine/     # Rust binary interface
│
├── engine/                # Rust matching engine crate
│   └── src/
│
├── internal/
│   ├── models/            # Shared Go models
│   └── utils/             # Config, logger, helpers
│
├── proto/                 # gRPC/message schema (optional)
├── config/                # .env, configs
├── scripts/               # DB migrations, admin tasks
└── deploy/                # Docker, K8s configs
```

---

## ✅ Workflow Example

### Order Placement Flow:

1. **Client** sends `POST /order` to **Go API Gateway**
2. API forwards to **Rust Matching Engine**
3. Engine emits events: `TradeExecuted`, `OrderMatched`
4. **Go Settlement Service** listens and updates:

   * Balances
   * Orders
   * Trade history
5. **WebSocket server** broadcasts updates to client

---

## 👋 Example PostgreSQL Tables

* `users(id, email, password_hash, ...)`
* `wallets(user_id, currency, balance, locked)`
* `orders(id, user_id, market, price, qty, side, status)`
* `trades(id, buy_order_id, sell_order_id, price, qty)`
* `markets(symbol, base, quote, min_qty, status)`
* `deposits`, `withdrawals`, `sessions`

---

## 📊 Performance Goals

| Component       | Target TPS / Latency |
| --------------- | -------------------- |
| Matching Engine | > 10,000 orders/sec  |
| Go REST API     | > 3,000 req/sec      |
| PostgreSQL      | Tuned with indexing  |
| WebSocket Push  | < 200ms latency      |

---

## 🚀 Setup

Coming soon: Docker Compose & Kubernetes manifests to spin up full local stack.

---

## 🔒 Security Best Practices

* Enforce JWT expiration + 2FA
* Rate limit login/order APIs
* TLS for internal services
* Withdrawals require multi-sig/manual approval
* Engine is stateless and sandboxed

---

## 📝 License

MIT License. Build your exchange freely, safely, and securely.

---

## 🤝 Contribute

PRs and feedback welcome! Fork and star if this helps your project.

---
