import { create } from "zustand"

interface AuthState {
  token: string | null
  isAuthenticated: boolean
  hydrated: boolean

  hydrate: () => void
  setToken: (token: string) => void
  logout: () => void
}

export const useAuthStore = create<AuthState>((set) => ({
  token: null,
  isAuthenticated: false,
  hydrated: false,

  hydrate: () => {
    const token = localStorage.getItem("auth_token")

    set({
      token,
      isAuthenticated: !!token,
      hydrated: true,
    })
  },

  setToken: (token: string) => {
    localStorage.setItem("auth_token", token)

    set({
      token,
      isAuthenticated: true,
    })
  },

  logout: () => {
    localStorage.removeItem("auth_token")

    set({
      token: null,
      isAuthenticated: false,
    })
  },
}))