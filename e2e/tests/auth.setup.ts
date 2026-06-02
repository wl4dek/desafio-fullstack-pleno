import { test as setup, expect } from '@playwright/test'
import path from 'path'
import { fileURLToPath } from 'url'

const __dirname = path.dirname(fileURLToPath(import.meta.url))
const AUTH_FILE = path.resolve(__dirname, '../playwright/.auth/user.json')

setup('authenticate', async ({ page }) => {
  await page.goto('/login')
  await page.getByLabel('E-mail').fill('tecnico@prefeitura.rio')
  await page.getByLabel('Senha').fill('painel@2024')
  await page.getByRole('button', { name: /Entrar/i }).click()

  await page.waitForURL('/dashboard')
  await expect(page.getByText('Total de Crianças')).toBeVisible()

  await page.context().storageState({ path: AUTH_FILE })
})
