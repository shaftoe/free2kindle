<script lang="ts">
	import { resolve } from '$app/paths';

	type Route = '/' | '/settings';

	let { href = '/' as Route }: { href?: Route } = $props();
</script>

<nav>
	<a href={resolve(href)}>← Back</a>
</nav>
