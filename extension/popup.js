document.addEventListener('DOMContentLoaded', () => {
  const sendBtn = document.getElementById('sendBtn');
  const statusEl = document.getElementById('status');
  const optionsLink = document.getElementById('optionsLink');

  function showStatus(message, isError = false) {
    statusEl.textContent = message;
    statusEl.className = `status ${isError ? 'error' : 'success'}`;
    sendBtn.disabled = true;
    setTimeout(() => {
      statusEl.className = 'status';
      sendBtn.disabled = false;
    }, 3000);
  }

  sendBtn.addEventListener('click', async () => {
    try {
      const [tab] = await chrome.tabs.query({ active: true, currentWindow: true });

      sendBtn.disabled = true;
      sendBtn.textContent = 'Sending...';

      const response = await makeApiRequest('/v1/articles', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({ url: tab.url })
      });

      if (response.ok) {
        showStatus('Article sent to Kindle!');
      } else {
        const error = await response.text();
        showStatus(error || 'Failed to send article', true);
      }
    } catch (error) {
      showStatus('Network error: ' + error.message, true);
    } finally {
      sendBtn.textContent = 'Send to Kindle';
    }
  });

  optionsLink.addEventListener('click', (e) => {
    e.preventDefault();
    chrome.runtime.openOptionsPage();
  });
});
