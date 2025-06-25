import { test, expect } from '@playwright/test';

test('homepage loads successfully', async ({ page }) => {
  await page.goto('/');
  
  // Check if the page loads
  await expect(page).toHaveTitle(/Stockle/);
  
  // Add more specific tests as the application develops
  // await expect(page.getByRole('heading', { name: 'Welcome to Stockle' })).toBeVisible();
});

test('navigation works', async ({ page }) => {
  await page.goto('/');
  
  // Test navigation functionality once implemented
  // await page.getByRole('button', { name: 'Sign In' }).click();
  // await expect(page).toHaveURL(/.*auth/);
});