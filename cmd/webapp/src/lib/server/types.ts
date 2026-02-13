// TODO: https://github.com/savetoink/savetoink/issues/2
export interface Article {
	account: string;
	id: string;
	url: string;
	createdAt: string;
	title?: string;
	content?: string;
	author?: string;
	siteName?: string;
	sourceDomain?: string;
	excerpt?: string;
	imageUrl?: string;
	contentType?: string;
	language?: string;
	error?: string;
	wordCount?: number;
	readingTimeMinutes?: number;
	publishedAt?: string;
	deliveryStatus?: 'pending' | 'delivered' | 'failed';
	deliveredFrom?: string;
	deliveredTo?: string;
	deliveredEmailUUID?: string;
	deliveredBy?: string;
}

export interface CreateArticleRequest {
	url: string;
}

export interface CreateArticleResponse {
	id: string;
	title: string;
	url: string;
	message: string;
	deliveryStatus?: string;
}

export interface ListArticlesResponse {
	articles: Article[];
	page: number;
	pageSize: number;
	total: number;
	hasMore: boolean;
}

export interface DeleteArticleResponse {
	deleted: number;
}

export interface HealthResponse {
	status: string;
}

export interface ErrorResponse {
	error: string;
}
