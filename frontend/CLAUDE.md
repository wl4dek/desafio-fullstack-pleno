# Frontend — Painel de Monitoramento Infantil

## Tecnologias

- **Next.js 16.2.6** (App Router) + **React 19.2.4**
- **TypeScript** ~5
- **Tailwind CSS v4** com `@theme inline` + `cn()` (clsx + tailwind-merge)
- **shadcn/ui** via `@radix-ui/*` + `class-variance-authority`
- **Zustand** (auth store), **SWR** (data fetching/cache)
- **react-hook-form** + **zod** (formulários + validação)
- **react-leaflet** + **leaflet** (mapa estatístico)
- **next-themes** (dark/light mode)
- **lucide-react** (ícones)
- **Vitest** + **@testing-library/react** + **@testing-library/user-event** + **jsdom**

## Path alias

`@/` → `./src/`

## Estrutura

```
src/
  __tests__/              # Testes Vitest (*.test.tsx)
    setup.ts
    schemas.test.ts
    LoginForm.test.tsx
    ChildrenTable.test.tsx
    SummaryCards.test.tsx
    ChildrenFilters.test.tsx
    ReviewButton.test.tsx
    ChildDetail.test.tsx
    Header.test.tsx
  app/                    # App Router (pages + layouts)
    layout.tsx
    client-layout.tsx
    page.tsx              # redirect /dashboard
    globals.css
    login/page.tsx
    dashboard/page.tsx
    children/page.tsx
    children/[id]/page.tsx
    statistics/page.tsx
  components/
    ui/                   # shadcn/ui (button, card, input, select, table, tabs, etc.)
    Header.tsx
    AuthGuard.tsx
  data/
    LimiteBairros.json    # GeoJSON dos bairros do Rio
  features/               # Feature modules (domain-driven)
    auth/components/LoginForm.tsx
    children/components/
      ChildrenTable.tsx
      ChildrenFilters.tsx
      ChildDetail.tsx
      ReviewButton.tsx
    dashboard/components/SummaryCards.tsx
    statistics/components/AlertsArea.tsx
  hooks/
    useAuth.ts            # useLogin (auth + redirect)
    useChildren.ts        # useChildren, useChild, useNeighborhoods
    useSummary.ts         # useSummary
    useStatistic.ts       # useStatistic
    use-toast.ts          # Toast system
  lib/
    api.ts               # HTTP client (Bearer token, 401 handling)
    utils.ts             # cn(), formatDateBR()
  schemas/
    index.ts             # loginSchema (zod)
  services/
    auth.ts              # POST /auth/token
    children.ts          # children CRUD + summary + review
    statistic.ts         # GET /api/v1/statistics
  stores/
    auth.ts              # useAuthStore (Zustand)
    alert.ts             # Alerts map (código → label)
  types/
    index.ts             # Interfaces TS (Child, Summary, etc.)
```

## Convenções

- **Componentes:** PascalCase, named exports
- **Hooks:** camelCase com prefixo `use`
- **Serviços/stores:** camelCase, named exports
- **UI components:** lowercase kebab (`button.tsx`, `alert-dialog.tsx`)
- **"use client"** em componentes com hooks/interatividade; páginas simples são server components
- **Importações:** paths absolutos com `@/` (ex: `@/features/auth/components/LoginForm`)

## Páginas

| Rota | Componente | Descrição |
|------|-----------|-----------|
| `/` | redirect → `/dashboard` | |
| `/login` | `LoginForm` | Email/senha com validação zod |
| `/dashboard` | `SummaryCards` | 4 cards (total, revisadas, pendentes, alertas) + alertas por área |
| `/children` | `ChildrenFilters` + `ChildrenTable` | Tabela com filtros (busca, bairro, alerta, revisão) |
| `/children/[id]` | `ChildDetail` | Detalhes + tabs (Saúde, Assistência Social, Educação) + marcar revisão |
| `/statistics` | `AlertsArea` | Mapa Leaflet com GeoJSON por bairro |

## Estado

- **Auth:** `useAuthStore` (Zustand) — token em `localStorage`, `hydrate` via `AuthGuard` (`GET /auth/session`)
- **Data fetching:** SWR com keys `/children`, `/children/:id`, `/children/neighborhood`, `/summary`, `/api/v1/statistics`
- **Filtros:** mantidos em URL search params com debounce de 500ms

## API

- `src/lib/api.ts` — client HTTP com Bearer token injetado do Zustand
- `NEXT_PUBLIC_API_URL` env var (default `http://localhost:8080`)
- Rewrites no `next.config.ts`
- Tratamento automático de 401 (limpa token e redireciona)

## Serviços

| Serviço | Endpoint | Descrição |
|---------|----------|-----------|
| `login()` | `POST /auth/token` | Autenticação |
| `fetchChildren()` | `GET /api/v1/children` | Lista paginada com filtros |
| `fetchChild()` | `GET /api/v1/children/:id` | Detalhes da criança |
| `fetchSummary()` | `GET /api/v1/summary` | Resumo do dashboard |
| `markReviewed()` | `PATCH /api/v1/children/:id/review` | Marcar como revisado |
| `listNeighborhood()` | `GET /api/v1/children/neighborhood` | Lista de bairros |
| `fetchStatistic()` | `GET /api/v1/statistics` | Dados estatísticos por bairro |

## Testes

```bash
npm run test        # vitest run
npm run test:watch  # vitest
```

### Arquivos de teste (8 arquivos, 63 testes)

| Arquivo | Testes | O que cobre |
|---------|--------|-------------|
| `schemas.test.ts` | 5 | Validação Zod (email válido, inválido, senha vazia, campos ausentes) |
| `LoginForm.test.tsx` | 6 | Render, erros validação (vazio + email inválido), submit, loading, erro servidor |
| `ChildrenTable.test.tsx` | 8 | Loading skeleton, erro+retry, vazio, dados, badges, clique navegação |
| `SummaryCards.test.tsx` | 7 | Loading skeleton, erro+retry, null, cards valores, alerts badges, percentuais, vazio |
| `ChildrenFilters.test.tsx` | 8 | Input busca, 3 selects, debounce 500ms, reseta page, clear, URL params |
| `ReviewButton.test.tsx` | 7 | Render botão, diálogo, chama API, mutate cache, toast sucesso/erro, loading |
| `ChildDetail.test.tsx` | 14 | Loading, erro, não encontrado, info, badge, tabs, AreaSection, ReviewButton, voltar |
| `Header.test.tsx` | 7 | Não autenticado (null), links nav, logout, menu mobile, toggle tema |

### Padrões

- **Mocks mutáveis:** `let mockReturn` permite testar diferentes estados no mesmo arquivo
- **Reset:** `beforeEach` com `vi.clearAllMocks()` + reset dos mocks
- **Queries:** `getByText`, `getByRole`, `findByText` (async), `queryByText`
- **Interação:** `userEvent.setup()` com `await user.type/click/clear`
- **Form submit com validação:** `fireEvent.submit(form)` (em vez de clicar no botão)
- **Busca por placeholder nos selects Radix:** `getByText("Bairro")`, `getByText("Alerta")`, etc.

### Limitações conhecidas

- **Radix UI Select** usa `PointerEvent` → `hasPointerCapture` não implementado em jsdom. Interações de seleção são testadas apenas via rendering (verificar placeholders). A lógica de filtro é validada através do input de busca.
- **userEvent.type + react-hook-form + jsdom:** em alguns casos pode ser necessário `fireEvent.submit` diretamente no form para acionar a validação.

## Idioma

Português (pt-BR) — labels, mensagens, toast, formatos de data.
