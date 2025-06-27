import { test, expect } from '@playwright/test';

test.describe('Auth Page', () => {
  test('User can register', async ({ page }) => {
    await page.goto('http://localhost:5173/auth-page');
    await page.getByRole('button', { name: /^Register$/ }).click();
    await expect(page.getByLabel('Email', { exact: true })).toBeVisible();

    const unique = Date.now();
    await page.getByLabel('Username', { exact: true }).fill(`testuser${unique}`);
    await page.getByLabel('Email', { exact: true }).fill(`testuser${unique}@example.com`);
    await page.getByLabel('Password', { exact: true }).fill('testpassword123');
    await page.getByLabel('Confirm Password', { exact: true }).fill('testpassword123');
    await page.locator('button[type="submit"][aria-label="Register"]').click();

    // Wait for either the toast or the login tab to be active again
    await expect(page.getByRole('button', { name: /^Login$/ })).toHaveClass(/active/);
  });

  test('User can login', async ({ page }) => {
    await page.goto('http://localhost:5173/auth-page');
    await page.getByRole('button', { name: /^Login$/ }).click();
    await expect(page.getByLabel('Username', { exact: true })).toBeVisible();

    // Use a known user or create one in a beforeAll hook
    await page.getByLabel('Username', { exact: true }).fill('your_existing_user');
    await page.getByLabel('Password', { exact: true }).fill('your_existing_password');
    await page.locator('button[type="submit"][aria-label="Login"]').click();

    await expect(page.getByText(/Login successful!/i)).toBeVisible({ timeout: 7000 });
  });
});