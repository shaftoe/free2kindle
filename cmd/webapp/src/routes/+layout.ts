import { env } from "$env/dynamic/public";

export function load() {
  return {
    apiUrl: env.PUBLIC_API_URL || "http://localhost:8080",
  };
}
