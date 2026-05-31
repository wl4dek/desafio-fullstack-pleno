"use client"

import { ThemeProvider } from "next-themes"
import { Header } from "@/components/Header"
import { AuthGuard } from "@/components/AuthGuard"
import { Toaster } from "@/components/ui/toaster"

export function ClientLayout({ children }: { children: React.ReactNode }) {
  return (
    <ThemeProvider attribute="class" defaultTheme="system" enableSystem>
      <Toaster />

      <Header />

      <AuthGuard>
        <main className="flex-1 p-4 md:p-6 max-w-7xl w-full mx-auto">
          {children}
        </main>
      </AuthGuard>
    </ThemeProvider>
  )
}