import { test, expect } from '@playwright/test'

test.describe('Navigation', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/dashboard')
  })

  test('should display header with app title', async ({ page }) => {
    await expect(page.getByText('Painel Infantil')).toBeVisible()
  })

  test('should show navigation links on desktop', async ({ page }) => {
    const viewport = page.viewportSize()
    const isDesktop = viewport ? viewport.width >= 1024 : false

    if (isDesktop) {
      const nav = page.locator('nav.hidden.md\\:flex')
      await expect(nav.getByText('Dashboard')).toBeVisible()
      await expect(nav.getByText('Crianças')).toBeVisible()
      await expect(nav.getByText('Estatísticas')).toBeVisible()
    }
  })

  test('should show hamburger menu on mobile and tablet', async ({ page }) => {
    const viewport = page.viewportSize()
    const isMobile = viewport ? viewport.width < 768 : true

    if (isMobile) {
      const menuButton = page.getByLabel('Abrir menu')
      await expect(menuButton).toBeVisible()

      await menuButton.click()
      const sheet = page.locator('[role="dialog"]')
      await expect(sheet).toBeVisible()
      await expect(sheet.getByText('Dashboard')).toBeVisible()
      await expect(sheet.getByText('Crianças')).toBeVisible()
      await expect(sheet.getByText('Estatísticas')).toBeVisible()
    }
  })

  test('should navigate to children page via header link', async ({ page }) => {
    const viewport = page.viewportSize()
    const isDesktop = viewport ? viewport.width >= 1024 : false

    if (isDesktop) {
      await page.getByText('Crianças').first().click()
      await page.waitForURL('/children')
      await expect(page.getByText('Crianças').first()).toBeVisible()
    }
  })

  test('should navigate to statistics page via header link', async ({ page }) => {
    const viewport = page.viewportSize()
    const isDesktop = viewport ? viewport.width >= 1024 : false

    if (isDesktop) {
      await page.getByText('Estatísticas').first().click()
      await page.waitForURL('/statistics')
    }
  })

  test('should highlight active route', async ({ page }) => {
    await page.goto('/children')
    await page.waitForURL('/children')

    const viewport = page.viewportSize()
    const isDesktop = viewport ? viewport.width >= 1024 : false

    if (isDesktop) {
      const activeLink = page.locator('nav.hidden.md\\:flex button[aria-pressed="true"], nav.hidden.md\\:flex [data-state="active"]')
      await expect(activeLink.or(page.locator('nav.hidden.md\\:flex button') .filter({ hasText: 'Crianças' }))).toBeVisible()
    }
  })

  test('should toggle theme', async ({ page }) => {
    const themeButton = page.getByLabel('Alternar tema')
    await expect(themeButton).toBeVisible()
    await themeButton.click()
  })

  test('should logout and redirect to login', async ({ page }) => {
    const viewport = page.viewportSize()
    const isDesktop = viewport ? viewport.width >= 768 : false

    if (isDesktop) {
      await page.getByText('Sair').click()
    } else {
      await page.getByLabel('Abrir menu').click()
      await page.getByText('Sair').click()
    }

    await page.waitForURL('/login')
    await expect(page.getByText('Acesso ao Painel')).toBeVisible({ timeout: 10000 })
  })
})
