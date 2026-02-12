import { fetchArticle } from "$lib/services/articles";
import { error } from "@sveltejs/kit";

export async function load({ params }) {
  const { id } = params;

  if (!id) {
    error(400, "Article ID is required");
  }

  try {
    const article = await fetchArticle(id);
    return { article };
  } catch (err) {
    const message = err instanceof Error ? err.message : "Failed to fetch article";
    error(500, message);
  }
}
