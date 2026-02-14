import { page } from 'vitest/browser';
import { describe, expect, it } from 'vitest';
import { render } from 'vitest-browser-svelte';
import ArticleCard from './ArticleCard.svelte';
import type { Article as ArticleType } from '$lib/server/types';

describe('ArticleCard.svelte', () => {
	it('should render article with all fields', async () => {
		const article: ArticleType = {
			account: 'test-account',
			id: '1',
			url: 'https://example.com/article',
			createdAt: '2024-01-01T00:00:00Z',
			title: 'Test Article',
			content: 'Test content',
			author: 'Test Author',
			siteName: 'Example Site',
			excerpt: 'Test excerpt',
			imageUrl: 'https://example.com/image.jpg',
			wordCount: 500,
			readingTimeMinutes: 2,
			deliveryStatus: 'pending'
		};

		render(ArticleCard, { article });

		const titleLink = page.getByRole('link', { name: 'Test Article' });
		await expect.element(titleLink).toBeInTheDocument();

		const excerpt = page.getByText('Test excerpt');
		await expect.element(excerpt).toBeInTheDocument();

		const author = page.getByText('by Test Author');
		await expect.element(author).toBeInTheDocument();

		const source = page.getByText('source: Example Site');
		await expect.element(source).toBeInTheDocument();

		const wordCount = page.getByText('500 words');
		await expect.element(wordCount).toBeInTheDocument();

		const readingTime = page.getByText('2 min read');
		await expect.element(readingTime).toBeInTheDocument();

		const status = page.getByText('status: pending');
		await expect.element(status).toBeInTheDocument();
	});

	it('should render article with minimal fields', async () => {
		const article: ArticleType = {
			account: 'test-account',
			id: '1',
			url: 'https://example.com/article',
			createdAt: '2024-01-01T00:00:00Z'
		};

		render(ArticleCard, { article });

		const heading = page.getByRole('heading');
		await expect.element(heading).toBeInTheDocument();

		const originalLink = page.getByRole('link', { name: 'https://example.com/article' }).nth(1);
		await expect.element(originalLink).toHaveAttribute('target', '_blank');
		await expect.element(originalLink).toHaveAttribute('rel', 'external');
	});

	it('should render article error', async () => {
		const article: ArticleType = {
			account: 'test-account',
			id: '1',
			url: 'https://example.com/article',
			createdAt: '2024-01-01T00:00:00Z',
			error: 'Failed to fetch content'
		};

		render(ArticleCard, { article });

		const error = page.getByText('error: Failed to fetch content');
		await expect.element(error).toBeInTheDocument();
	});
});
