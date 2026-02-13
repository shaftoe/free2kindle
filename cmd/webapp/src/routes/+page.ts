import { fetchArticles } from "$lib/services/articles";
import { getToken } from "$lib/stores/token";
import { redirect } from "@sveltejs/kit";
import { browser } from "$app/environment";
import type { ArticlesResponse } from "$lib/types";

export async function load({
  parent,
  url,
  fetch,
}: {
  parent: () => Promise<{ apiUrl: string }>;
  url: URL;
  fetch: typeof globalThis.fetch;
}): Promise<ArticlesResponse> {
  const token = getToken();

  if (!token) {
    if (browser) {
      redirect(302, "/settings");
    }

    return {
      articles: [],
      page: 1,
      pageSize: 20,
      total: 0,
      hasMore: false,
    };
  }

  const { apiUrl } = await parent();
  const page = parseInt(url.searchParams.get("page") || "1", 10);
  const pageSize = parseInt(url.searchParams.get("page_size") || "20", 10);

  try {
    return await fetchArticles(apiUrl, page, pageSize, fetch);
  } catch (err) {
    if (err instanceof Error && err.message.includes("401")) {
      redirect(302, "/settings");
    }
    throw err;
  }
}
