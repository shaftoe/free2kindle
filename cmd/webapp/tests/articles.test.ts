import { describe, it, expect, beforeEach, afterEach, vi } from 'vitest';
import type { Article, ArticlesResponse } from '../src/lib/types';
import { ApiClient, joinUrl } from '../src/lib/services/apiClient';
import { getToken, setToken, clearToken } from '../src/lib/stores/token';

describe('joinUrl', () => {
  it('should join base URL and path correctly with trailing slash in base', () => {
    expect(joinUrl('https://api.example.com/', '/v1/articles')).toBe('https://api.example.com/v1/articles');
  });

  it('should join base URL and path correctly without trailing slash in base', () => {
    expect(joinUrl('https://api.example.com', '/v1/articles')).toBe('https://api.example.com/v1/articles');
  });

  it('should join base URL and path correctly without leading slash in path', () => {
    expect(joinUrl('https://api.example.com', 'v1/articles')).toBe('https://api.example.com/v1/articles');
  });

  it('should join base URL and path correctly with neither trailing nor leading slashes', () => {
    expect(joinUrl('https://api.example.com', 'v1/articles')).toBe('https://api.example.com/v1/articles');
  });

  it('should handle multiple trailing slashes in base URL', () => {
    expect(joinUrl('https://api.example.com///', '/v1/articles')).toBe('https://api.example.com/v1/articles');
  });

  it('should handle multiple leading slashes in path', () => {
    expect(joinUrl('https://api.example.com/', '///v1/articles')).toBe('https://api.example.com/v1/articles');
  });

  it('should handle query strings in path', () => {
    expect(joinUrl('https://api.example.com/', '/v1/articles?page=1')).toBe('https://api.example.com/v1/articles?page=1');
  });
});

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

describe('Token store', () => {
  beforeEach(() => {
    clearToken();
  });

  afterEach(() => {
    clearToken();
  });

  it('should set and get token', () => {
    setToken('test-token');
    expect(getToken()).toBe('test-token');
  });

  it('should clear token', () => {
    setToken('test-token');
    clearToken();
    expect(getToken()).toBe(null);
  });

  it('should get null when no token is set', () => {
    expect(getToken()).toBe(null);
  });
});

describe('ApiClient', () => {
  const fetchSpy = vi.fn();

  beforeEach(() => {
    clearToken();
    fetchSpy.mockClear();
    global.fetch = fetchSpy;
  });

  afterEach(() => {
    clearToken();
  });

  it('should include Authorization header when token is set', async () => {
    setToken('test-token');
    const client = new ApiClient('http://localhost:8080');

    fetchSpy.mockResolvedValue({
      ok: true,
      json: () => Promise.resolve({ data: 'test' }),
    } as Response);

    await client.get('/test');

    expect(fetchSpy).toHaveBeenCalledWith(
      'http://localhost:8080/test',
      expect.objectContaining({
        headers: expect.objectContaining({
          'Content-Type': 'application/json',
          Authorization: 'Bearer test-token',
        }),
      }),
    );
  });

  it('should handle trailing slash in API URL', async () => {
    setToken('test-token');
    const client = new ApiClient('http://localhost:8080/');

    fetchSpy.mockResolvedValue({
      ok: true,
      json: () => Promise.resolve({ data: 'test' }),
    } as Response);

    await client.get('/v1/articles');

    expect(fetchSpy).toHaveBeenCalledWith(
      'http://localhost:8080/v1/articles',
      expect.objectContaining({
        headers: expect.objectContaining({
          'Content-Type': 'application/json',
          Authorization: 'Bearer test-token',
        }),
      }),
    );
  });

  it('should not include Authorization header when token is not set', async () => {
    const client = new ApiClient('http://localhost:8080');

    fetchSpy.mockResolvedValue({
      ok: true,
      json: () => Promise.resolve({ data: 'test' }),
    } as Response);

    await client.get('/test');

    expect(fetchSpy).toHaveBeenCalledWith(
      'http://localhost:8080/test',
      expect.objectContaining({
        headers: expect.objectContaining({
          'Content-Type': 'application/json',
        }),
      }),
    );

    const callHeaders = fetchSpy.mock.calls[0][1].headers;
    expect(callHeaders.Authorization).toBeUndefined();
  });

  it('should throw error on 401 response', async () => {
    setToken('test-token');
    const client = new ApiClient('http://localhost:8080');

    fetchSpy.mockResolvedValue({
      ok: false,
      status: 401,
      statusText: 'Unauthorized',
    } as Response);

    await expect(client.get('/test')).rejects.toThrow('request failed: 401 Unauthorized');
  });

  it('should throw error on non-OK response', async () => {
    setToken('test-token');
    const client = new ApiClient('http://localhost:8080');

    fetchSpy.mockResolvedValue({
      ok: false,
      status: 500,
      statusText: 'Internal Server Error',
    } as Response);

    await expect(client.get('/test')).rejects.toThrow('request failed: 500 Internal Server Error');
  });
});
