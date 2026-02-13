document.addEventListener("DOMContentLoaded", () => {
  const form = document.getElementById("settingsForm");
  const statusEl = document.getElementById("status");
  const apiUrlInput = document.getElementById("apiUrl");
  const apiKeyInput = document.getElementById("apiKey");

  chrome.storage.local.get(["apiKey", "apiUrl"], (data) => {
    if (data.apiUrl) apiUrlInput.value = data.apiUrl;
    if (data.apiKey) apiKeyInput.value = data.apiKey;
  });

  form.addEventListener("submit", async (e) => {
    e.preventDefault();

    let apiUrl = apiUrlInput.value.trim();
    const apiKey = apiKeyInput.value.trim();

    if (!apiUrl || !apiKey) {
      showStatus("Please fill in all fields", true);
      return;
    }

    apiUrl = apiUrl.replace(/\/$/, "");

    try {
      const response = await fetch(`${apiUrl}/v1/health`, {
        method: "GET",
        headers: {
          "Content-Type": "application/json",
        },
      });

      if (response.ok) {
        await chrome.storage.local.set({ apiUrl, apiKey });
        showStatus("Settings saved successfully!");
      } else {
        const text = await response.text();
        showStatus("API error: " + text, true);
      }
    } catch (error) {
      showStatus("Failed to connect to API: " + error.message, true);
    }
  });

  function showStatus(message, isError = false) {
    statusEl.textContent = message;
    statusEl.className = `status ${isError ? "error" : "success"}`;
    setTimeout(() => {
      statusEl.className = "status";
    }, 5000);
  }
});
