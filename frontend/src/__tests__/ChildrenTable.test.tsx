import { render, screen } from "@testing-library/react"
import { describe, it, expect, vi } from "vitest"
import { ChildrenTable } from "@/features/children/components/ChildrenTable"

const mockData = {
  data: [
    {
      id: "c001",
      name: "Ana Clara Mendes",
      age: 6,
      neighborhood: "Rocinha",
      area: "Saúde",
      has_alert: true,
      reviewed: false,
      reviewed_by: null,
      reviewed_at: null,
      notes: "",
      created_at: "2024-01-01T00:00:00Z",
    },
    {
      id: "c002",
      name: "Lucas Ferreira",
      age: 7,
      neighborhood: "Maré",
      area: "Educação",
      has_alert: false,
      reviewed: true,
      reviewed_by: "tecnico@prefeitura.rio",
      reviewed_at: "2024-06-01T00:00:00Z",
      notes: "",
      created_at: "2024-01-02T00:00:00Z",
    },
  ],
  pagination: { page: 1, per_page: 10, total: 2, total_pages: 1 },
}

let mockReturn = { data: mockData, isLoading: false, isError: false, error: null, refresh: vi.fn() }

vi.mock("@/hooks/useChildren", () => ({
  useChildren: () => mockReturn,
}))

vi.mock("next/navigation", () => ({
  useRouter: () => ({ push: vi.fn(), replace: vi.fn() }),
  useSearchParams: () => new URLSearchParams(),
}))

describe("ChildrenTable", () => {
  it("renders children data correctly", () => {
    render(<ChildrenTable filters={{}} />)

    expect(screen.getByText("Ana Clara Mendes")).toBeInTheDocument()
    expect(screen.getByText("Lucas Ferreira")).toBeInTheDocument()
    expect(screen.getByText("Rocinha")).toBeInTheDocument()
    expect(screen.getByText("Maré")).toBeInTheDocument()
  })

  it("shows alert badges", () => {
    render(<ChildrenTable filters={{}} />)

    const alertBadges = screen.getAllByText("Alerta")
    expect(alertBadges.length).toBeGreaterThanOrEqual(1)
  })

  it("shows empty state when no data", () => {
    mockReturn = { data: { data: [], pagination: { page: 1, per_page: 10, total: 0, total_pages: 0 } }, isLoading: false, isError: false, error: null, refresh: vi.fn() }
    render(<ChildrenTable filters={{}} />)

    expect(screen.getByText("Nenhuma criança encontrada")).toBeInTheDocument()
  })
})
