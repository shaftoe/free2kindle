<script lang="ts">
	import type { PageData } from './$types';
	let { data }: { data: PageData } = $props();
</script>

<h1>articles</h1>

{#if data.articles.length === 0}
	<p>no articles yet</p>
{:else}
	<ul>
		{#each data.articles as article (article.id)}
			<li>
				<article>
					<h2>{article.title || article.url}</h2>
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
