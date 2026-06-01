# Desafio Fullstack Pleno

Fullstack application for monitoring children tracked by city technicians.

## Stack

- **Backend:** Go 1.26+, Gin, PostgreSQL, pgx/v5, golang-jwt/v5, golang-migrate, zerolog, godotenv
- **Frontend:** Next.js 16, TypeScript, Tailwind CSS v4, shadcn/ui, Lucide React, SWR, Zustand, React Hook Form, Zod, react-leaflet, Vitest
- **Infra:** Docker, Docker Compose (PostgreSQL 17)

## Documentação por projeto

Cada subprojeto tem seu próprio `CLAUDE.md` com detalhes completos:

- [`frontend/CLAUDE.md`](frontend/CLAUDE.md) — estrutura, páginas, serviços, testes, convenções
- [`backend/CLAUDE.md`](backend/CLAUDE.md) — arquitetura, rotas, banco, testes

## Estrutura

```
├── backend/           # Go API
├── frontend/          # Next.js app
├── e2e/
├── specs/
└── docker-compose.yml
```

## Comandos

```bash
docker compose up --build          # Start everything
docker compose up --build api      # Start backend only
cd frontend && npm run dev         # Frontend dev
cd backend && go run ./cmd/api     # Backend dev
cd frontend && npm test            # Frontend tests
cd backend && make test            # Backend tests
docker compose down                # Tear down
docker compose down -v             # Full reset (wipes DB)
```

## Ambiente

| Variável | Default | Projeto |
|----------|---------|---------|
| `PORT` | `8080` | backend |
| `DATABASE_URL` | `postgres://postgres:postgres@localhost:5432/children?...` | backend |
| `JWT_SECRET` | `super-secret` | backend |
| `ALLOW_ORIGINS` | `http://localhost:3000` | backend |
| `API_URL` | `""` (proxy rewrites) | frontend |

## Auth

- Credenciais fixas: `tecnico@prefeitura.rio` / `painel@2024`
- JWT HS256, 1h expiry, enviado via `Authorization: Bearer` ou cookie `auth_token`
- Rota pública: `/login`; demais rotas protegidas por `AuthGuard` (client-side)

## API Endpoints

| Método | Path | Auth | Descrição |
|--------|------|------|-----------|
| POST | `/auth/token` | Não | Login |
| GET | `/auth/session` | Cookie | Valida sessão |
| DELETE | `/auth/session` | Cookie | Logout |
| GET | `/api/v1/children` | Sim | Lista paginada/filtrada |
| GET | `/api/v1/children/neighborhood` | Sim | Lista bairros |
| GET | `/api/v1/children/:id` | Sim | Detalhes da criança |
| PATCH | `/api/v1/children/:id/review` | Sim | Marcar revisão |
| GET | `/api/v1/summary` | Sim | Dashboard |
| GET | `/api/v1/statistics` | Sim | Alertas por bairro |

## Key Conventions

- **Backend:** Handler → Service → Repository com injeção de dependências
- **Frontend:** Feature modules em `features/`, componentes em `components/ui/`
- Erros: `{"error": "message"}`
- Timestamps em UTC
- Mobile-first (375px → 1440px)
