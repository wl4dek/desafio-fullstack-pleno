import { render, screen, waitFor } from "@testing-library/react"
import userEvent from "@testing-library/user-event"
import { describe, it, expect, vi, beforeEach } from "vitest"
import { ReviewButton } from "@/features/children/components/ReviewButton"

const mockMarkReviewed = vi.fn().mockResolvedValue({ message: "ok" })
const mockMutate = vi.fn()
const mockToast = vi.fn()

vi.mock("@/services/children", () => ({
  markReviewed: (...args: Parameters<typeof mockMarkReviewed>) =>
    mockMarkReviewed(...args),
}))

vi.mock("swr", () => ({
  useSWRConfig: () => ({ mutate: mockMutate }),
}))

vi.mock("@/hooks/use-toast", () => ({
  toast: (...args: Parameters<typeof mockToast>) => mockToast(...args),
}))

beforeEach(() => {
  vi.clearAllMocks()
  mockMarkReviewed.mockResolvedValue({ message: "ok" })
})

describe("ReviewButton", () => {
  it("renders the review button", () => {
    render(<ReviewButton childId="c001" />)

    expect(
      screen.getByRole("button", { name: /marcar como revisado/i }),
    ).toBeInTheDocument()
  })

  it("opens confirmation dialog on click", async () => {
    const user = userEvent.setup()
    render(<ReviewButton childId="c001" />)

    await user.click(
      screen.getByRole("button", { name: /marcar como revisado/i }),
    )

    expect(screen.getByText("Revisão")).toBeInTheDocument()
    expect(
      screen.getByText("O técnico está pronto para registrar a revisão?"),
    ).toBeInTheDocument()
  })

  it("calls markReviewed on confirm", async () => {
    const user = userEvent.setup()
    render(<ReviewButton childId="c001" />)

    await user.click(
      screen.getByRole("button", { name: /marcar como revisado/i }),
    )
    await user.click(screen.getByRole("button", { name: /sim/i }))

    await waitFor(() => {
      expect(mockMarkReviewed).toHaveBeenCalledWith("c001")
    })
  })

  it("mutates SWR cache on success", async () => {
    const user = userEvent.setup()
    render(<ReviewButton childId="c001" />)

    await user.click(
      screen.getByRole("button", { name: /marcar como revisado/i }),
    )
    await user.click(screen.getByRole("button", { name: /sim/i }))

    await waitFor(() => {
      expect(mockMutate).toHaveBeenCalledWith("/children/c001")
      expect(mockMutate).toHaveBeenCalledWith("/summary")
      expect(mockMutate).toHaveBeenCalledWith(expect.any(Function))
    })
  })

  it("shows success toast on review", async () => {
    const user = userEvent.setup()
    render(<ReviewButton childId="c001" />)

    await user.click(
      screen.getByRole("button", { name: /marcar como revisado/i }),
    )
    await user.click(screen.getByRole("button", { name: /sim/i }))

    await waitFor(() => {
      expect(mockToast).toHaveBeenCalledWith(
        expect.objectContaining({
          title: "Revisão registrada",
          variant: "success",
        }),
      )
    })
  })

  it("shows error toast when review fails", async () => {
    mockMarkReviewed.mockRejectedValueOnce(new Error("API Error"))
    const user = userEvent.setup()
    render(<ReviewButton childId="c001" />)

    await user.click(
      screen.getByRole("button", { name: /marcar como revisado/i }),
    )
    await user.click(screen.getByRole("button", { name: /sim/i }))

    await waitFor(() => {
      expect(mockToast).toHaveBeenCalledWith(
        expect.objectContaining({
          title: "Erro",
          variant: "destructive",
        }),
      )
    })
  })

  it("shows loading state while reviewing", async () => {
    mockMarkReviewed.mockImplementationOnce(
      () => new Promise(() => {}),
    )
    const user = userEvent.setup()
    render(<ReviewButton childId="c001" />)

    await user.click(
      screen.getByRole("button", { name: /marcar como revisado/i }),
    )
    await user.click(screen.getByRole("button", { name: /sim/i }))

    expect(await screen.findByText("Registrando...")).toBeInTheDocument()
  })
})
