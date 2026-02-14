<script lang="ts">
	import Article from '$lib/components/Article.svelte';
	import Navigator from '$lib/components/Navigator.svelte';
	import type { PageData } from './$types';
	let { data }: { data: PageData } = $props();
</script>

<h1>My List (<span>{data.total} articles)</span></h1>

{#if data.error}
	<p class="error">failed to load articles: {data.error}</p>
{:else if data.articles.length > 0}
	<ul>
		{#each data.articles as article (article.id)}
			<li>
				<Article {article} />
			</li>
		{/each}
	</ul>
{/if}

<p>total: {data.total}</p>

{#if !data.error && data.articles.length > 0}
	<Navigator page={data.page} has_more={data.has_more} />
{/if}
