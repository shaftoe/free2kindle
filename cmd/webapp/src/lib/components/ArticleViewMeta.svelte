<script lang="ts">
	import ArticleMeta from './ArticleMeta.svelte';
	import type { Article } from '$lib/server/types';
	let { article }: { article: Article } = $props();
</script>

<header>
	{#if article.title}
		<h1>{article.title}</h1>
	{:else}
		<h1>article</h1>
	{/if}

	<ArticleMeta {article} />
</header>
