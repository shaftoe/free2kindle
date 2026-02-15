<script lang="ts">
	import ArticleMeta from './ArticleMeta.svelte';
	import type { Article } from '$lib/server/types';
	let { article }: { article: Article } = $props();
</script>

<ArticleMeta {article} mode="card" />
