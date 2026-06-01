import { describe, it, expect } from "vitest"
import { loginSchema } from "@/schemas"

describe("loginSchema", () => {
  it("accepts valid email and password", () => {
    const result = loginSchema.safeParse({
      email: "tecnico@prefeitura.rio",
      password: "senha123",
    })
    expect(result.success).toBe(true)
  })

  it("rejects invalid email", () => {
    const result = loginSchema.safeParse({
      email: "invalido",
      password: "senha123",
    })
    expect(result.success).toBe(false)
    if (!result.success) {
      expect(result.error.issues[0].message).toBe("E-mail inválido")
    }
  })

  it("rejects empty password", () => {
    const result = loginSchema.safeParse({
      email: "tecnico@prefeitura.rio",
      password: "",
    })
    expect(result.success).toBe(false)
    if (!result.success) {
      expect(result.error.issues[0].message).toBe("Senha obrigatória")
    }
  })

  it("rejects empty email", () => {
    const result = loginSchema.safeParse({
      email: "",
      password: "senha123",
    })
    expect(result.success).toBe(false)
  })

  it("rejects missing fields", () => {
    const result = loginSchema.safeParse({})
    expect(result.success).toBe(false)
  })
})
