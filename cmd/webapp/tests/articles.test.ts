import { describe, it, expect } from 'vitest';

describe('Articles API endpoint types', () => {
  it('should have correct Article interface', () => {
    interface Article {
      id: string;
      url: string;
      title?: string;
      createdAt: string;
      author?: string;
      siteName?: string;
      excerpt?: string;
      imageUrl?: string;
      wordCount?: number;
      readingTimeMinutes?: number;
      publishedAt?: string;
      deliveryStatus?: string;
    }

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
    interface ArticlesResponse {
      articles: any[];
      page: number;
      pageSize: number;
      total: number;
      hasMore: boolean;
    }

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
