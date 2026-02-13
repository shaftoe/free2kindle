<script lang="ts">
	import { resolve } from '$app/paths';
	import BackNav from '$lib/components/BackNav.svelte';
	import type { PageData } from './$types';
	let { data }: { data: PageData } = $props();
</script>

<BackNav />
<h1>articles</h1>

{#if data.articles.length === 0}
	<p>no articles yet</p>
{:else}
	<ul>
		{#each data.articles as article (article.id)}
			<li>
				<article>
					{#if article.title}
						<h2><a href={resolve(`/articles/${article.id}`)}>{article.title}</a></h2>
					{:else}
						<h2><a href={resolve(`/articles/${article.id}`)}>{article.url}</a></h2>
					{/if}
					{#if article.imageUrl}
						<img src={article.imageUrl} alt={article.title} width="30%" />
					{/if}
					{#if article.excerpt}
						<p>{article.excerpt}</p>
					{/if}
					<p>
						Original link: <a href={article.url} target="_blank" rel="external">{article.url}</a>
					</p>
					{#if article.author}
						<p>by {article.author}</p>
					{/if}
					{#if article.siteName}
						<p>source: {article.siteName}</p>
					{/if}
					<p>added: {new Date(article.createdAt).toLocaleDateString()}</p>
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
				</article>
			</li>
		{/each}
	</ul>
{/if}

<p>total: {data.total}</p>
