import { page } from 'vitest/browser';
import { describe, expect, it } from 'vitest';
import { render } from 'vitest-browser-svelte';
import Page from './+page.svelte';

describe('/+page.svelte', () => {
	it('should render nav', async () => {
		render(Page);

		const nav = page.getByRole('navigation');
		await expect.element(nav).toBeInTheDocument();
	});
});
