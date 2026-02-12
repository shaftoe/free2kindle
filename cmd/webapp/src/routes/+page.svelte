<script lang="ts">
  import Article from "$lib/components/Article.svelte";
  import type { PageData } from "./$types";
  import type { Article as ArticleType } from "$lib/types";

  let { data }: { data: PageData } = $props();

  const articles = $derived(data.articles as ArticleType[]);
  const page = $derived(data.page);
  const pageSize = $derived(data.pageSize);
  const total = $derived(data.total);
  const hasMore = $derived(data.hasMore);
</script>

<svelte:head>
  <meta charset="UTF-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1.0" />
  <title>Articles - Save to Ink</title>
</svelte:head>

<main>
  <header class="page-header">
    <h1>My List</h1>

    <nav class="page-info" aria-label="Pagination information">
      <p class="page-stats">
        Total articles: <strong>{total}</strong>
      </p>
    </nav>
  </header>

  {#if articles.length === 0}
    <section class="empty-state">
      <p>No articles found.</p>
    </section>
  {:else}
    <section class="articles-list" aria-label="Articles list">
      {#each articles as article}
        <Article {...article} />
      {/each}
    </section>

    <nav class="pagination" aria-label="Pagination">
      <ul class="pagination-list">
        {#if page > 1}
          <li class="pagination-item">
            <a href="/?page={page - 1}&page_size={pageSize}" rel="prev">
              ← Previous page
            </a>
          </li>
        {/if}

        {#if hasMore}
          <li class="pagination-item">
            <a href="/?page={page + 1}&page_size={pageSize}" rel="next">
              Next page →
            </a>
          </li>
        {/if}
      </ul>
    </nav>
  {/if}
</main>
