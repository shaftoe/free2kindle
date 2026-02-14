<script lang="ts">
	import { enhance } from '$app/forms';
	import type { Article } from '$lib/server/types';
	let { article }: { article: Article } = $props();

	async function handleEnhance({ cancel }: { cancel: () => void }) {
		if (!window.confirm('Are you sure you want to delete this article?')) {
			cancel();
		}
	}
</script>

<form method="POST" action="?/delete" use:enhance={handleEnhance}>
	<input type="hidden" name="id" value={article.id} />
	<button type="submit">Delete</button>
</form>
