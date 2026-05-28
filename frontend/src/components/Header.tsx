"use client"

import { useEffect, useState } from "react"
import Link from "next/link"
import { usePathname } from "next/navigation"

import { useAuthStore } from "@/stores/auth"

import { Button } from "@/components/ui/button"

import {
  LogOut,
  LayoutDashboard,
  Users,
} from "lucide-react"

const navLinks = [
  {
    href: "/dashboard",
    label: "Dashboard",
    icon: LayoutDashboard,
  },
  {
    href: "/children",
    label: "Crianças",
    icon: Users,
  },
]

export function Header() {
  const [mounted, setMounted] = useState(false)

  useEffect(() => {
    setMounted(true)
  }, [])

  const pathname = usePathname()

  const logout = useAuthStore((s) => s.logout)

  const isAuthenticated = useAuthStore(
    (s) => s.isAuthenticated
  )

  // evita hydration mismatch
  if (!mounted) {
    return (
      <header className="sticky top-0 z-40 w-full border-b bg-white dark:bg-neutral-950">
        <div className="h-14" />
      </header>
    )
  }

  if (!isAuthenticated) return null

  return (
    <header className="sticky top-0 z-40 w-full border-b bg-white dark:bg-neutral-950">
      <div className="flex h-14 items-center px-4 gap-4">
        <Link
          href="/dashboard"
          className="font-semibold text-sm mr-4"
        >
          Painel Infantil
        </Link>

        <nav className="flex items-center gap-1">
          {navLinks.map((link) => (
            <Link key={link.href} href={link.href}>
              <Button
                variant={
                  pathname.startsWith(link.href)
                    ? "secondary"
                    : "ghost"
                }
                size="sm"
                className="gap-2"
              >
                <link.icon className="h-4 w-4" />
                {link.label}
              </Button>
            </Link>
          ))}
        </nav>

        <div className="ml-auto">
          <Button
            variant="ghost"
            size="sm"
            onClick={logout}
            className="gap-2"
          >
            <LogOut className="h-4 w-4" />
            Sair
          </Button>
        </div>
      </div>
    </header>
  )
}