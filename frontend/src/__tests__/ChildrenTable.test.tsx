import { render, screen } from "@testing-library/react"
import userEvent from "@testing-library/user-event"
import { describe, it, expect, vi, beforeEach } from "vitest"
import { ChildrenTable } from "@/features/children/components/ChildrenTable"

const mockPush = vi.fn()

const mockData = {
  data: [
    {
      id: "c001",
      name: "Ana Clara Mendes",
      age: 6,
      neighborhood: "Rocinha",
      alert_categories: ["health"],
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
      alert_categories: [],
      reviewed: true,
      reviewed_by: "tecnico@prefeitura.rio",
      reviewed_at: "2024-06-01T00:00:00Z",
      notes: "",
      created_at: "2024-01-02T00:00:00Z",
    },
  ],
  pagination: { page: 1, per_page: 10, total: 2, total_pages: 1 },
}

let mockReturn: {
  data: typeof mockData | { data: []; pagination: { page: number; per_page: number; total: number; total_pages: number } }
  isLoading: boolean
  isError: boolean
  error: Error | null
  refresh: ReturnType<typeof vi.fn>
}

vi.mock("@/hooks/useChildren", () => ({
  useChildren: () => mockReturn,
}))

vi.mock("next/navigation", () => ({
  useRouter: () => ({ push: mockPush, replace: vi.fn() }),
  useSearchParams: () => new URLSearchParams(),
}))

beforeEach(() => {
  vi.clearAllMocks()
  mockReturn = {
    data: mockData,
    isLoading: false,
    isError: false,
    error: null,
    refresh: vi.fn(),
  }
})

describe("ChildrenTable", () => {
  it("shows loading skeleton", () => {
    mockReturn.isLoading = true
    const { container } = render(<ChildrenTable filters={{}} />)

    const skeletons = container.querySelectorAll(".h-12")
    expect(skeletons.length).toBe(5)
  })

  it("shows error state with retry button", async () => {
    mockReturn.isError = true
    mockReturn.error = new Error("Erro ao carregar")
    const user = userEvent.setup()
    render(<ChildrenTable filters={{}} />)

    expect(screen.getByText("Erro ao carregar crianças")).toBeInTheDocument()

    await user.click(screen.getByRole("button", { name: /tentar novamente/i }))
    expect(mockReturn.refresh).toHaveBeenCalled()
  })

  it("shows empty state when no data", () => {
    mockReturn.data = {
      data: [],
      pagination: { page: 1, per_page: 10, total: 0, total_pages: 0 },
    }
    render(<ChildrenTable filters={{}} />)

    expect(screen.getByText("Nenhuma criança encontrada")).toBeInTheDocument()
  })

  it("renders children data correctly", () => {
    render(<ChildrenTable filters={{}} />)

    expect(screen.getByText("Ana Clara Mendes")).toBeInTheDocument()
    expect(screen.getByText("Lucas Ferreira")).toBeInTheDocument()
    expect(screen.getByText("Rocinha")).toBeInTheDocument()
    expect(screen.getByText("Maré")).toBeInTheDocument()
  })

  it("shows alert badges for children with alerts", () => {
    render(<ChildrenTable filters={{}} />)

    const alertBadge = screen.getByText("Saúde")
    expect(alertBadge).toBeInTheDocument()
  })

  it("shows OK badge for children without alerts", () => {
    render(<ChildrenTable filters={{}} />)

    expect(screen.getByText("OK")).toBeInTheDocument()
  })

  it("shows review status badges", () => {
    render(<ChildrenTable filters={{}} />)

    expect(screen.getByText("Pendente")).toBeInTheDocument()
    expect(screen.getByText("Revisado")).toBeInTheDocument()
  })

  it("navigates to child detail on row click", async () => {
    const user = userEvent.setup()
    render(<ChildrenTable filters={{}} />)

    await user.click(screen.getByText("Ana Clara Mendes"))
    expect(mockPush).toHaveBeenCalledWith("/children/c001")
  })
})
