<script lang="ts">
  import type { PageData } from "./$types";
  import type { Article } from "$lib/types";

  let { data }: { data: PageData } = $props();

  const article = $derived(data.article as Article);
</script>

<svelte:head>
  <meta charset="UTF-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1.0" />
  <title>{article.title || 'Article'} - Free2Kindle</title>
</svelte:head>

<main>
  <nav class="back-nav">
    <a href="/">← Back to List</a>
  </nav>

  {#if article}
    <article>
      <header>
        <h1>{article.title || 'Untitled'}</h1>

        <div class="meta">
          {#if article.siteName || article.sourceDomain}
            <span>{article.siteName || article.sourceDomain}</span>
          {/if}

          {#if article.author}
            <span>by {article.author}</span>
          {/if}

          {#if article.readingTimeMinutes}
            <span>{article.readingTimeMinutes} min read</span>
          {/if}

          <time datetime={article.createdAt}>
            {new Date(article.createdAt).toLocaleDateString()}
          </time>
        </div>

        {#if article.url}
          <p class="original-url">
            <a href={article.url} target="_blank" rel="noopener noreferrer">
              Original article
            </a>
          </p>
        {/if}
      </header>

      <div class="content">
        {#if article.content}
          {@html article.content}
        {:else}
          <p>No content available.</p>
        {/if}
      </div>
    </article>
  {/if}
</main>
