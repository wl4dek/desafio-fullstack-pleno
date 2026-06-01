import { render, screen, waitFor } from "@testing-library/react"
import userEvent from "@testing-library/user-event"
import { describe, it, expect, vi, beforeEach } from "vitest"
import { Header } from "@/components/Header"

let mockIsAuthenticated = false
let mockPathname = "/dashboard"
let mockTheme = "light"
const mockSetTheme = vi.fn()
const mockLogout = vi.fn()
const mockPush = vi.fn()

vi.mock("next/navigation", () => ({
  usePathname: () => mockPathname,
  useRouter: () => ({ push: mockPush }),
}))

vi.mock("next-themes", () => ({
  useTheme: () => ({ theme: mockTheme, setTheme: mockSetTheme }),
}))

vi.mock("@/stores/auth", () => ({
  useAuthStore: (selector: (s: any) => any) => {
    const store = {
      isAuthenticated: mockIsAuthenticated,
      logout: mockLogout,
    }
    return selector(store)
  },
}))

beforeEach(() => {
  vi.clearAllMocks()
  mockIsAuthenticated = false
  mockPathname = "/dashboard"
  mockTheme = "light"
  global.fetch = vi.fn().mockResolvedValue({ ok: true })
})

describe("Header", () => {
  it("returns null when not authenticated after mount", async () => {
    render(<Header />)

    await waitFor(() => {
      expect(
        screen.queryByRole("link", { name: /dashboard/i }),
      ).not.toBeInTheDocument()
    })
  })

  it("renders navigation links when authenticated", async () => {
    mockIsAuthenticated = true
    render(<Header />)

    expect(await screen.findByText("Dashboard")).toBeInTheDocument()
    expect(screen.getByText("Crianças")).toBeInTheDocument()
    expect(screen.getByText("Estatísticas")).toBeInTheDocument()
  })

  it("highlights active route", async () => {
    mockIsAuthenticated = true
    mockPathname = "/children"
    render(<Header />)

    await waitFor(() => {
      const childrenLink = screen.getByText("Crianças")
      expect(childrenLink).toBeInTheDocument()
    })
  })

  it("calls logout and navigates to login", async () => {
    mockIsAuthenticated = true
    const user = userEvent.setup()
    render(<Header />)

    const logoutButton = await screen.findByRole("button", { name: /sair/i })
    await user.click(logoutButton)

    await waitFor(() => {
      expect(global.fetch).toHaveBeenCalledWith(
        expect.stringContaining("/auth/session"),
        expect.objectContaining({ method: "DELETE" }),
      )
    })

    expect(mockLogout).toHaveBeenCalled()
    expect(mockPush).toHaveBeenCalledWith("/login")
  })

  it("renders mobile menu button", async () => {
    mockIsAuthenticated = true
    render(<Header />)

    expect(
      await screen.findByRole("button", { name: /abrir menu/i }),
    ).toBeInTheDocument()
  })

  it("renders theme toggle button", async () => {
    mockIsAuthenticated = true
    render(<Header />)

    expect(
      await screen.findByRole("button", { name: /alternar tema/i }),
    ).toBeInTheDocument()
  })

  it("toggles theme on click", async () => {
    mockIsAuthenticated = true
    const user = userEvent.setup()
    render(<Header />)

    const themeButton = await screen.findByRole("button", {
      name: /alternar tema/i,
    })
    await user.click(themeButton)

    expect(mockSetTheme).toHaveBeenCalledWith("dark")
  })
})
