<svelte:head>
  <meta charset="UTF-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1.0" />
  <title>Articles</title>
</svelte:head>

<script lang="ts">
  import type { PageData } from './$types';

  interface Article {
    id: string;
    url: string;
    title?: string;
    createdAt: string;
    author?: string;
    siteName?: string;
    excerpt?: string;
    imageUrl?: string;
    wordCount?: number;
    readingTimeMinutes?: number;
    publishedAt?: string;
    deliveryStatus?: string;
  }

  let { data }: { data: PageData } = $props();

  const articles = $derived(data.articles as Article[]);
  const page = $derived(data.page);
  const pageSize = $derived(data.pageSize);
  const total = $derived(data.total);
  const hasMore = $derived(data.hasMore);
</script>

<h1>Articles</h1>

<p>Total: {total}</p>
<p>Page: {page} of {Math.ceil(total / pageSize)}</p>

{#if articles.length === 0}
  <p>No articles found.</p>
{:else}
  <ul>
    {#each articles as article}
      <li>
        <strong>{article.title || article.id}</strong>
        <br />
        <a href={article.url}>{article.url}</a>
        <br />
        Created: {article.createdAt}
        {#if article.author}
          <br />Author: {article.author}
        {/if}
        {#if article.siteName}
          <br />Site: {article.siteName}
        {/if}
        {#if article.wordCount}
          <br />Word count: {article.wordCount}
        {/if}
        {#if article.readingTimeMinutes}
          <br />Reading time: {article.readingTimeMinutes} min
        {/if}
      </li>
    {/each}
  </ul>
{/if}

{#if hasMore}
  <p><a href="/?page={page + 1}&page_size={pageSize}">Next page</a></p>
{/if}

{#if page > 1}
  <p><a href="/?page={page - 1}&page_size={pageSize}">Previous page</a></p>
{/if}
