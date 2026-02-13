<script lang="ts">
	import type { PageProps } from './$types';
	let { form }: PageProps = $props();
</script>

<h1>Settings</h1>
<p>Enter your API key to access the article management system.</p>

{#if form?.error}
	<p class="error">{form.error}</p>
{/if}

<form method="POST">
	<label>
		API Key
		<input type="password" name="apiKey" required autocomplete="current-password" />
	</label>
	<button type="submit">Save</button>
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
