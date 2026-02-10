async function getApiCredentials() {
  const data = await chrome.storage.local.get(['apiKey', 'apiUrl']);
  
  if (!data.apiKey || !data.apiUrl) {
    throw new Error('API credentials not configured');
  }

  return {
    apiKey: data.apiKey,
    apiUrl: data.apiUrl.replace(/\/$/, '')
  };
}

async function makeApiRequest(endpoint, options = {}) {
  const { apiKey, apiUrl } = await getApiCredentials();

  const headers = {
    'Authorization': 'Bearer ' + apiKey,
    ...options.headers
  };

  const response = await fetch(`${apiUrl}${endpoint}`, {
    ...options,
    headers
  });

  return response;
}
