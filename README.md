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
| Middleware support                    | ✅        | Auth, Rate-limit |
| Rate limiting                         | ✅        | Custom middleware |
| MySQL integration                     | ✅        | With migrations |
| Dockerized deployment                 | ✅        | Separate dev & prod configs |
| Migration management                  | ✅        | CLI migration scripts |
| API versioning                        | ✅        | e.g., `/api/v1` routes |
| Refresh token revocation              | ✅      | Not implemented |
| Multi-role support                    | ❌        | Only one role per user |
| Swagger / API documentation           | ❌        | Planned for future |
| Logging                               | ❌        | Planned for future |
| seeder                                | ❌        | Planned for future |

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
```
## development mode
### in development make sure you have go in your local machine
```bash
#1. make sure you have air on your go path, if not run this
go install github.com/air-verse/air@latest

#set up db using docker compose
docker compose -p dev --env-file .env.dev -f docker-compose.dev.yml up -d

#migrate db
./run.sh migrate-up

#then
./run.sh run-dev

##url
http://localhost:8080/api/v1/posts

```

## production mode
``` bash
./run-prod.sh

./migrate-prod.sh

##url
http://localhost:8081/api/v1/posts

```
## ⭐ Don't Forget to Star This Repository

If you find this project useful, please give it a ⭐ on GitHub to show your support!  

