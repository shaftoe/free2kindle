import { fetchArticles } from "$lib/services/articles";
import { error } from "@sveltejs/kit";

export async function load({ url }) {
  const page = parseInt(url.searchParams.get("page") || "1", 10);
  const pageSize = parseInt(url.searchParams.get("page_size") || "20", 10);

  try {
    return await fetchArticles(page, pageSize);
  } catch (err) {
    const message =
      err instanceof Error ? err.message : "Failed to fetch articles";
    error(500, message);
  }
}
