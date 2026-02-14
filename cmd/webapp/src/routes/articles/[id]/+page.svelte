<script lang="ts">
	import ArticleView from '$lib/components/ArticleView.svelte';
	import type { PageData } from './$types';
	let { data }: { data: PageData } = $props();
</script>

<ArticleView article={data.article} />
