<script lang="ts">
	import { enhance } from '$app/forms';

	async function handleEnhance({ cancel }: { cancel: () => void }) {
		if (!window.confirm('Are you sure you want to delete this article?')) {
			cancel();
		}
	}
</script>

<form method="POST" action="?/delete" use:enhance={handleEnhance}>
	<button type="submit">Delete</button>
</form>
