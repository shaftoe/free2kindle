<script lang="ts">
	import ArticleControls from '$lib/components/ArticleControls.svelte';
	import ArticleViewMeta from '$lib/components/ArticleViewMeta.svelte';
	import type { PageData } from './$types';
	let { data }: { data: PageData } = $props();
</script>

<ArticleViewMeta article={data.article} />

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

<ArticleControls article={data.article} />
