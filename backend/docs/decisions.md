# Architectural Decisions

## PostgreSQL

PostgreSQL was chosen for its reliability, JSON support, and strong ecosystem. The relational model fits the structured child monitoring data well, and PostgreSQL's mature tooling (golang-migrate, pgx) integrates smoothly with Go.

## Seed Strategy

Seed data in `data/seed.json` is loaded automatically at startup only when the database is empty. This avoids duplication while keeping the initial dataset reproducible. A full reset is achieved by removing Docker volumes.

## Docker

Docker Compose provides a consistent environment across development and testing. The API and PostgreSQL run in separate containers with a network bridge, eliminating environment-specific issues.

## Layered Architecture

Handler → Service → Repository separation keeps HTTP concerns, business logic, and data access decoupled. This makes the codebase testable and maintainable, with each layer independently mockable.

## Future Improvements

- Add database indexes on frequently filtered columns (neighborhood, has_alert, reviewed)
- Implement soft delete for children records
- Add an audit log for review actions
- Introduce refresh tokens and RBAC for multiple technician roles
- Add OpenTelemetry instrumentation for observability
- Generate OpenAPI/Swagger documentation
