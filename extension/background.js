chrome.runtime.onInstalled.addListener(() => {
  chrome.contextMenus.create({
    id: 'sendToKindle',
    title: 'Send to Kindle',
    contexts: ['link']
  });
});

chrome.contextMenus.onClicked.addListener(async (info, tab) => {
  if (info.menuItemId === 'sendToKindle' && info.linkUrl) {
    try {
      const data = await chrome.storage.local.get(['apiKey', 'apiUrl']);
      
      if (!data.apiKey || !data.apiUrl) {
        chrome.tabs.sendMessage(tab.id, { 
          type: 'error', 
          message: 'Please configure API settings in extension options' 
        });
        return;
      }

      const response = await fetch(`${data.apiUrl}/api/v1/articles`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'X-API-Key': data.apiKey
        },
        body: JSON.stringify({ url: info.linkUrl })
      });

      if (response.ok) {
        chrome.tabs.sendMessage(tab.id, { 
          type: 'success', 
          message: 'Link sent to Kindle!' 
        });
      } else {
        const error = await response.text();
        chrome.tabs.sendMessage(tab.id, { 
          type: 'error', 
          message: error || 'Failed to send link' 
        });
      }
    } catch (error) {
      chrome.tabs.sendMessage(tab.id, { 
        type: 'error', 
        message: 'Network error: ' + error.message 
      });
    }
  }
});

chrome.runtime.onMessage.addListener((request, sender, sendResponse) => {
  if (request.type === 'success' || request.type === 'error') {
    alert(request.message);
  }
});
