import { api } from "@/lib/api"
import type { LoginResponse } from "@/types"

export function login(email: string, password: string) {
  return api.post<LoginResponse>("/auth/token", { email, password })
}
