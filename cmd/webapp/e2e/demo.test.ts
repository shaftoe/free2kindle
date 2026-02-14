import { expect, test } from '@playwright/test';

test('home page has expected nav', async ({ page }) => {
	await page.goto('/');
	await expect(page.getByRole('navigation').filter({ hasText: 'My List' })).toBeVisible();
});
