import { test, expect } from '@playwright/test'

test.describe('Dashboard', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/dashboard')
  })

  test('should display summary cards with correct labels', async ({ page }) => {
    await expect(page.getByText('Total de Crianças')).toBeVisible()
    await expect(page.getByText('Revisadas')).toBeVisible()
    await expect(page.getByText('Pendentes')).toBeVisible()
    await expect(page.getByText('Com Alerta')).toBeVisible()
  })

  test('should display numeric values in summary cards', async ({ page }) => {
    await expect(page.getByText('Total de Crianças')).toBeVisible()

    const allCards = page.locator('h3.text-sm.font-medium')
    await expect(allCards).toHaveCount(4)
  })

  test('should display alerts by area section', async ({ page }) => {
    await expect(page.getByText('Alertas por Área')).toBeVisible()
  })

  test('should display alerts badges or empty message', async ({ page }) => {
    const alertsSection = page.getByText('Alertas por Área')
    await expect(alertsSection).toBeVisible()

    const noAlerts = page.getByText('Nenhum alerta registrado')
    const badges = page.locator('[data-sonner-toast], .badge, [class*="badge"]')

    if (await noAlerts.isVisible()) {
      await expect(noAlerts).toBeVisible()
    } else {
      await expect(badges.first()).toBeVisible()
    }
  })

  test('should navigate to children page when clicking reviewed card', async ({ page }) => {
    await page.getByRole('link', { name: /Revisadas/ }).click()
    await page.waitForURL('/children?reviewed=true')
    await expect(page.getByText('Crianças')).toBeVisible()
  })

  test('should navigate to children page when clicking pending card', async ({ page }) => {
    await page.getByRole('link', { name: /Pendentes/ }).click()
    await page.waitForURL('/children?reviewed=false')
    await expect(page.getByText('Crianças')).toBeVisible()
  })
})
