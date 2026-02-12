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
}

export interface ArticlesResponse {
  articles: Article[];
  page: number;
  pageSize: number;
  total: number;
  hasMore: boolean;
}
