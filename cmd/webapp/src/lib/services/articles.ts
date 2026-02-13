import { ApiClient, getApiClient } from "./apiClient";
import type { ArticlesResponse, Article, CreateArticleRequest, CreateArticleResponse } from "$lib/types";

export async function fetchArticles(
  apiUrl: string,
  page: number = 1,
  pageSize: number = 20,
  fetchFn: typeof fetch = fetch,
): Promise<ArticlesResponse> {
  const client = getApiClient() ?? new ApiClient(apiUrl);
  return client.get<ArticlesResponse>(
    `/v1/articles?page=${page}&page_size=${pageSize}`,
    fetchFn,
  );
}

export async function fetchArticle(
  apiUrl: string,
  id: string,
  fetchFn: typeof fetch = fetch,
): Promise<Article> {
  const client = getApiClient() ?? new ApiClient(apiUrl);
  try {
    return client.get<Article>(`/v1/articles/${id}`, fetchFn);
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
  fetchFn: typeof fetch = fetch,
): Promise<CreateArticleResponse> {
  const client = getApiClient() ?? new ApiClient(apiUrl);
  const req: CreateArticleRequest = { url };
  return client.post<CreateArticleResponse>("/v1/articles", req, fetchFn);
}
