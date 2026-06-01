# API de Acompanhamento Infantil

API REST em Go para gerenciamento e acompanhamento de crianças monitoradas por técnicos da prefeitura.

## Stack

- Go 1.26+
- Gin (HTTP framework) + gin-contrib/cors
- PostgreSQL (pgx/v5)
- JWT (golang-jwt/v5)
- golang-migrate (migrations)
- zerolog (logging)
- godotenv (config)

## Estrutura

```
backend/
├── cmd/api/main.go
├── internal/
│   ├── auth/                # JWT, middleware, session
│   ├── children/            # CRUD com filtros e paginação
│   ├── statistics/          # Agregações (summary + alertas por bairro)
│   ├── config/              # Config via env vars
│   ├── database/            # Conexão, migrations, seed + testutil
│   └── server/              # Router, CORS, DI
├── migrations/              # SQL migrations
├── data/seed.json           # Dados iniciais
├── Dockerfile               # Multi-stage
├── Makefile
├── go.mod / go.sum
└── .env
```

## Setup

```bash
docker compose up --build
```

### Dev manual

```bash
# Backend
cd backend && go run ./cmd/api

# Frontend
cd frontend && npm run dev
```

## Derrubar ambiente

```bash
docker compose down
```

## Reset completo

```bash
docker compose down -v
```

## Autenticação

- JWT assinado com HS256, expira em 1h
- Token aceito via `Authorization: Bearer` **ou** cookie `auth_token`
- Cookie `auth_token` é definido no login e limpo no logout
- Credenciais fixas: `tecnico@prefeitura.rio` / `painel@2024`

### Endpoints de autenticação (públicos)

| Método | Path | Descrição |
|--------|------|-----------|
| POST | `/auth/token` | Login — retorna JWT + seta cookie |
| GET | `/auth/session` | Valida cookie e retorna info do token |
| DELETE | `/auth/session` | Logout — limpa cookie |

## API

Rotas protegidas usam prefixo `/api/v1` e exigem token JWT.

| Método | Path | Descrição |
|--------|------|-----------|
| GET | `/api/v1/children` | Lista paginada com filtros |
| GET | `/api/v1/children/neighborhood` | Lista bairros disponíveis |
| GET | `/api/v1/children/:id` | Detalhes da criança + áreas + alertas |
| PATCH | `/api/v1/children/:id/review` | Marcar como revisado |
| GET | `/api/v1/summary` | Total, revisados, pendentes, alertas por área |
| GET | `/api/v1/statistics` | Alertas agregados por bairro |

### Query params — `GET /api/v1/children`

`childName`, `neighborhood`, `reviewed` (bool), `has_alert` (bool), `alert` (código), `page` (≥1), `per_page` (10–50)

## Testes

```bash
# Todos os testes
make test

# Unitários (sem banco)
make test-unit

# Integração (precisa PostgreSQL rodando)
make test-integration
```

> Use `-p=1` para evitar concorrência entre pacotes que compartilham o mesmo banco.

## Ambiente

| Variável | Default |
|----------|---------|
| `PORT` | `8080` |
| `DATABASE_URL` | `postgres://postgres:postgres@localhost:5432/full-stack?sslmode=disable` |
| `JWT_SECRET` | `super-secret` |
| `ALLOW_ORIGINS` | `http://localhost:3000` |
| `BASE_PATH` | `.` |

## Exemplos

```bash
# Login
curl -X POST localhost:8080/auth/token \
  -H "Content-Type: application/json" \
  -d '{"email":"tecnico@prefeitura.rio","password":"painel@2024"}'

# Listar crianças (com token)
TOKEN="seu-jwt-aqui"
curl localhost:8080/api/v1/children \
  -H "Authorization: Bearer $TOKEN"

# Verificar sessão
curl localhost:8080/auth/session \
  -b "auth_token=$TOKEN"

# Logout
curl -X DELETE localhost:8080/auth/session \
  -b "auth_token=$TOKEN"
```
