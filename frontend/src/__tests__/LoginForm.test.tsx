import { render, screen } from "@testing-library/react"
import { describe, it, expect, vi } from "vitest"
import { LoginForm } from "@/features/auth/components/LoginForm"

const mockPush = vi.fn()

vi.mock("next/navigation", () => ({
  useRouter: () => ({ push: mockPush, replace: vi.fn() }),
}))

vi.mock("@/hooks/useAuth", () => ({
  useLogin: () => ({
    submit: vi.fn(),
    loading: false,
    error: null,
  }),
}))

describe("LoginForm", () => {
  it("renders the login form", () => {
    render(<LoginForm />)
    expect(screen.getByText("Acesso ao Painel")).toBeInTheDocument()
    expect(screen.getByLabelText("E-mail")).toBeInTheDocument()
    expect(screen.getByLabelText("Senha")).toBeInTheDocument()
    expect(screen.getByRole("button", { name: /entrar/i })).toBeInTheDocument()
  })
})
