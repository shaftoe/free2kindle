<script lang="ts">
  import token, { getToken, setToken } from "$lib/stores/token";
  import { goto } from "$app/navigation";

  let tokenInput = $state("");
  let saved = $state(false);
  let error = $state("");

  $effect(() => {
    const currentToken = getToken();
    tokenInput = currentToken || "";
  });

  function handleSave() {
    const trimmed = tokenInput.trim();
    if (!trimmed) {
      error = "token cannot be empty";
      return;
    }

    setToken(trimmed);
    saved = true;
    error = "";

    setTimeout(() => {
      saved = false;
    }, 2000);
  }

  function handleClear() {
    setToken(null);
    tokenInput = "";
    saved = false;
    error = "";
  }

  function handleBack() {
    goto("/");
  }
</script>

<svelte:head>
  <meta charset="UTF-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1.0" />
  <title>Settings - Save to Ink</title>
</svelte:head>

<main>
  <header class="page-header">
    <h1>Settings</h1>
  </header>

  <section class="settings-section">
    <h2>API Token</h2>

    <form onsubmit={(e) => e.preventDefault()}>
      <div class="form-group">
        <label for="token">API Token</label>
        <input
          id="token"
          type="password"
          bind:value={tokenInput}
          placeholder="Enter your API token"
          autocomplete="off"
        />
      </div>

      {#if error}
        <p class="error-message">{error}</p>
      {/if}

      {#if saved}
        <p class="success-message">Token saved successfully!</p>
      {/if}

      <div class="form-actions">
        <button type="submit" class="btn btn-primary" onclick={handleSave}>
          Save Token
        </button>
        <button type="button" class="btn btn-secondary" onclick={handleClear}>
          Clear Token
        </button>
        <button type="button" class="btn btn-secondary" onclick={handleBack}>
          Back
        </button>
      </div>
    </form>

    <div class="info-box">
      <h3>About API Tokens</h3>
      <p>
        Your API token is used to authenticate with the Save to Ink backend. It
        is stored locally in your browser's localStorage and is never sent to
        any third-party servers.
      </p>
      <p>
        Without a valid API token, you will not be able to fetch articles or
        interact with the service.
      </p>
    </div>
  </section>
</main>
