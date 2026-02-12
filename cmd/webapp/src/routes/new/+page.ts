import { getToken } from "$lib/stores/token";
import { redirect } from "@sveltejs/kit";
import { browser } from "$app/environment";

export async function load({
  parent,
}: {
  parent: () => Promise<{ apiUrl: string }>;
}): Promise<{ apiUrl: string }> {
  const token = getToken();

  if (!token) {
    if (browser) {
      redirect(302, "/settings");
    }
  }

  const { apiUrl } = await parent();
  return { apiUrl };
}
