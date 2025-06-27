import { test, expect } from '@playwright/test';

test.describe('Regex Token Generation', () => {
  test('User can generate tokens from regex rules', async ({ page }) => {
    await page.addInitScript(() => {
      localStorage.setItem('user_id', '6833695df0819feda85033e8');
    });

    await page.goto('http://localhost:5173/main-workspace');

    // 1. Click the toolbox button to create a Source Code node
    await page.getByRole('button', { name: /source code.*start here/i }).click();

    // 2. Click the Source Code node in the canvas (SVG group)
    await page.locator('g[id^="N-source"]').first().click();

    // 3. Fill in the source code and confirm
    await page.getByPlaceholder('Paste or type your source code here…').fill('int blue = 13 + 22;');
    await page.getByRole('button', { name: 'Confirm Code' }).click();

    // 4. Click the "Lexer" node in the canvas
    await page.getByRole('button', { name: /lexer/i }).click();
    await page.locator('g[id^="N-lexer"]').first().click();

    // 5. Select "Regular Expression" mode
    await page.getByRole('button', { name: /regular expression/i }).click();

    // 6. Fill in at least one regex rule
    await page.getByPlaceholder('Enter type...').fill('integer');
    await page.getByPlaceholder('Enter regex pattern...').fill('[0-9]+');

    // 7. Submit the regex rules
    await page.getByRole('button', { name: 'Submit' }).click();

    // 8. Wait for the "Generate Tokens" button to appear and click it
    await page.getByRole('button', { name: 'Generate Tokens' }).click();

    // 9. Assert that tokens are displayed (adjust selector as needed)
    await expect(page.getByText(/Tokens generated successfully!/i)).toBeVisible();
  });
});