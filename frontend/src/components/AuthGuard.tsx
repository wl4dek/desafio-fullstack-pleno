"use client"

import { useEffect, useState } from "react"
import { useRouter, usePathname } from "next/navigation"
import { useAuthStore } from "@/stores/auth"

const BASE_URL = process.env.API_URL || "http://localhost:8080"
const publicRoutes = ["/login"]

export function AuthGuard({ children }: { children: React.ReactNode }) {
  const isAuthenticated = useAuthStore((s) => s.isAuthenticated)
  const pathname = usePathname()

  if (isAuthenticated || publicRoutes.includes(pathname)) {
    return <>{children}</>
  }

  return <SessionGuard>{children}</SessionGuard>
}

function SessionGuard({ children }: { children: React.ReactNode }) {
  const [loading, setLoading] = useState(true)
  const setToken = useAuthStore((s) => s.setToken)
  const router = useRouter()

  useEffect(() => {
    const controller = new AbortController()

    fetch(`${BASE_URL}/auth/session`, {
      credentials: "include",
      signal: controller.signal,
    })
      .then((res) => {
        if (!res.ok) throw new Error("no session")
        return res.json() as Promise<{ access_token: string }>
      })
      .then((data) => {
        setToken(data.access_token)
        setLoading(false)
      })
      .catch(() => {
        router.replace("/login")
      })

    return () => controller.abort()
  }, [setToken, router])

  if (loading) return null

  return <>{children}</>
}
