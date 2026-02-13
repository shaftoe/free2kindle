import { fetchArticle } from "$lib/services/articles";
import { getToken } from "$lib/stores/token";
import { redirect } from "@sveltejs/kit";
import { browser } from "$app/environment";
import type { Article } from "$lib/types";

export async function load({
  parent,
  params,
  fetch,
}: {
  parent: () => Promise<{ apiUrl: string }>;
  params: { id: string };
  fetch: typeof globalThis.fetch;
}): Promise<{ article: Article }> {
  const token = getToken();
  
  if (!token) {
    if (browser) {
      redirect(302, "/settings");
    }
    redirect(302, "/");
  }

  const { id } = params;
  if (!id) {
    redirect(302, "/");
  }

  const { apiUrl } = await parent();

  try {
    const article = await fetchArticle(apiUrl, id, fetch);
    return { article };
  } catch (err) {
    if (err instanceof Error && err.message.includes("401")) {
      redirect(302, "/settings");
    }
    throw err;
  }
}
