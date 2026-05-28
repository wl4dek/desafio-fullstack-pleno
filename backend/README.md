# API de Acompanhamento Infantil

API REST em Go para gerenciamento e acompanhamento de crianças monitoradas por técnicos da prefeitura.

## Stack

- Go 1.26+
- Gin (HTTP framework)
- PostgreSQL
- Docker / Docker Compose
- JWT (golang-jwt/v5)
- pgx/v5 (PostgreSQL driver)
- golang-migrate (migrations)
- zerolog (logging)

## Setup

```bash
docker compose up --build
```

## Derrubar ambiente

```bash
docker compose down
```

## Reset completo

```bash
docker compose down -v
```

## Rodar testes

```bash
go test ./...
```

## Exemplo de autenticação

```bash
curl -X POST localhost:8080/auth/token \
  -H "Content-Type: application/json" \
  -d '{
    "email":"tecnico@prefeitura.rio",
    "password":"painel@2024"
  }'
```
