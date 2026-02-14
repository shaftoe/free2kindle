<script lang="ts">
	import type { Article } from '$lib/server/types';
	let { article }: { article: Article } = $props();
</script>

<header>
	{#if article.title}
		<h1>{article.title}</h1>
	{:else}
		<h1>article</h1>
	{/if}

	<p>added: {new Date(article.createdAt).toLocaleDateString()}</p>

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
</header>
