<script lang="ts">
	import BackNav from '$lib/components/BackNav.svelte';
	import type { PageProps } from './$types';
	let { form }: PageProps = $props();
</script>

<BackNav href="/" />
<h1>Add Article</h1>

{#if form?.error}
	<p class="error">{form.error}</p>
{/if}

<form method="POST">
	<label>
		URL
		<input type="url" name="url" required placeholder="https://example.com/article" />
	</label>
	<button type="submit">Add</button>
</form>

<style>
	form {
		display: flex;
		flex-direction: column;
		gap: 1rem;
		max-width: 400px;
	}
	label {
		display: flex;
		flex-direction: column;
		gap: 0.25rem;
	}
	.error {
		color: red;
	}
</style>
