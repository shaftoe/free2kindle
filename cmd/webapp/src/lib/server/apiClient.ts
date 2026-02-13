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

	private async request<T>(endpoint: string, options: RequestInit = {}): Promise<T> {
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

	async healthCheck(): Promise<HealthResponse> {
		return this.request<HealthResponse>('/v1/health');
	}

	async createArticle(data: CreateArticleRequest): Promise<CreateArticleResponse> {
		return this.request<CreateArticleResponse>('/v1/articles', {
			method: 'POST',
			body: JSON.stringify(data)
		});
	}

	async getArticles(page?: number, pageSize?: number): Promise<ListArticlesResponse> {
		const params = new URLSearchParams();
		if (page) params.set('page', page.toString());
		if (pageSize) params.set('page_size', pageSize.toString());
		const query = params.toString();
		return this.request<ListArticlesResponse>(`/v1/articles${query ? `?${query}` : ''}`);
	}

	async getArticle(id: string): Promise<Article> {
		return this.request<Article>(`/v1/articles/${id}`);
	}

	async deleteArticle(id: string): Promise<DeleteArticleResponse> {
		return this.request<DeleteArticleResponse>(`/v1/articles/${id}`, {
			method: 'DELETE'
		});
	}

	async deleteAllArticles(): Promise<DeleteArticleResponse> {
		return this.request<DeleteArticleResponse>('/v1/articles', {
			method: 'DELETE'
		});
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
