import { page } from 'vitest/browser';
import { describe, expect, it } from 'vitest';
import { render } from 'vitest-browser-svelte';
import Page from './+page.svelte';

describe('/+page.svelte', () => {
	it('should render article list', async () => {
		render(Page, {
			data: {
				articles: [],
				total: 0,
				page: 1,
				pageSize: 10,
				hasMore: false,
				error: undefined
			}
		});

		const heading = page.getByRole('heading', { name: 'articles' });
		await expect.element(heading).toBeInTheDocument();

		const noArticles = page.getByText('no articles yet');
		await expect.element(noArticles).toBeInTheDocument();
	});
});
