import { test, expect } from '@playwright/test'

test.describe('Child Detail', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/children')
    await page.waitForTimeout(1000)

    const empty = page.getByText('Nenhuma criança encontrada')
    if (await empty.isVisible()) {
      test.skip(true, 'No children data available')
      return
    }

    const firstRow = page.locator('table tbody tr').first()
    await firstRow.locator('td').first().click()
    await page.waitForURL(/\/children\/(?!$)/)
  })

  test('should display child name, neighborhood and age', async ({ page }) => {
    await expect(page.getByText(/Idade:/)).toBeVisible({ timeout: 10000 })
    const neighborhood = page.locator('p.text-sm.text-neutral-500')
    await expect(neighborhood).toBeVisible()
  })

  test('should show review status badge', async ({ page }) => {
    await page.waitForTimeout(1000)

    const revisado = page.getByText('Revisado')
    const pendente = page.getByText('Pendente')

    const hasRevisado = await revisado.isVisible().catch(() => false)
    const hasPendente = await pendente.isVisible().catch(() => false)

    expect(hasRevisado || hasPendente).toBeTruthy()
  })

  test('should render tabs for alert categories', async ({ page }) => {
    await page.waitForTimeout(1000)

    const tabs = page.locator('[role="tab"]')
    const count = await tabs.count()

    if (count > 0) {
      const firstTab = tabs.first()
      await expect(firstTab).toBeVisible()
      await firstTab.click()
    }
  })

  test('should display area section content when clicking a tab', async ({ page }) => {
    await page.waitForTimeout(1000)

    const tabs = page.locator('[role="tab"]')
    const count = await tabs.count()

    if (count > 0) {
      await tabs.first().click()
      const tabPanel = page.locator('[role="tabpanel"]')
      await expect(tabPanel.first()).toBeVisible()
    }
  })

  test('should show review button for unreviewed children', async ({ page }) => {
    await page.waitForTimeout(1000)

    const pendente = page.getByText('Pendente')
    const reviewButton = page.getByRole('button', { name: /Marcar como Revisado/i })

    if (await pendente.isVisible()) {
      await expect(reviewButton).toBeVisible()
    }
  })

  test('should complete review flow: dialog confirm and toast', async ({ page }) => {
    await page.waitForTimeout(1000)

    const pendente = page.getByText('Pendente')
    if (!(await pendente.isVisible())) {
      test.skip(true, 'Child already reviewed')
      return
    }

    const reviewButton = page.getByRole('button', { name: /Marcar como Revisado/i })
    await reviewButton.click()

    const dialog = page.getByRole('alertdialog')
    await expect(dialog).toBeVisible()
    await expect(page.getByText('O técnico está pronto para registrar a revisão?')).toBeVisible()

    await page.getByRole('button', { name: /Sim/i }).click()

    await expect(page.getByText('Revisão registrada')).toBeVisible({ timeout: 10000 })
  })

  test('should navigate back to children list', async ({ page }) => {
    await page.getByRole('button', { name: /Voltar/i }).first().click()
    await page.waitForURL('/children')
  })
})
