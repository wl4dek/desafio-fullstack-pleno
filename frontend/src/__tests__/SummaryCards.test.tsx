import { render, screen } from "@testing-library/react"
import userEvent from "@testing-library/user-event"
import { describe, it, expect, vi, beforeEach } from "vitest"
import { SummaryCards } from "@/features/dashboard/components/SummaryCards"

const mockRefresh = vi.fn()

const mockSummary = {
  total_children: 25,
  reviewed: 14,
  pending_review: 11,
  alerts_by_area: {
    health: 5,
    education: 2,
    social_assistance: 4,
  },
}

let mockReturn: {
  summary: typeof mockSummary | null
  isLoading: boolean
  isError: boolean
  error: Error | null
  refresh: ReturnType<typeof vi.fn>
}

vi.mock("@/hooks/useSummary", () => ({
  useSummary: () => mockReturn,
}))

beforeEach(() => {
  vi.clearAllMocks()
  mockReturn = {
    summary: mockSummary,
    isLoading: false,
    isError: false,
    error: null,
    refresh: mockRefresh,
  }
})

describe("SummaryCards", () => {
  it("shows loading skeleton", () => {
    mockReturn.isLoading = true
    const { container } = render(<SummaryCards />)

    const skeletons = container.querySelectorAll(".h-8")
    expect(skeletons.length).toBe(4)
  })

  it("shows error state with retry button", async () => {
    mockReturn.isError = true
    mockReturn.error = new Error("Erro ao carregar")
    const user = userEvent.setup()
    render(<SummaryCards />)

    expect(screen.getByText("Erro ao carregar indicadores")).toBeInTheDocument()

    await user.click(screen.getByRole("button", { name: /tentar novamente/i }))
    expect(mockReturn.refresh).toHaveBeenCalled()
  })

  it("returns null when summary is not available", () => {
    mockReturn.summary = null
    const { container } = render(<SummaryCards />)

    expect(container.innerHTML).toBe("")
  })

  it("renders summary cards with correct values", () => {
    render(<SummaryCards />)

    expect(screen.getByText("Total de Crianças")).toBeInTheDocument()
    expect(screen.getByText("25")).toBeInTheDocument()
    expect(screen.getByText("Revisadas")).toBeInTheDocument()
    expect(screen.getByText("14")).toBeInTheDocument()
    expect(screen.getByText("Pendentes")).toBeInTheDocument()
    expect(screen.getAllByText("11").length).toBe(2)
    expect(screen.getByText("Com Alerta")).toBeInTheDocument()
  })

  it("renders alerts by area badges", () => {
    render(<SummaryCards />)

    expect(screen.getByText("Saúde: 5")).toBeInTheDocument()
    expect(screen.getByText("Educação: 2")).toBeInTheDocument()
    expect(screen.getByText("Assistência Social: 4")).toBeInTheDocument()
  })

  it("renders correct review percentage", () => {
    render(<SummaryCards />)

    expect(screen.getByText("56.0%")).toBeInTheDocument()
    expect(screen.getByText("44.0%")).toBeInTheDocument()
  })

  it("shows empty alerts message when no alerts", () => {
    mockReturn.summary = {
      total_children: 10,
      reviewed: 5,
      pending_review: 5,
      alerts_by_area: {},
    }
    render(<SummaryCards />)

    expect(screen.getByText("Nenhum alerta registrado")).toBeInTheDocument()
  })
})
