import { render, screen, waitFor } from "@testing-library/react"
import userEvent from "@testing-library/user-event"
import { describe, it, expect, vi, beforeEach } from "vitest"
import { ChildrenFilters } from "@/features/children/components/ChildrenFilters"

const mockReplace = vi.fn()
let mockSearchParams = new URLSearchParams()

vi.mock("next/navigation", () => ({
  useRouter: () => ({ replace: mockReplace, push: vi.fn() }),
  useSearchParams: () => mockSearchParams,
}))

vi.mock("@/hooks/useChildren", () => ({
  useNeighborhoods: () => ({
    neighborhoods: ["Rocinha", "Maré", "Copacabana"],
    isLoading: false,
    isError: false,
    error: null,
  }),
}))

beforeEach(() => {
  vi.clearAllMocks()
  mockSearchParams = new URLSearchParams()
})

describe("ChildrenFilters", () => {
  it("renders search input with placeholder", () => {
    render(<ChildrenFilters />)

    expect(
      screen.getByPlaceholderText("Buscar por criança..."),
    ).toBeInTheDocument()
  })

  it("renders neighborhood select", () => {
    render(<ChildrenFilters />)

    expect(screen.getByText("Bairro")).toBeInTheDocument()
  })

  it("renders alert type select", () => {
    render(<ChildrenFilters />)

    expect(screen.getByText("Alerta")).toBeInTheDocument()
  })

  it("renders review status select", () => {
    render(<ChildrenFilters />)

    expect(screen.getByText("Revisão")).toBeInTheDocument()
  })

  it("debounces search input and updates filter", async () => {
    const user = userEvent.setup()
    render(<ChildrenFilters />)

    const input = screen.getByPlaceholderText("Buscar por criança...")
    await user.type(input, "Ana")

    await waitFor(
      () => {
        expect(mockReplace).toHaveBeenLastCalledWith(
          expect.stringContaining("childName=Ana"),
        )
      },
      { timeout: 1000 },
    )
  })

  it("resets page to 1 when filter changes", async () => {
    const user = userEvent.setup()
    render(<ChildrenFilters />)

    const input = screen.getByPlaceholderText("Buscar por criança...")
    await user.type(input, "Maria")

    await waitFor(
      () => {
        expect(mockReplace).toHaveBeenCalledWith(
          expect.stringMatching(/page=1/),
        )
      },
      { timeout: 1000 },
    )
  })

  it("removes childName param when search is cleared", async () => {
    const user = userEvent.setup()
    render(<ChildrenFilters />)

    const input = screen.getByPlaceholderText("Buscar por criança...")
    await user.type(input, "A")
    await user.clear(input)

    await waitFor(
      () => {
        expect(mockReplace).toHaveBeenCalledWith(
          expect.stringContaining("/children?"),
        )
      },
      { timeout: 1000 },
    )
  })

  it("renders all select placeholders", async () => {
    render(<ChildrenFilters />)

    expect(screen.getByText("Bairro")).toBeInTheDocument()
    expect(screen.getByText("Alerta")).toBeInTheDocument()
    expect(screen.getByText("Revisão")).toBeInTheDocument()
  })

  it("initializes search from URL params", () => {
    mockSearchParams = new URLSearchParams("childName=Ana")
    render(<ChildrenFilters />)

    const input = screen.getByPlaceholderText(
      "Buscar por criança...",
    ) as HTMLInputElement
    expect(input.value).toBe("Ana")
  })
})
