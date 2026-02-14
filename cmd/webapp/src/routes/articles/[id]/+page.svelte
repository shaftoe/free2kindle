<script lang="ts">
	import BackNav from '$lib/components/BackNav.svelte';
	import type { PageData } from './$types';
	let { data }: { data: PageData } = $props();
</script>

<BackNav href="/" />
<header>
	{#if data.article.title}
		<h1>{data.article.title}</h1>
	{:else}
		<h1>article</h1>
	{/if}

	<p>added: {new Date(data.article.createdAt).toLocaleDateString()}</p>

	{#if data.article.author}
		<p>by {data.article.author}</p>
	{/if}

	{#if data.article.siteName}
		<p>source: {data.article.siteName}</p>
	{/if}

	{#if data.article.wordCount}
		<p>{data.article.wordCount} words</p>
	{/if}

	{#if data.article.readingTimeMinutes}
		<p>{data.article.readingTimeMinutes} min read</p>
	{/if}

	{#if data.article.deliveryStatus}
		<p>status: {data.article.deliveryStatus}</p>
	{/if}

	{#if data.article.error}
		<p class="error">error: {data.article.error}</p>
	{/if}
</header>

<hr />
{#if data.article.content}
	<div>
		<!-- eslint-disable-next-line svelte/no-at-html-tags -->
		{@html data.article.content}
	</div>
{:else}
	<p>content not yet available</p>
{/if}

<p>
	<!-- eslint-disable-next-line svelte/no-navigation-without-resolve -->
	<a href={data.article.url} target="_blank" rel="noopener noreferrer">original article</a>
</p>
