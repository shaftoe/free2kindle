(function(){
  const apiUrl = 'YOUR_FUNCTION_URL_HERE';
  const apiKey = 'YOUR_API_KEY_HERE';
  const currentUrl = window.location.href;
  
  fetch(apiUrl + '/api/v1/articles', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'X-API-Key': apiKey
    },
    body: JSON.stringify({url: currentUrl})
  })
  .then(r=>r.json())
  .then(d=>alert('✓ ' + d.title + '\nSent to Kindle successfully'))
  .catch(e=>alert('Error: ' + e.message));
})();
