export interface Article {
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
  contentType?: string;
  language?: string;
  sourceDomain?: string;
  deliveredFrom?: string;
  deliveredTo?: string;
  deliveredBy?: string;
  deliveredEmailUUID?: string;
  content?: string;
}

export interface ArticlesResponse {
  articles: Article[];
  page: number;
  pageSize: number;
  total: number;
  hasMore: boolean;
}

export interface CreateArticleRequest {
  url: string;
}

export interface CreateArticleResponse {
  id: string;
  title: string;
  url: string;
  message: string;
  delivery_status?: string;
}

import { ApiClient } from "$lib/services/apiClient";

export const API_CLIENT_KEY = Symbol("apiClient");
