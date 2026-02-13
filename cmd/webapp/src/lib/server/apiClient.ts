import { env } from '$env/dynamic/public';
import type {
	Article,
	CreateArticleRequest,
	CreateArticleResponse,
	ListArticlesResponse,
	DeleteArticleResponse,
	HealthResponse,
	ErrorResponse
} from './types';

export class ApiError extends Error {
	constructor(
		public status: number,
		message: string
	) {
		super(message);
		this.name = 'ApiError';
	}
}

export class ApiClient {
	constructor(
		private apiKey: string,
		private baseUrl: string
	) {}

	private async request<T>(
		endpoint: string,
		options: RequestInit = {},
		fetch: typeof globalThis.fetch
	): Promise<T> {
		const url = `${this.baseUrl}${endpoint}`;
		const response = await fetch(url, {
			...options,
			headers: {
				'Content-Type': 'application/json',
				Authorization: `Bearer ${this.apiKey}`,
				...options.headers
			}
		});

		if (!response.ok) {
			const error: ErrorResponse = await response.json().catch(() => ({ error: 'request failed' }));
			throw new ApiError(response.status, error.error);
		}

		return response.json();
	}

	async healthCheck(fetch: typeof globalThis.fetch): Promise<HealthResponse> {
		return this.request<HealthResponse>('/v1/health', {}, fetch);
	}

	async createArticle(
		data: CreateArticleRequest,
		fetch: typeof globalThis.fetch
	): Promise<CreateArticleResponse> {
		return this.request<CreateArticleResponse>(
			'/v1/articles',
			{
				method: 'POST',
				body: JSON.stringify(data)
			},
			fetch
		);
	}

	async getArticles(
		fetch: typeof globalThis.fetch,
		page?: number,
		pageSize?: number
	): Promise<ListArticlesResponse> {
		const params = new URLSearchParams();
		if (page) params.set('page', page.toString());
		if (pageSize) params.set('page_size', pageSize.toString());
		const query = params.toString();
		return this.request<ListArticlesResponse>(`/v1/articles${query ? `?${query}` : ''}`, {}, fetch);
	}

	async getArticle(id: string, fetch: typeof globalThis.fetch): Promise<Article> {
		return this.request<Article>(`/v1/articles/${id}`, {}, fetch);
	}

	async deleteArticle(id: string, fetch: typeof globalThis.fetch): Promise<DeleteArticleResponse> {
		return this.request<DeleteArticleResponse>(
			`/v1/articles/${id}`,
			{
				method: 'DELETE'
			},
			fetch
		);
	}

	async deleteAllArticles(fetch: typeof globalThis.fetch): Promise<DeleteArticleResponse> {
		return this.request<DeleteArticleResponse>(
			'/v1/articles',
			{
				method: 'DELETE'
			},
			fetch
		);
	}
}

export function createApiClient(apiKey: string, baseUrl?: string): ApiClient {
	const resolvedBaseUrl = baseUrl || env.PUBLIC_API_URL;

	if (!resolvedBaseUrl) {
		throw new ApiError(500, 'PUBLIC_API_URL environment variable is not set');
	}

	try {
		new URL(resolvedBaseUrl);
	} catch {
		throw new ApiError(500, `invalid base url: ${resolvedBaseUrl} is not a valid url`);
	}

	return new ApiClient(apiKey, resolvedBaseUrl);
}
