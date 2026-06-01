# Backend — Desafio Fullstack Pleno

## Stack
Go 1.26, Gin, pgx/v5, golang-jwt/v5, golang-migrate/v4, zerolog, godotenv, gin-contrib/cors

## Estrutura

```
backend/
├── cmd/api/main.go          # Entrypoint
├── internal/
│   ├── auth/                # JWT, middleware, session
│   │   ├── handler.go       # POST /auth/token, GET/DELETE /auth/session
│   │   ├── middleware.go    # Bearer + cookie extraction
│   │   ├── jwt.go           # GenerateToken, ValidateToken (HS256)
│   │   └── service.go       # Credenciais fixas, delega ao jwt
│   ├── children/            # CRUD com filtros e paginação
│   │   ├── handler.go       # List, GetByID, MarkReviewed, ListNeighborhood
│   │   ├── service.go       # Orquestra repo, monta response com alertas
│   │   ├── repository.go    # SQL com QueryBuilder dinâmico
│   │   ├── models.go        # Child, ChildById, Health, Education, SocialAssistance
│   │   ├── filters.go       # Name, Neighborhood, Alert, Reviewed, HasAlert, Page, PerPage
│   │   └── queries.go       # QueryBuilder — constrói WHERE/LIMIT/OFFSET
│   ├── statistics/          # Agregações (summary + alertas por bairro)
│   │   ├── handler.go       # GET /summary, GET /statistics
│   │   ├── service.go       # Summary + StatisticsByNeighborhood
│   │   ├── repository.go    # COUNTs + JOIN grouped by neighborhood
│   │   └── model.go         # Summary, NeighborhoodAlertCount, StatisticsResponse
│   ├── config/config.go     # Load() — PORT, DATABASE_URL, JWT_SECRET, ALLOW_ORIGINS
│   ├── database/
│   │   ├── postgres.go      # Connect — pgxpool
│   │   ├── migrations.go    # RunMigrations — golang-migrate
│   │   ├── seed.go          # LoadSeed — JSON → children + áreas + alertas
│   │   ├── inserts.go       # insertChild, insertHealth, insertEducation, insertSocialAssistance, alert inserts
│   │   └── testutil/        # testdb.go — helpers para testes de integração
│   └── server/router.go     # SetupRouter — CORS, logger, rotas, DI
├── migrations/
│   ├── 001_create_children.up.sql
│   └── 002_create_areas.up.sql   # health, education, social_assistance + alert_* tables
├── data/seed.json
├── Dockerfile               # Multi-stage: golang:1.26-alpine → alpine:3.21
├── Makefile                 # run, build, test, clean
├── go.mod / go.sum
└── .env
```

## Arquitetura
Handler → Service → Repository. Injeção de dependências em `server/router.go`.

## Autenticação
- JWT HS256, 1h de expiração
- Token aceito via `Authorization: Bearer` **ou** cookie `auth_token` (definido no login, limpo no logout)
- Credenciais fixas: `tecnico@prefeitura.rio` / `painel@2024`

## Rotas da API

### Públicas
| Método | Path | Descrição |
|--------|------|-----------|
| POST | `/auth/token` | Login — retorna JWT + seta cookie |
| GET | `/auth/session` | Valida cookie e retorna token info |
| DELETE | `/auth/session` | Logout — limpa cookie |

### Protegidas (prefixo `/api/v1`)
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

## Banco de Dados

- **PostgreSQL** via pgxpool
- Migrations automáticas com `golang-migrate` na startup
- Seed automático via `data/seed.json` (executado apenas se DB vazio)
- Reset: `docker compose down -v`

### Tabelas
- `children` — id, name, age, neighborhood, reviewed, reviewed_by, reviewed_at, notes, created_at
- `health` — child_id, vaccinations_up_to_date, last_consultation
- `education` — child_id, school_name, frequency_percent
- `social_assistance` — child_id, cad_unico, active_benefit
- `alert_health`, `alert_education`, `alert_social_assistance` — área_id, code, message, created_at

## Testes

### Unitários (sem banco)
```bash
make test-unit
# ou
go test ./internal/auth/... ./internal/children/... ./internal/statistics/... ./internal/config/... -v
```

### Integração (precisa PostgreSQL rodando)
```bash
make test-integration
# ou
TEST_DATABASE_URL="postgres://postgres:postgres@localhost:5432/full-stack-test?sslmode=disable" \
  go test ./internal/database/... ./internal/children/... ./internal/statistics/... -v -p=1 -run Integration
```

### Todos os testes
```bash
make test
# ou
go test ./... -v -p=1
```

> `-p=1` evita concorrência entre pacotes que compartilham o mesmo banco de dados.

### Estrutura de testes
```
internal/
├── auth/
│   ├── jwt_test.go              # GenerateToken, ValidateToken
│   ├── middleware_test.go       # extractToken, AuthMiddleware
│   ├── service_test.go          # Authenticate, ValidateToken
│   └── handler_test.go          # Token, Session, Logout
├── children/
│   ├── mock_repository.go       # Mock da interface ChildRepository
│   ├── filters_test.go          # Normalize, Offset
│   ├── queries_test.go          # QueryBuilder
│   ├── service_test.go          # List, GetByID, MarkReviewed
│   ├── handler_test.go          # Handlers com mock
│   └── repository_test.go       # Integração com PostgreSQL (package children_test)
├── statistics/
│   ├── mock_repository.go       # Mock da interface StatisticRepository
│   ├── service_test.go          # GetSummary, GetStatistics
│   ├── handler_test.go          # Handlers com mock
│   └── repository_test.go       # Integração com PostgreSQL
├── config/
│   └── config_test.go           # Load defaults + env vars
└── database/
    ├── testutil/testdb.go       # Helpers: Connect, Migrate, Truncate, Seed, Insert*
    └── database_test.go         # Integração: Connect, Migrations, LoadSeed
```

## Ambiente
| Variável | Default |
|----------|---------|
| `PORT` | `8080` |
| `DATABASE_URL` | `postgres://postgres:postgres@localhost:5432/full-stack?sslmode=disable` |
| `JWT_SECRET` | `super-secret` |
| `ALLOW_ORIGINS` | `http://localhost:3000` |
| `BASE_PATH` | `.` |
