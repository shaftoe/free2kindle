<script lang="ts">
  import { createArticle } from "$lib/services/articles";
  import { goto } from "$app/navigation";
  import type { PageData } from "./$types";

  let { data }: { data: PageData } = $props();

  let urlInput = $state("");
  let loading = $state(false);
  let error = $state("");
  let success = $state(false);

  async function handleSubmit(e: Event) {
    e.preventDefault();
    const trimmedUrl = urlInput.trim();

    if (!trimmedUrl) {
      error = "url cannot be empty";
      return;
    }

    loading = true;
    error = "";
    success = false;

    try {
      await createArticle(data.apiUrl, trimmedUrl);
      success = true;
      setTimeout(() => {
        goto("/");
      }, 1000);
    } catch (err) {
      if (err instanceof Error) {
        error = err.message;
      } else {
        error = "failed to add article";
      }
    } finally {
      loading = false;
    }
  }
</script>

<svelte:head>
  <meta charset="UTF-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1.0" />
  <title>Add Article - Save to Ink</title>
</svelte:head>

<main>
  <header class="page-header">
    <h1>Add Article</h1>
  </header>

  <section>
    <form onsubmit={handleSubmit} style="max-width: 600px; margin: 0 auto;">
      <div style="margin-bottom: 20px;">
        <label for="url" style="display: block; margin-bottom: 8px; font-weight: bold;">Article URL</label>
        <input
          id="url"
          type="url"
          bind:value={urlInput}
          placeholder="https://example.com/article"
          disabled={loading}
          style="width: 100%; padding: 12px; font-size: 16px; border: 1px solid #ddd; border-radius: 4px;"
        />
      </div>

      {#if error}
        <p style="color: #d32f2f; margin-bottom: 20px; padding: 10px; background: #ffebee; border-radius: 4px; border: 1px solid #ffcdd2;">
          {error}
        </p>
      {/if}

      {#if success}
        <p style="color: #2e7d32; margin-bottom: 20px; padding: 10px; background: #e8f5e9; border-radius: 4px; border: 1px solid #c8e6c9;">
          article added successfully! redirecting...
        </p>
      {/if}

      <div style="display: flex; gap: 10px;">
        <button
          type="submit"
          disabled={loading}
          style="padding: 12px 24px; font-size: 16px; background: #1976d2; color: white; border: none; border-radius: 4px; cursor: pointer; {loading ? 'opacity: 0.6; cursor: not-allowed;' : ''}"
        >
          {loading ? "adding..." : "add article"}
        </button>
        <button
          type="button"
          onclick={() => goto("/")}
          disabled={loading}
          style="padding: 12px 24px; font-size: 16px; background: #6c757d; color: white; border: none; border-radius: 4px; cursor: pointer;"
        >
          cancel
        </button>
      </div>
    </form>
  </section>
</main>
