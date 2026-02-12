import { env } from '$env/dynamic/private';
import { error } from '@sveltejs/kit';

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

interface ArticlesResponse {
  articles: Article[];
  page: number;
  pageSize: number;
  total: number;
  hasMore: boolean;
}

export async function load({ url }) {
  const apiUrl = env.API_URL || 'http://localhost:8080';
  const apiKey = env.API_KEY;

  try {
    const page = url.searchParams.get('page') || '1';
    const pageSize = url.searchParams.get('page_size') || '20';

    const response = await fetch(
      `${apiUrl}/v1/articles?page=${page}&page_size=${pageSize}`,
      {
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${apiKey}`,
        },
      }
    );

    if (!response.ok) {
      error(response.status, `Error fetching articles: ${response.status} ${response.statusText}`);
    }

    const data = (await response.json()) as ArticlesResponse;

    return {
      articles: data.articles,
      page: data.page,
      pageSize: data.pageSize,
      total: data.total,
      hasMore: data.hasMore,
    };
  } catch (err) {
    console.error('Error fetching articles:', err);
    error(500, `Error fetching articles: ${err instanceof Error ? err.message : 'Unknown error'}`);
  }
}
