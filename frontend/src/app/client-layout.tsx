"use client"

import { ThemeProvider } from "next-themes"
import { Header } from "@/components/Header"
import { AuthGuard } from "@/components/AuthGuard"
import { Toaster } from "@/components/ui/toaster"
import { usePathname } from "next/navigation"

export function ClientLayout({ children }: { children: React.ReactNode }) {
  const pathname = usePathname()
  const isLogin = pathname === "/login"

  return (
    <ThemeProvider attribute="class" defaultTheme="system" enableSystem>
      <Toaster />

      <Header />

      <AuthGuard>
        <main className={`flex-1 w-full mx-auto ${isLogin ? "" : "p-4 md:p-6 max-w-7xl"}`}>
          {children}
        </main>
      </AuthGuard>
    </ThemeProvider>
  )
}