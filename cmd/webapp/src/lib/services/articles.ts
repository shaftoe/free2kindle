import { env } from "$env/dynamic/private";
import type { ArticlesResponse, Article } from "$lib/types";

export async function fetchArticles(
  page: number = 1,
  pageSize: number = 20,
): Promise<ArticlesResponse> {
  const apiUrl = env.API_URL || "http://localhost:8080";
  const apiKey = env.API_KEY;

  const response = await fetch(
    `${apiUrl}/v1/articles?page=${page}&page_size=${pageSize}`,
    {
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${apiKey}`,
      },
    },
  );

  if (!response.ok) {
    throw new Error(
      `Failed to fetch articles: ${response.status} ${response.statusText}`,
    );
  }

  return (await response.json()) as ArticlesResponse;
}

export async function fetchArticle(id: string): Promise<Article> {
  const apiUrl = env.API_URL || "http://localhost:8080";
  const apiKey = env.API_KEY;

  const response = await fetch(`${apiUrl}/v1/articles/${id}`, {
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${apiKey}`,
    },
  });

  if (!response.ok) {
    if (response.status === 404) {
      throw new Error("Article not found");
    }
    throw new Error(
      `Failed to fetch article: ${response.status} ${response.statusText}`,
    );
  }

  return (await response.json()) as Article;
}
