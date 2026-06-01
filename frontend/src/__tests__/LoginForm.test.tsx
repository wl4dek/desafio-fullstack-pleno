import { render, screen, waitFor, fireEvent } from "@testing-library/react"
import userEvent from "@testing-library/user-event"
import { describe, it, expect, vi, beforeEach } from "vitest"
import { LoginForm } from "@/features/auth/components/LoginForm"

const mockPush = vi.fn()
let mockLogin: { submit: ReturnType<typeof vi.fn>; loading: boolean; error: string | null }

vi.mock("next/navigation", () => ({
  useRouter: () => ({ push: mockPush, replace: vi.fn() }),
}))

vi.mock("@/hooks/useAuth", () => ({
  useLogin: () => mockLogin,
}))

beforeEach(() => {
  vi.clearAllMocks()
  mockLogin = { submit: vi.fn(), loading: false, error: null }
})

describe("LoginForm", () => {
  it("renders form fields and submit button", () => {
    render(<LoginForm />)

    expect(screen.getByText("Acesso ao Painel")).toBeInTheDocument()
    expect(screen.getByLabelText("E-mail")).toBeInTheDocument()
    expect(screen.getByLabelText("Senha")).toBeInTheDocument()
    expect(screen.getByRole("button", { name: /entrar/i })).toBeInTheDocument()
  })

  it("shows validation errors when submitting empty form", async () => {
    const user = userEvent.setup()
    render(<LoginForm />)

    await user.click(screen.getByRole("button", { name: /entrar/i }))

    expect(await screen.findByText("E-mail inválido")).toBeInTheDocument()
    expect(await screen.findByText("Senha obrigatória")).toBeInTheDocument()
  })

  it("shows invalid email error", async () => {
    const user = userEvent.setup()
    render(<LoginForm />)

    await user.type(screen.getByLabelText("E-mail"), "invalido")
    await user.type(screen.getByLabelText("Senha"), "123")

    const form = screen.getByRole("button", { name: /entrar/i }).closest("form")!
    fireEvent.submit(form)

    await waitFor(() => {
      expect(screen.getByText("E-mail inválido")).toBeInTheDocument()
    })
  })

  it("calls submit with email and password on valid form", async () => {
    const user = userEvent.setup()
    render(<LoginForm />)

    await user.type(screen.getByLabelText("E-mail"), "tec@prefeitura.rio")
    await user.type(screen.getByLabelText("Senha"), "senha123")
    await user.click(screen.getByRole("button", { name: /entrar/i }))

    await waitFor(() => {
      expect(mockLogin.submit).toHaveBeenCalledWith(
        "tec@prefeitura.rio",
        "senha123",
      )
    })
  })

  it("shows loading state and disables button", () => {
    mockLogin.loading = true
    render(<LoginForm />)

    expect(screen.getByRole("button", { name: /entrando/i })).toBeInTheDocument()
    expect(screen.getByRole("button")).toBeDisabled()
  })

  it("shows server error message", () => {
    mockLogin.error = "Credenciais inválidas"
    render(<LoginForm />)

    expect(screen.getByText("Credenciais inválidas")).toBeInTheDocument()
  })
})
