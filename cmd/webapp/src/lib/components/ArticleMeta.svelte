<script lang="ts">
	import { formatDate } from '$lib/utils/date';
	import type { Article } from '$lib/server/types';
	let { article }: { article: Article } = $props();
</script>

{#if article.author}
	<p>by {article.author}</p>
{/if}

{#if article.siteName}
	<p>source: {article.siteName}</p>
{/if}

{#if article.wordCount}
	<p>{article.wordCount} words</p>
{/if}

{#if article.readingTimeMinutes}
	<p>{article.readingTimeMinutes} min read</p>
{/if}

{#if article.deliveryStatus}
	<p>status: {article.deliveryStatus}</p>
{/if}

{#if article.error}
	<p class="error">error: {article.error}</p>
{/if}

{#if article.createdAt}
	<p>added: {formatDate(article.createdAt)}</p>
{/if}
