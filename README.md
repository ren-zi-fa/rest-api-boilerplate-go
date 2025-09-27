# REST API with Go and Chi

A clean, modular REST API built with **Go (Golang)** using the [Chi router](https://github.com/go-chi/chi).  
Includes authentication, role-based access control, middleware, and MySQL integration with migrations.

---

## 📌 Feature Table

| Feature                               | Available | Notes |
|---------------------------------------|:---------:|-------|
| Authentication (JWT)                 | ✅        | Access & Refresh token support |
| Password hashing (bcrypt)            | ✅        | Secure password storage |
| Role-based Access Control (RBAC)     | ✅        | Single-role per user (not multi-role) |
| Middleware support                    | ✅        | Auth, Rate-limit, Logging, etc. |
| Rate limiting                         | ✅        | Custom middleware |
| MySQL integration                     | ✅        | With migrations |
| Dockerized deployment                 | ✅        | Separate dev & prod configs |
| Migration management                  | ✅        | CLI migration scripts |
| API versioning                        | ✅        | e.g., `/api/v1` routes |
| Multi-role support                    | ❌        | Only one role per user |
| Refresh token revocation              | ❌        | Not implemented |
| Swagger / API documentation           | ❌        | Planned for future |

---

## ⚙️ Tech Stack

| Component      | Technology |
|----------------|------------|
| Language       | Go (Golang) |
| Router         | [Chi](https://github.com/go-chi/chi) |
| Database       | MySQL |
| Auth           | JWT |
| Password Hash  | bcrypt |
| Container      | Docker & Docker Compose |
| Env Management | `.env` file |

---

## 🚀 Getting Started

### Clone this repository

```bash
git clone https://github.com/ren-zi-fa/rest-api-boilerplate-go
cd rest-api-boilerplate-go

./run-prod.sh