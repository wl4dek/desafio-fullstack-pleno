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

  return <SessionGuard pathname={pathname}>{children}</SessionGuard>
}

function SessionGuard({
  children,
  pathname,
}: {
  children: React.ReactNode
  pathname: string
}) {
  const [loading, setLoading] = useState(true)
  const setToken = useAuthStore((s) => s.setToken)
  const router = useRouter()

  useEffect(() => {
    fetch(`${BASE_URL}/auth/session`, {
      credentials: "include",
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
        setLoading(false)
        router.replace("/login")
      })
  }, [setToken, router, pathname])

  if (loading) return null

  return <>{children}</>
}
