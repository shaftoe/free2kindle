import { page } from 'vitest/browser';
import { describe, expect, it } from 'vitest';
import { render } from 'vitest-browser-svelte';
import BackNav from './BackNav.svelte';

describe('BackNav.svelte', () => {
	it('should render back link with default href', async () => {
		render(BackNav);

		const link = page.getByRole('link', { name: '← Back' });
		await expect.element(link).toBeInTheDocument();
		await expect.element(link).toHaveAttribute('href', '/');
	});

	it('should render back link with custom href', async () => {
		render(BackNav, { href: '/articles' });

		const link = page.getByRole('link', { name: '← Back' });
		await expect.element(link).toBeInTheDocument();
		await expect.element(link).toHaveAttribute('href', '/articles');
	});
});
