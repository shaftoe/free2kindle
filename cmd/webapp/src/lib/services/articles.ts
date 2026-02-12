import { ApiClient } from "./apiClient";
import type { ArticlesResponse, Article } from "$lib/types";

export async function fetchArticles(
  apiUrl: string,
  page: number = 1,
  pageSize: number = 20,
): Promise<ArticlesResponse> {
  const client = new ApiClient(apiUrl);
  return client.get<ArticlesResponse>(
    `/v1/articles?page=${page}&page_size=${pageSize}`,
  );
}

export async function fetchArticle(
  apiUrl: string,
  id: string,
): Promise<Article> {
  const client = new ApiClient(apiUrl);
  try {
    return client.get<Article>(`/v1/articles/${id}`);
  } catch (err) {
    if (err instanceof Error && err.message.includes("404")) {
      throw new Error("Article not found");
    }
    throw err;
  }
}
