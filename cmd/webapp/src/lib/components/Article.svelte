<script lang="ts">
  import type { Article } from "$lib/types";

  let {
    id,
    url,
    title,
    createdAt,
    author,
    siteName,
    excerpt,
    imageUrl,
    wordCount,
    readingTimeMinutes,
    publishedAt,
    deliveryStatus,
    language,
    sourceDomain,
    deliveredFrom,
    deliveredTo,
    deliveredBy,
    deliveredEmailUUID,
  }: Article = $props();

  const displaySource = $derived(siteName || sourceDomain);
</script>

<article class="article-card" data-article-id={id}>
  <div class="article-content">
    <h2 class="article-title">
      <a href={`/article/${id}`}>{title}</a>
    </h2>

    <div class="article-image">
      {#if imageUrl}
        <img src={imageUrl} alt={title} loading="lazy" width="30%" />
      {/if}
    </div>

    {#if excerpt}
      <p class="article-excerpt">{excerpt}</p>
    {/if}

    <div class="article-meta">
      {#if displaySource}
        <span class="article-source">{displaySource}</span>
      {/if}

      {#if readingTimeMinutes}
        <span class="article-time">{readingTimeMinutes} min read</span>
      {/if}

      {#if createdAt}
        <span class="article-date">
          Added <time datetime={createdAt}>
            {new Date(createdAt).toLocaleDateString()}
          </time>
        </span>
      {/if}
    </div>

    <div class="article-tags">
      {#if author}
        <span class="tag">{author}</span>
      {/if}
      {#if language}
        <span class="tag">{language}</span>
      {/if}
    </div>
  </div>

  {#if wordCount || publishedAt || deliveryStatus}
    <div class="article-details">
      {#if publishedAt}
        <span class="detail-item">
          <span class="detail-label">Published</span>
          <time datetime={publishedAt}
            >{new Date(publishedAt).toLocaleDateString()}</time
          >
        </span>
      {/if}

      {#if wordCount}
        <span class="detail-item">
          <span class="detail-label">Words</span>
          {wordCount.toLocaleString()}
        </span>
      {/if}

      {#if deliveryStatus}
        <span class="detail-item">
          <span class="detail-label">Status</span>
          <span class="status-{deliveryStatus}">{deliveryStatus}</span>
        </span>
      {/if}
    </div>
  {/if}

  {#if deliveredFrom || deliveredTo || deliveredBy || deliveredEmailUUID}
    <details class="article-delivery">
      <summary>Delivery Information</summary>
      <table>
        <tbody>
          {#if deliveredFrom}
            <tr>
              <th scope="row">From</th>
              <td>{deliveredFrom}</td>
            </tr>
          {/if}
          {#if deliveredTo}
            <tr>
              <th scope="row">To</th>
              <td>{deliveredTo}</td>
            </tr>
          {/if}
          {#if deliveredBy}
            <tr>
              <th scope="row">Provider</th>
              <td>{deliveredBy}</td>
            </tr>
          {/if}
          {#if deliveredEmailUUID}
            <tr>
              <th scope="row">Email ID</th>
              <td><code>{deliveredEmailUUID}</code></td>
            </tr>
          {/if}
        </tbody>
      </table>
    </details>
  {/if}
</article>
