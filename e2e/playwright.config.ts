import { defineConfig } from '@playwright/test';

export default defineConfig({
  testDir: './tests',
  fullyParallel: true,
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 2 : 0,
  workers: process.env.CI ? 1 : undefined,
  reporter: 'html',
  use: {
    baseURL: process.env.BASE_URL || 'http://localhost:3000',
    screenshot: 'only-on-failure',
    video: 'retain-on-failure',
    trace: 'retain-on-failure',
  },

  reporters: [
    ['html', { outputFolder: 'playwright/reports', open: 'never' }],
  ],

  projects: [
    {
      name: 'setup',
      testMatch: 'auth.setup.ts',
    },
    {
      name: 'Mobile (375px)',
      use: { viewport: { width: 375, height: 812 } },
      dependencies: ['setup'],
      storageState: 'playwright/.auth/user.json',
    },
    {
      name: 'Tablet (640px)',
      use: { viewport: { width: 640, height: 900 } },
      dependencies: ['setup'],
      storageState: 'playwright/.auth/user.json',
    },
    {
      name: 'Desktop (1024px)',
      use: { viewport: { width: 1024, height: 768 } },
      dependencies: ['setup'],
      storageState: 'playwright/.auth/user.json',
    },
    {
      name: 'Desktop (1440px)',
      use: { viewport: { width: 1440, height: 900 } },
      dependencies: ['setup'],
      storageState: 'playwright/.auth/user.json',
    },
  ],
});
