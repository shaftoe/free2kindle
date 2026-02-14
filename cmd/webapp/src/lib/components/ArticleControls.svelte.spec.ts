import { page } from 'vitest/browser';
import { describe, expect, it, vi } from 'vitest';
import { render } from 'vitest-browser-svelte';
import ArticleControls from './ArticleControls.svelte';
import type { Article as ArticleType } from '$lib/server/types';

describe('ArticleControls.svelte', () => {
	it('should render delete button', async () => {
		const article: ArticleType = {
			account: 'test-account',
			id: '1',
			url: 'https://example.com/article',
			createdAt: '2024-01-01T00:00:00Z',
			title: 'Test Article'
		};

		render(ArticleControls, { article });

		const deleteButton = page.getByRole('button', { name: 'Delete' });
		await expect.element(deleteButton).toBeInTheDocument();
	});

	it('should render form with correct action and hidden input', async () => {
		const article: ArticleType = {
			account: 'test-account',
			id: 'test-id-123',
			url: 'https://example.com/article',
			createdAt: '2024-01-01T00:00:00Z',
			title: 'Test Article'
		};

		const { container } = render(ArticleControls, { article });

		const form = container.querySelector('form');
		expect(form).toBeDefined();
		expect(form?.getAttribute('method')).toBe('POST');
		expect(form?.getAttribute('action')).toBe('?/delete');

		const hiddenInput = container.querySelector('input[type="hidden"][name="id"]');
		expect(hiddenInput).toBeDefined();
		expect(hiddenInput?.getAttribute('value')).toBe('test-id-123');
	});

	it('should show confirmation dialog before submit', async () => {
		const article: ArticleType = {
			account: 'test-account',
			id: '1',
			url: 'https://example.com/article',
			createdAt: '2024-01-01T00:00:00Z',
			title: 'Test Article'
		};

		const confirmSpy = vi.spyOn(window, 'confirm').mockReturnValue(false);

		render(ArticleControls, { article });

		const deleteButton = page.getByRole('button', { name: 'Delete' });
		await deleteButton.click();

		expect(confirmSpy).toHaveBeenCalledWith('Are you sure you want to delete this article?');

		confirmSpy.mockRestore();
	});
});
