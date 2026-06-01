import { render, screen, waitFor } from "@testing-library/react"
import userEvent from "@testing-library/user-event"
import { describe, it, expect, vi, beforeEach } from "vitest"
import { ChildDetail } from "@/features/children/components/ChildDetail"

const mockPush = vi.fn()

const mockChildData = {
  id: "c001",
  name: "Ana Clara Mendes",
  age: 6,
  neighborhood: "Rocinha",
  alert_categories: ["health", "education"],
  reviewed: false,
  reviewed_by: null,
  reviewed_at: null,
  notes: "",
  created_at: "2024-01-01T00:00:00Z",
  health: {
    alerts: ["Vacina atrasada"],
    vaccinationsUpToDate: false,
    lastConsultation: "2024-03-15",
  },
  social_assistance: {
    alerts: [],
    cadUnico: true,
    activeBenefit: true,
  },
  education: {
    alerts: [],
    schoolName: "Escola Municipal",
    frequenciaPercent: 95,
  },
}

let mockReturn: {
  child: typeof mockChildData | null
  isLoading: boolean
  isError: boolean
  error: Error | null
}

vi.mock("@/hooks/useChildren", () => ({
  useChild: () => mockReturn,
}))

vi.mock("next/navigation", () => ({
  useRouter: () => ({ push: mockPush }),
}))

beforeEach(() => {
  vi.clearAllMocks()
  mockReturn = {
    child: mockChildData,
    isLoading: false,
    isError: false,
    error: null,
  }
})

describe("ChildDetail", () => {
  it("shows loading skeleton", () => {
    mockReturn.isLoading = true
    const { container } = render(<ChildDetail id="c001" />)

    const skeletons = container.querySelectorAll(".h-8")
    expect(skeletons.length).toBeGreaterThanOrEqual(1)
  })

  it("shows error state", () => {
    mockReturn.isError = true
    mockReturn.error = new Error("Falha ao carregar dados")
    render(<ChildDetail id="c001" />)

    expect(screen.getByText("Falha ao carregar dados")).toBeInTheDocument()
  })

  it("shows not found state", () => {
    mockReturn.child = null
    render(<ChildDetail id="c001" />)

    expect(screen.getByText("Criança não encontrada")).toBeInTheDocument()
  })

  it("renders child basic information", () => {
    render(<ChildDetail id="c001" />)

    expect(screen.getByText("Ana Clara Mendes")).toBeInTheDocument()
    expect(screen.getByText("Rocinha")).toBeInTheDocument()
    expect(screen.getByText("Idade: 6 anos")).toBeInTheDocument()
  })

  it("shows Pendente badge when not reviewed", () => {
    render(<ChildDetail id="c001" />)

    expect(screen.getByText("Pendente")).toBeInTheDocument()
  })

  it("shows Revisado badge when reviewed", () => {
    mockReturn.child = { ...mockChildData, reviewed: true }
    render(<ChildDetail id="c001" />)

    expect(screen.getByText("Revisado")).toBeInTheDocument()
  })

  it("shows tabs for each alert category", () => {
    render(<ChildDetail id="c001" />)

    expect(screen.getByText("Saúde")).toBeInTheDocument()
    expect(screen.getByText("Educação")).toBeInTheDocument()
  })

  it("renders AreaSection for health tab", () => {
    render(<ChildDetail id="c001" />)

    expect(screen.getByText("Vacina atrasada")).toBeInTheDocument()
    expect(screen.getByText(/Última consulta/)).toBeInTheDocument()
  })

  it("renders Education area section when tab is clicked", async () => {
    const user = userEvent.setup()
    render(<ChildDetail id="c001" />)

    await user.click(screen.getByText("Educação"))

    await waitFor(() => {
      expect(screen.getByText(/Escola Municipal/)).toBeInTheDocument()
    })
  })

  it("shows review button when not reviewed", () => {
    render(<ChildDetail id="c001" />)

    expect(
      screen.getByRole("button", { name: /marcar como revisado/i }),
    ).toBeInTheDocument()
  })

  it("hides review button when already reviewed", () => {
    mockReturn.child = { ...mockChildData, reviewed: true }
    render(<ChildDetail id="c001" />)

    expect(
      screen.queryByRole("button", { name: /marcar como revisado/i }),
    ).not.toBeInTheDocument()
  })

  it("shows back button that navigates to children list", async () => {
    const user = userEvent.setup()
    render(<ChildDetail id="c001" />)

    await user.click(screen.getByText("Voltar"))
    expect(mockPush).toHaveBeenCalledWith("/children")
  })

  it("renders SocialAssistance area section", () => {
    mockReturn.child = {
      ...mockChildData,
      alert_categories: ["health", "education", "social_assistance"],
    }
    render(<ChildDetail id="c001" />)

    expect(screen.getByText("Assistência Social")).toBeInTheDocument()
  })

  it("shows reviewed date when available", () => {
    mockReturn.child = {
      ...mockChildData,
      reviewed: true,
      reviewed_at: "2024-06-01T00:00:00Z",
    }
    render(<ChildDetail id="c001" />)

    expect(screen.getByText(/Última revisão/)).toBeInTheDocument()
  })
})
