import { describe, it, expect, vi, beforeEach } from 'vitest';
import { createApiClient, ApiError } from './apiClient';
import type { HealthResponse, CreateArticleResponse, ListArticlesResponse, Article } from './types';

const mockPublicApiUrl = vi.hoisted(() => ({
	value: 'http://localhost:8080' as string | undefined
}));

vi.mock('$env/dynamic/public', () => ({
	env: {
		get PUBLIC_API_URL() {
			return mockPublicApiUrl.value as unknown as string;
		}
	}
}));

export function mockEnvVar(value: string | undefined) {
	mockPublicApiUrl.value = value;
}

describe('ApiClient', () => {
	let mockFetch: ReturnType<typeof vi.fn>;

	beforeEach(() => {
		mockFetch = vi.fn();
		global.fetch = mockFetch as unknown as typeof global.fetch;
		mockEnvVar('http://localhost:8080');
	});

	describe('healthCheck', () => {
		it('should create health check request', async () => {
			mockFetch.mockResolvedValue({
				ok: true,
				json: async () => ({ status: 'ok' }) as HealthResponse
			});

			const client = createApiClient('test-key', 'http://localhost:8080');
			const result = await client.healthCheck(mockFetch as unknown as typeof globalThis.fetch);

			expect(mockFetch).toHaveBeenCalledWith('http://localhost:8080/v1/health', {
				headers: {
					'Content-Type': 'application/json',
					Authorization: 'Bearer test-key'
				}
			});
			expect(result).toEqual({ status: 'ok' });
		});

		it('should throw ApiError on failed request', async () => {
			mockFetch.mockResolvedValue({
				ok: false,
				status: 401,
				json: async () => ({ error: 'unauthorized' })
			});

			const client = createApiClient('invalid-key', 'http://localhost:8080');

			await expect(
				client.healthCheck(mockFetch as unknown as typeof globalThis.fetch)
			).rejects.toThrow(ApiError);
		});
	});

	describe('createArticle', () => {
		it('should create article with valid data', async () => {
			const mockResponse: CreateArticleResponse = {
				id: '123',
				title: 'Test Article',
				url: 'https://example.com',
				message: 'article created'
			};

			mockFetch.mockResolvedValue({
				ok: true,
				json: async () => mockResponse
			});

			const client = createApiClient('test-key', 'http://localhost:8080');
			const result = await client.createArticle(
				{ url: 'https://example.com' },
				mockFetch as unknown as typeof globalThis.fetch
			);

			expect(mockFetch).toHaveBeenCalledWith('http://localhost:8080/v1/articles', {
				method: 'POST',
				body: JSON.stringify({ url: 'https://example.com' }),
				headers: {
					'Content-Type': 'application/json',
					Authorization: 'Bearer test-key'
				}
			});
			expect(result).toEqual(mockResponse);
		});
	});

	describe('getArticles', () => {
		it('should fetch articles with pagination params', async () => {
			const mockResponse: ListArticlesResponse = {
				articles: [],
				page: 1,
				page_size: 20,
				total: 0,
				has_more: false
			};

			mockFetch.mockResolvedValue({
				ok: true,
				json: async () => mockResponse
			});

			const client = createApiClient('test-key', 'http://localhost:8080');
			const result = await client.getArticles(
				mockFetch as unknown as typeof globalThis.fetch,
				1,
				20
			);

			expect(mockFetch).toHaveBeenCalledWith(
				'http://localhost:8080/v1/articles?page=1&page_size=20',
				{
					headers: {
						'Content-Type': 'application/json',
						Authorization: 'Bearer test-key'
					}
				}
			);
			expect(result).toEqual(mockResponse);
		});

		it('should fetch articles without pagination params', async () => {
			const mockResponse: ListArticlesResponse = {
				articles: [],
				page: 1,
				page_size: 20,
				total: 0,
				has_more: false
			};

			mockFetch.mockResolvedValue({
				ok: true,
				json: async () => mockResponse
			});

			const client = createApiClient('test-key', 'http://localhost:8080');
			const result = await client.getArticles(mockFetch as unknown as typeof globalThis.fetch);

			expect(mockFetch).toHaveBeenCalledWith('http://localhost:8080/v1/articles', {
				headers: {
					'Content-Type': 'application/json',
					Authorization: 'Bearer test-key'
				}
			});
			expect(result).toEqual(mockResponse);
		});
	});

	describe('getArticle', () => {
		it('should fetch single article by id', async () => {
			const mockArticle: Article = {
				account: 'test-account',
				id: '123',
				url: 'https://example.com',
				createdAt: '2024-01-01T00:00:00Z',
				title: 'Test Article'
			};

			mockFetch.mockResolvedValue({
				ok: true,
				json: async () => mockArticle
			});

			const client = createApiClient('test-key', 'http://localhost:8080');
			const result = await client.getArticle(
				'123',
				mockFetch as unknown as typeof globalThis.fetch
			);

			expect(mockFetch).toHaveBeenCalledWith('http://localhost:8080/v1/articles/123', {
				headers: {
					'Content-Type': 'application/json',
					Authorization: 'Bearer test-key'
				}
			});
			expect(result).toEqual(mockArticle);
		});
	});

	describe('deleteArticle', () => {
		it('should delete article by id', async () => {
			mockFetch.mockResolvedValue({
				ok: true,
				json: async () => ({ deleted: 1 })
			});

			const client = createApiClient('test-key', 'http://localhost:8080');
			const result = await client.deleteArticle(
				'123',
				mockFetch as unknown as typeof globalThis.fetch
			);

			expect(mockFetch).toHaveBeenCalledWith('http://localhost:8080/v1/articles/123', {
				method: 'DELETE',
				headers: {
					'Content-Type': 'application/json',
					Authorization: 'Bearer test-key'
				}
			});
			expect(result).toEqual({ deleted: 1 });
		});
	});

	describe('deleteAllArticles', () => {
		it('should delete all articles', async () => {
			mockFetch.mockResolvedValue({
				ok: true,
				json: async () => ({ deleted: 5 })
			});

			const client = createApiClient('test-key', 'http://localhost:8080');
			const result = await client.deleteAllArticles(
				mockFetch as unknown as typeof globalThis.fetch
			);

			expect(mockFetch).toHaveBeenCalledWith('http://localhost:8080/v1/articles', {
				method: 'DELETE',
				headers: {
					'Content-Type': 'application/json',
					Authorization: 'Bearer test-key'
				}
			});
			expect(result).toEqual({ deleted: 5 });
		});
	});

	describe('createApiClient', () => {
		it('should throw ApiError for invalid baseUrl', () => {
			expect(() => createApiClient('test-key', 'not-a-url')).toThrow(ApiError);
		});

		it('should throw ApiError for baseUrl without protocol', () => {
			expect(() => createApiClient('test-key', 'example.com')).toThrow(ApiError);
		});

		it('should accept valid baseUrl', () => {
			const client = createApiClient('test-key', 'http://localhost:8080');
			expect(client).toBeDefined();
		});

		it('should accept valid https baseUrl', () => {
			const client = createApiClient('test-key', 'https://api.example.com');
			expect(client).toBeDefined();
		});

		it('should throw ApiError with status 400', () => {
			try {
				createApiClient('test-key', 'invalid-url');
			} catch (e) {
				expect(e).toBeInstanceOf(ApiError);
				expect((e as ApiError).status).toBe(500);
			}
		});

		it('should throw ApiError when PUBLIC_API_URL is not set', () => {
			mockEnvVar(undefined);
			expect(() => createApiClient('test-key')).toThrow(ApiError);
		});

		it('should throw ApiError with correct message when PUBLIC_API_URL is not set', () => {
			mockEnvVar(undefined);
			try {
				createApiClient('test-key');
			} catch (e) {
				expect(e).toBeInstanceOf(ApiError);
				expect((e as ApiError).status).toBe(500);
				expect((e as ApiError).message).toBe('PUBLIC_API_URL environment variable is not set');
			}
		});

		it('should throw ApiError when PUBLIC_API_URL is empty string', () => {
			mockEnvVar('');
			expect(() => createApiClient('test-key')).toThrow(ApiError);
		});

		it('should throw ApiError when PUBLIC_API_URL is invalid URL', () => {
			mockEnvVar('not-a-url');
			expect(() => createApiClient('test-key')).toThrow(ApiError);
		});

		it('should throw ApiError with correct message when PUBLIC_API_URL is invalid URL', () => {
			mockEnvVar('invalid-url');
			try {
				createApiClient('test-key');
			} catch (e) {
				expect(e).toBeInstanceOf(ApiError);
				expect((e as ApiError).status).toBe(500);
				expect((e as ApiError).message).toBe('invalid base url: invalid-url is not a valid url');
			}
		});

		it('should use PUBLIC_API_URL when baseUrl is not provided', () => {
			mockEnvVar('https://api.example.com');
			const client = createApiClient('test-key');
			expect(client).toBeDefined();
		});
	});
});
