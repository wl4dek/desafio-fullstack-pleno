import { render, screen } from "@testing-library/react"
import { describe, it, expect, vi } from "vitest"
import { SummaryCards } from "@/features/dashboard/components/SummaryCards"

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

const mockRefresh = vi.fn()

vi.mock("@/hooks/useSummary", () => ({
  useSummary: () => ({
    summary: mockSummary,
    isLoading: false,
    isError: false,
    error: null,
    refresh: mockRefresh,
  }),
}))

describe("SummaryCards", () => {
  it("renders summary cards with correct values", () => {
    render(<SummaryCards />)

    expect(screen.getByText("Total de Crianças")).toBeInTheDocument()
    expect(screen.getByText("25")).toBeInTheDocument()
    expect(screen.getByText("Revisadas")).toBeInTheDocument()
    expect(screen.getByText("14")).toBeInTheDocument()
    expect(screen.getByText("Pendentes")).toBeInTheDocument()
    expect(screen.getByText("Com Alerta")).toBeInTheDocument()
  })

  it("renders alerts by area", () => {
    render(<SummaryCards />)

    expect(screen.getByText("Saúde: 5")).toBeInTheDocument()
    expect(screen.getByText("Educação: 2")).toBeInTheDocument()
    expect(screen.getByText("Assistência Social: 4")).toBeInTheDocument()
  })
})
