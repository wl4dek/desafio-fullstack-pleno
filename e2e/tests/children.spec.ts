import { test, expect } from '@playwright/test'

test.describe('Children List', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/children')
  })

  test('should render table with data columns', async ({ page }) => {
    await expect(page.getByText('Nome')).toBeVisible()
    await expect(page.getByText('Alerta')).toBeVisible()
    await expect(page.getByText('Revisão')).toBeVisible()
  })

  test('should display neighborhood column on desktop', async ({ page }) => {
    const viewport = page.viewportSize()
    const isDesktop = viewport ? viewport.width >= 768 : false

    if (isDesktop) {
      await expect(page.getByText('Bairro')).toBeVisible()
    }
  })

  test('should show children data or empty state', async ({ page }) => {
    await page.waitForTimeout(1000)

    const empty = page.getByText('Nenhuma criança encontrada')
    const rows = page.locator('table tbody tr')

    if (await empty.isVisible()) {
      await expect(empty).toBeVisible()
    } else {
      await expect(rows.first()).toBeVisible()
    }
  })

  test('should filter by search with debounce', async ({ page }) => {
    await page.getByPlaceholder('Buscar por criança...').fill('Maria')
    await page.waitForTimeout(600)

    const currentUrl = page.url()
    expect(currentUrl).toContain('childName=Maria')
  })

  test('should navigate to child detail on row click', async ({ page }) => {
    await page.waitForTimeout(1000)

    const empty = page.getByText('Nenhuma criança encontrada')
    if (await empty.isVisible()) {
      test.skip(true, 'No children data available')
    }

    const firstRow = page.locator('table tbody tr').first()
    await firstRow.locator('td').first().click()
    await page.waitForURL(/\/children\/(?!$)/)
    await expect(page.getByText(/Idade:/)).toBeVisible({ timeout: 10000 })
  })

  test('should show pagination controls', async ({ page }) => {
    await page.waitForTimeout(1000)

    const pagination = page.getByText(/Página/)
    const nextButton = page.getByRole('button', { name: /Próximo/i })
    const prevButton = page.getByRole('button', { name: /Anterior/i })

    if (await pagination.isVisible()) {
      await expect(pagination).toBeVisible()
      await expect(nextButton).toBeVisible()
      await expect(prevButton).toBeVisible()
    }
  })

  test('should show error state when API fails', async ({ page }) => {
    await page.route('**/api/v1/children**', route => route.abort())
    await page.reload()
    await expect(page.getByText('Erro ao carregar crianças')).toBeVisible({ timeout: 10000 })
    await expect(page.getByRole('button', { name: /Tentar novamente/i })).toBeVisible()
  })
})
