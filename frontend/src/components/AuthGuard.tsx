"use client"

import { useEffect } from "react"
import { useRouter, usePathname } from "next/navigation"
import { useAuthStore } from "@/stores/auth"

const publicRoutes = ["/login"]

export function AuthGuard({ children }: { children: React.ReactNode }) {
  const isAuthenticated = useAuthStore((s) => s.isAuthenticated)
  const router = useRouter()
  const pathname = usePathname()

  useEffect(() => {
    if (!isAuthenticated && !publicRoutes.includes(pathname)) {
      router.replace("/login")
    }
  }, [isAuthenticated, pathname, router])

  if (!isAuthenticated && !publicRoutes.includes(pathname)) {
    return null
  }

  return <>{children}</>
}
