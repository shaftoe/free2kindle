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
    contentType,
    language,
    sourceDomain,
    deliveredFrom,
    deliveredTo,
    deliveredBy,
    deliveredEmailUUID,
  }: Article = $props();
</script>

<article class="article-card" data-article-id={id}>
  <h2 class="article-title">
    {#if title}
      {title}
    {:else}
      Article {id}
    {/if}
  </h2>

  {#if excerpt}
    <p class="article-excerpt">{excerpt}</p>
  {/if}

  <table class="article-table">
    <tbody>
      <tr>
        <th scope="row">ID</th>
        <td>{id}</td>
      </tr>
      <tr>
        <th scope="row">URL</th>
        <td
          ><a href={url} rel="noopener noreferrer" target="_blank">{url}</a></td
        >
      </tr>
      {#if publishedAt}
        <tr>
          <th scope="row">Published</th>
          <td
            ><time datetime={publishedAt}
              >{new Date(publishedAt).toLocaleDateString()}</time
            ></td
          >
        </tr>
      {/if}
      <tr>
        <th scope="row">Added</th>
        <td
          ><time datetime={createdAt}
            >{new Date(createdAt).toLocaleDateString()}</time
          ></td
        >
      </tr>
      {#if author}
        <tr>
          <th scope="row">Author</th>
          <td>{author}</td>
        </tr>
      {/if}
      {#if siteName}
        <tr>
          <th scope="row">Source</th>
          <td>{siteName}</td>
        </tr>
      {/if}
      {#if sourceDomain}
        <tr>
          <th scope="row">Source domain</th>
          <td>{sourceDomain}</td>
        </tr>
      {/if}
      {#if language}
        <tr>
          <th scope="row">Language</th>
          <td>{language}</td>
        </tr>
      {/if}
      {#if contentType}
        <tr>
          <th scope="row">Type</th>
          <td>{contentType}</td>
        </tr>
      {/if}
      {#if wordCount}
        <tr>
          <th scope="row">Word count</th>
          <td>{wordCount.toLocaleString()}</td>
        </tr>
      {/if}
      {#if readingTimeMinutes}
        <tr>
          <th scope="row">Reading time</th>
          <td>{readingTimeMinutes} min</td>
        </tr>
      {/if}
      {#if deliveryStatus}
        <tr>
          <th scope="row">Delivery status</th>
          <td>{deliveryStatus}</td>
        </tr>
      {/if}
    </tbody>
  </table>

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

  <p class="article-link">
    <a href={url} rel="noopener noreferrer" target="_blank">Read article →</a>
  </p>
</article>
