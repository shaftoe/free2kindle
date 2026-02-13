import { apiClient } from "$lib/stores/apiClient";
import type { ArticlesResponse, Article, CreateArticleRequest, CreateArticleResponse } from "$lib/types";

export async function fetchArticles(
  page: number = 1,
  pageSize: number = 20,
): Promise<ArticlesResponse> {
  return apiClient.get<ArticlesResponse>(
    `/v1/articles?page=${page}&page_size=${pageSize}`,
  );
}

export async function fetchArticle(
  id: string,
): Promise<Article> {
  try {
    return apiClient.get<Article>(`/v1/articles/${id}`);
  } catch (err) {
    if (err instanceof Error && err.message.includes("404")) {
      throw new Error("Article not found");
    }
    throw err;
  }
}

export async function createArticle(
  url: string,
): Promise<CreateArticleResponse> {
  const req: CreateArticleRequest = { url };
  return apiClient.post<CreateArticleResponse>("/v1/articles", req);
}
