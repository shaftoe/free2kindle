<script lang="ts">
	import { goto } from '$app/navigation';

	let { page, has_more }: { page: number; has_more: boolean } = $props();

	function navigateTo(newPage: number) {
		// eslint-disable-next-line svelte/no-navigation-without-resolve
		goto(`/?page=${newPage}`, { replaceState: true, noScroll: true });
	}
</script>

<nav>
	{#if page > 1}
		<button onclick={() => navigateTo(page - 1)}>Previous</button>
	{/if}
	{#if has_more}
		<button onclick={() => navigateTo(page + 1)}>Next</button>
	{/if}
</nav>
