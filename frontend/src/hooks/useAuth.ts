"use client"

import { useState } from "react"
import { useRouter } from "next/navigation"
import { login } from "@/services/auth"
import { useAuthStore } from "@/stores/auth"

export function useLogin() {
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const setToken = useAuthStore((s) => s.setToken)
  const router = useRouter()

  const submit = async (email: string, password: string) => {
    setLoading(true)
    setError(null)
    try {
      const res = await login(email, password)
      setToken(res.access_token)
      router.push("/dashboard")
    } catch (err: unknown) {
      const message =
        err instanceof Error ? err.message : "Erro ao fazer login"
      setError(message)
    } finally {
      setLoading(false)
    }
  }

  return { submit, loading, error }
}
