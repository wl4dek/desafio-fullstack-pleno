"use client"

import { useEffect, useState } from "react"
import Link from "next/link"
import { usePathname, useRouter } from "next/navigation"
import { useTheme } from "next-themes"

import { useAuthStore } from "@/stores/auth"

import { Button } from "@/components/ui/button"
import {
  Sheet,
  SheetContent,
  SheetTrigger,
  SheetClose,
  SheetTitle,
} from "@/components/ui/sheet"

import {
  LogOut,
  LayoutDashboard,
  Users,
  Sun,
  Moon,
  Menu,
} from "lucide-react"

const BASE_URL = process.env.API_URL || "http://localhost:8080"

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
  const router = useRouter()

  const { theme, setTheme } = useTheme()

  const logout = useAuthStore((s) => s.logout)

  const isAuthenticated = useAuthStore(
    (s) => s.isAuthenticated
  )

  const handleLogout = async () => {
    try {
      await fetch(`${BASE_URL}/auth/session`, {
        method: "DELETE",
        credentials: "include",
      })
    } finally {
      logout()
      router.push("/login")
    }
  }

  if (!mounted) {
    return (
      <header className="sticky top-0 z-40 w-full border-b bg-background">
        <div className="h-14" />
      </header>
    )
  }

  if (!isAuthenticated) return null

  return (
    <header className="sticky top-0 z-40 w-full border-b border-neutral-300 bg-background">
      <div className="flex h-14 items-center px-4 gap-4">
        <Sheet>
          <SheetTrigger asChild className="md:hidden">
            <Button variant="ghost" size="icon" aria-label="Abrir menu">
              <Menu className="h-5 w-5" />
            </Button>
          </SheetTrigger>

          <SheetContent side="left" className="w-64 p-0">
            <SheetTitle className="sr-only">Menu de navegação</SheetTitle>
            <div className="flex flex-col h-full">
              <div className="px-6 pt-6 pb-4 border-b border-neutral-200 dark:border-neutral-700">
                <SheetClose asChild>
                  <Link
                    href="/dashboard"
                    className="font-semibold text-sm"
                  >
                    Painel Infantil
                  </Link>
                </SheetClose>
              </div>

              <nav className="flex flex-col gap-1 p-4 flex-1">
                {navLinks.map((link) => (
                  <SheetClose key={link.href} asChild>
                    <Link href={link.href}>
                      <Button
                        variant={
                          pathname.startsWith(link.href)
                            ? "secondary"
                            : "ghost"
                        }
                        size="sm"
                        className="gap-2 w-full justify-start"
                      >
                        <link.icon className="h-4 w-4" />
                        {link.label}
                      </Button>
                    </Link>
                  </SheetClose>
                ))}
              </nav>

              <div className="p-4 border-t border-neutral-200 dark:border-neutral-700 flex items-center gap-2">
                <Button
                  variant="ghost"
                  size="icon"
                  onClick={() => setTheme(theme === "dark" ? "light" : "dark")}
                  aria-label="Alternar tema"
                >
                  {theme === "dark" ? (
                    <Sun className="h-4 w-4" />
                  ) : (
                    <Moon className="h-4 w-4" />
                  )}
                </Button>

                <SheetClose asChild>
                  <Button
                    variant="ghost"
                    size="sm"
                    onClick={handleLogout}
                    className="gap-2"
                  >
                    <LogOut className="h-4 w-4" />
                    Sair
                  </Button>
                </SheetClose>
              </div>
            </div>
          </SheetContent>
        </Sheet>

        <Link
          href="/dashboard"
          className="font-semibold text-sm mr-4"
        >
          Painel Infantil
        </Link>

        <nav className="hidden md:flex items-center gap-1">
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

        <div className="ml-auto flex items-center gap-1">
          <Button
            variant="ghost"
            size="icon"
            onClick={() => setTheme(theme === "dark" ? "light" : "dark")}
            aria-label="Alternar tema"
          >
            {theme === "dark" ? (
              <Sun className="h-4 w-4" />
            ) : (
              <Moon className="h-4 w-4" />
            )}
          </Button>

          <Button
            variant="ghost"
            size="sm"
            onClick={handleLogout}
            className="hidden md:inline-flex gap-2"
          >
            <LogOut className="h-4 w-4" />
            Sair
          </Button>
        </div>
      </div>
    </header>
  )
}