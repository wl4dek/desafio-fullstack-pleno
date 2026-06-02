import { test, expect } from '@playwright/test'

test.describe('Statistics', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/statistics')
  })

  test('should display metric toggle buttons', async ({ page }) => {
    await expect(page.getByText('Saúde').or(page.getByText('Saude'))).toBeVisible()
    await expect(page.getByText('Assistência Social')).toBeVisible()
    await expect(page.getByText('Educação')).toBeVisible()
  })

  test('should have first metric active by default', async ({ page }) => {
    const saudeButton = page.getByRole('button', { name: /Saú?de/i })
    await expect(saudeButton).toBeVisible()
  })

  test('should switch metric when clicking another button', async ({ page }) => {
    const educacaoButton = page.getByRole('button', { name: /Educação/i })
    await educacaoButton.click()

    const assistenciaButton = page.getByRole('button', { name: /Assistência/i })
    await assistenciaButton.click()
  })

  test('should render the map container', async ({ page }) => {
    await page.waitForTimeout(2000)

    const mapContainer = page.locator('.leaflet-container')
    await expect(mapContainer).toBeVisible({ timeout: 15000 })
  })

  test('should display map tiles', async ({ page }) => {
    await page.waitForTimeout(3000)

    const mapTile = page.locator('.leaflet-tile-loaded').first()
    await expect(mapTile).toBeVisible({ timeout: 20000 })
  })
})
