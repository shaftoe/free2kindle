<script lang="ts">
	import { resolve } from '$app/paths';
	import ArticleControls from './ArticleControls.svelte';
	import ArticleMeta from './ArticleMeta.svelte';
	import type { Article } from '$lib/server/types';
	let { article }: { article: Article } = $props();
</script>

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
	<p>added: {new Date(article.createdAt).toLocaleDateString()}</p>
	<ArticleMeta {article} />
	<ArticleControls {article} />
</article>
