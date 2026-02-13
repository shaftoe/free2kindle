import { ApiClient, getApiClient } from "./apiClient";
import type { ArticlesResponse, Article, CreateArticleRequest, CreateArticleResponse } from "$lib/types";

export async function fetchArticles(
  apiUrl: string,
  page: number = 1,
  pageSize: number = 20,
): Promise<ArticlesResponse> {
  const client = getApiClient() ?? new ApiClient(apiUrl);
  return client.get<ArticlesResponse>(
    `/v1/articles?page=${page}&page_size=${pageSize}`,
  );
}

export async function fetchArticle(
  apiUrl: string,
  id: string,
): Promise<Article> {
  const client = getApiClient() ?? new ApiClient(apiUrl);
  try {
    return client.get<Article>(`/v1/articles/${id}`);
  } catch (err) {
    if (err instanceof Error && err.message.includes("404")) {
      throw new Error("Article not found");
    }
    throw err;
  }
}

export async function createArticle(
  apiUrl: string,
  url: string,
): Promise<CreateArticleResponse> {
  const client = getApiClient() ?? new ApiClient(apiUrl);
  const req: CreateArticleRequest = { url };
  return client.post<CreateArticleResponse>("/v1/articles", req);
}
