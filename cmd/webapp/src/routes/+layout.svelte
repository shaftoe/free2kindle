<script lang="ts">
	import favicon from '$lib/assets/favicon.svg';
	import { setContext } from 'svelte';
	import { ApiClient } from '$lib/services/apiClient';
	import { API_CLIENT_KEY } from '$lib/types';
	import type { LayoutData } from './$types';

	let { children, data }: { children: any; data: LayoutData } = $props();

	$effect(() => {
		setContext(API_CLIENT_KEY, new ApiClient(data.apiUrl));
	});
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
</svelte:head>

{@render children()}
