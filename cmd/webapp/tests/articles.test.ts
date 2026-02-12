import { describe, it, expect } from 'vitest';
import type { Article, ArticlesResponse } from '../src/lib/types';

describe('Articles API endpoint types', () => {
  it('should have correct Article interface', () => {
    const article: Article = {
      id: '123',
      url: 'https://example.com',
      title: 'Test Article',
      createdAt: '2024-01-01T00:00:00Z',
    };

    expect(article.id).toBe('123');
    expect(article.url).toBe('https://example.com');
    expect(article.title).toBe('Test Article');
    expect(article.createdAt).toBe('2024-01-01T00:00:00Z');
  });

  it('should have correct ArticlesResponse interface', () => {
    const response: ArticlesResponse = {
      articles: [],
      page: 1,
      pageSize: 20,
      total: 0,
      hasMore: false,
    };

    expect(response.page).toBe(1);
    expect(response.pageSize).toBe(20);
    expect(response.total).toBe(0);
    expect(response.hasMore).toBe(false);
  });
});
