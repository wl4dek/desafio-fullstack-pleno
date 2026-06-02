import { test, expect } from '@playwright/test'

test.describe('Login', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/login')
  })

  test('should display login form with title, fields and button', async ({ page }) => {
    await expect(page.getByText('Acesso ao Painel')).toBeVisible()
    await expect(page.getByText('Informe suas credenciais')).toBeVisible()
    await expect(page.getByLabel('E-mail')).toBeVisible()
    await expect(page.getByLabel('Senha')).toBeVisible()
    await expect(page.getByRole('button', { name: /Entrar/i })).toBeVisible()
  })

  test('should show validation errors when submitting empty form', async ({ page }) => {
    await page.getByRole('button', { name: /Entrar/i }).click()
    await expect(page.getByText('E-mail inválido')).toBeVisible()
    await expect(page.getByText('Senha obrigatória')).toBeVisible()
  })

  test('should show validation error for invalid email format', async ({ page }) => {
    await page.getByLabel('E-mail').fill('email-invalido')
    await page.getByLabel('Senha').fill('password')
    await page.getByRole('button', { name: /Entrar/i }).click()
    await expect(page.getByText('E-mail inválido')).toBeVisible()
  })

  test('should show server error for wrong credentials', async ({ page }) => {
    await page.getByLabel('E-mail').fill('errado@email.com')
    await page.getByLabel('Senha').fill('senha-errada')
    await page.getByRole('button', { name: /Entrar/i }).click()

    await expect(page.locator('p.text-red-500.text-center')).toBeVisible({ timeout: 10000 })
    await expect(page).toHaveURL('/login')
  })

  test('should login successfully and redirect to dashboard', async ({ page }) => {
    await page.getByLabel('E-mail').fill('tecnico@prefeitura.rio')
    await page.getByLabel('Senha').fill('painel@2024')
    await page.getByRole('button', { name: /Entrar/i }).click()

    await page.waitForURL('/dashboard')
    await expect(page.getByText('Total de Crianças')).toBeVisible({ timeout: 10000 })
  })
})
