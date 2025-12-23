package main

import (
	"fmt"
	"net/http"
)

// ServeUI writes a small single-file web UI to the response.
func ServeUI(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, uiHTML)
}

var uiHTML = `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>URL Shortener - Fast & Simple</title>
  <style>
    * { margin: 0; padding: 0; box-sizing: border-box; }

    :root {
      --primary: rgb(52, 73, 94);
      --primary-light: rgb(209, 226, 240);
      --primary-lighter: rgb(240, 245, 250);
      --success: #5dade2;
      --error: #e74c3c;
      --bg: rgb(240, 245, 250);
      --card-bg: #ffffff;
      --text: rgb(52, 73, 94);
      --text-muted: #7f8c8d;
      --border: rgb(209, 226, 240);
      --shadow: 0 2px 8px rgba(52, 73, 94, 0.08);
      --shadow-lg: 0 8px 24px rgba(52, 73, 94, 0.12);
    }

    body {
      font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
      background: linear-gradient(180deg, rgb(240, 245, 250) 0%, rgb(209, 226, 240) 100%);
      min-height: 100vh;
      padding: 20px;
      color: var(--text);
    }

    .container {
      max-width: 800px;
      margin: 0 auto;
      animation: fadeIn 0.5s ease-in;
    }

    @keyframes fadeIn {
      from { opacity: 0; transform: translateY(20px); }
      to { opacity: 1; transform: translateY(0); }
    }

    .header {
      text-align: center;
      color: var(--primary);
      margin-bottom: 40px;
    }

    .header h1 {
      font-size: 2.5rem;
      font-weight: 700;
      margin-bottom: 8px;
      letter-spacing: -0.5px;
    }

    .header p {
      font-size: 1.1rem;
      color: var(--text-muted);
    }

    .card {
      background: var(--card-bg);
      border-radius: 16px;
      padding: 32px;
      box-shadow: var(--shadow-lg);
      margin-bottom: 24px;
    }

    .input-group {
      margin-bottom: 20px;
    }

    .input-group label {
      display: block;
      margin-bottom: 8px;
      font-weight: 500;
      color: var(--text);
      font-size: 0.9rem;
    }

    input[type="url"],
    input[type="text"] {
      width: 100%;
      padding: 14px 16px;
      font-size: 16px;
      border: 2px solid var(--border);
      border-radius: 12px;
      transition: all 0.2s;
      font-family: inherit;
    }

    input:focus {
      outline: none;
      border-color: var(--primary);
      box-shadow: 0 0 0 4px rgba(99, 102, 241, 0.1);
    }

    .button-group {
      display: flex;
      gap: 12px;
      margin-top: 24px;
    }

    button {
      flex: 1;
      padding: 14px 24px;
      font-size: 16px;
      font-weight: 600;
      border: none;
      border-radius: 12px;
      cursor: pointer;
      transition: all 0.2s;
      font-family: inherit;
    }

    .btn-primary {
      background: var(--primary);
      color: white;
    }

    .btn-primary:hover {
      background: #2c3e50;
      transform: translateY(-1px);
      box-shadow: 0 4px 12px rgba(52, 73, 94, 0.3);
    }

    .btn-secondary {
      background: white;
      color: var(--text);
      border: 2px solid var(--border);
    }

    .btn-secondary:hover {
      background: var(--primary-lighter);
      border-color: var(--primary);
    }

    .result-card {
      display: none;
      background: var(--primary-lighter);
      border: 2px solid var(--primary-light);
      border-radius: 12px;
      padding: 24px;
      margin-top: 24px;
      animation: slideIn 0.3s ease-out;
    }

    @keyframes slideIn {
      from { opacity: 0; transform: translateY(-10px); }
      to { opacity: 1; transform: translateY(0); }
    }

    .result-card.show {
      display: block;
    }

    .result-header {
      display: flex;
      align-items: center;
      justify-content: space-between;
      margin-bottom: 16px;
    }

    .result-header h3 {
      color: var(--primary);
      font-size: 1.1rem;
    }

    .short-url {
      display: flex;
      align-items: center;
      gap: 12px;
      background: white;
      padding: 16px;
      border-radius: 8px;
      margin-bottom: 12px;
    }

    .short-url a {
      flex: 1;
      color: var(--primary);
      text-decoration: none;
      font-weight: 600;
      font-size: 1.1rem;
      word-break: break-all;
    }

    .copy-btn {
      padding: 8px 16px;
      background: var(--primary);
      color: white;
      border: none;
      border-radius: 8px;
      cursor: pointer;
      font-size: 14px;
      font-weight: 500;
      transition: all 0.2s;
    }

    .copy-btn:hover {
      background: var(--primary-dark);
    }

    .copy-btn.copied {
      background: var(--success);
    }

    .original-url {
      color: var(--text-muted);
      font-size: 0.9rem;
      padding: 12px 16px;
      background: white;
      border-radius: 8px;
      word-break: break-all;
    }

    .error-message {
      display: none;
      background: #fef2f2;
      border: 2px solid var(--error);
      color: #991b1b;
      padding: 16px;
      border-radius: 12px;
      margin-top: 16px;
      animation: shake 0.4s;
    }

    @keyframes shake {
      0%, 100% { transform: translateX(0); }
      25% { transform: translateX(-10px); }
      75% { transform: translateX(10px); }
    }

    .error-message.show {
      display: block;
    }

    .stats {
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
      gap: 16px;
      margin-bottom: 24px;
    }

    .stat-card {
      background: white;
      padding: 20px;
      border-radius: 12px;
      text-align: center;
      border: 2px solid var(--border);
    }

    .stat-value {
      font-size: 2rem;
      font-weight: 700;
      color: var(--primary);
      margin-bottom: 4px;
    }

    .stat-label {
      color: var(--text-muted);
      font-size: 0.9rem;
    }

    .url-list {
      display: none;
      margin-top: 24px;
    }

    .url-list.show {
      display: block;
    }

    .url-item {
      background: white;
      padding: 20px;
      border-radius: 12px;
      margin-bottom: 12px;
      border: 2px solid var(--border);
      transition: all 0.2s;
    }

    .url-item:hover {
      border-color: var(--primary);
      transform: translateX(4px);
    }

    .url-item-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 12px;
    }

    .url-item-code {
      font-weight: 600;
      color: var(--primary);
      font-size: 1.1rem;
    }

    .url-item-clicks {
      background: var(--bg);
      padding: 4px 12px;
      border-radius: 20px;
      font-size: 0.85rem;
      color: var(--text-muted);
    }

    .url-item-original {
      color: var(--text-muted);
      font-size: 0.9rem;
      word-break: break-all;
      margin-bottom: 8px;
    }

    .url-item-date {
      color: var(--text-muted);
      font-size: 0.8rem;
    }

    .loading {
      display: none;
      text-align: center;
      padding: 20px;
    }

    .spinner {
      border: 3px solid var(--border);
      border-top: 3px solid var(--primary);
      border-radius: 50%;
      width: 40px;
      height: 40px;
      animation: spin 1s linear infinite;
      margin: 0 auto;
    }

    @keyframes spin {
      0% { transform: rotate(0deg); }
      100% { transform: rotate(360deg); }
    }

    .footer {
      text-align: center;
      color: var(--text-muted);
      margin-top: 40px;
      font-size: 0.9rem;
    }

    @media (max-width: 640px) {
      .header h1 { font-size: 2rem; }
      .card { padding: 24px 20px; }
      .button-group { flex-direction: column; }
      .result-header { flex-direction: column; align-items: flex-start; gap: 8px; }
    }
  </style>
</head>
<body>
  <div class="container">
    <div class="header">
      <h1>URL Shortener</h1>
      <p>Transform your long URLs into short, shareable links</p>
    </div>

    <div class="card">
      <div class="input-group">
        <label for="url">Enter your long URL</label>
        <input id="url" type="url" placeholder="https://example.com/very/long/url" />
      </div>

      <div class="input-group">
        <label for="code">Custom short code (optional)</label>
        <input id="code" type="text" placeholder="my-custom-code" />
      </div>

      <div class="button-group">
        <button id="shorten" class="btn-primary">Shorten URL</button>
        <button id="toggle-list" class="btn-secondary">View All URLs</button>
      </div>

      <div id="result" class="result-card">
        <div class="result-header">
          <h3>Your shortened URL is ready</h3>
        </div>
        <div class="short-url">
          <a id="short-link" href="#" target="_blank"></a>
          <button class="copy-btn" id="copy-btn">Copy</button>
        </div>
        <div class="original-url" id="original-url"></div>
      </div>

      <div id="error" class="error-message"></div>

      <div class="loading" id="loading">
        <div class="spinner"></div>
      </div>
    </div>

    <div id="url-list-section" class="url-list">
      <div class="card">
        <h2 style="margin-bottom: 20px; color: var(--text);">All Shortened URLs</h2>
        <div class="stats" id="stats"></div>
        <div id="url-list-content"></div>
      </div>
    </div>

    <div class="footer">
      Built with Go • Fast & Reliable • Open Source
    </div>
  </div>

  <script>
    const urlInput = document.getElementById('url');
    const codeInput = document.getElementById('code');
    const shortenBtn = document.getElementById('shorten');
    const toggleListBtn = document.getElementById('toggle-list');
    const resultCard = document.getElementById('result');
    const errorMsg = document.getElementById('error');
    const loading = document.getElementById('loading');
    const shortLink = document.getElementById('short-link');
    const originalUrl = document.getElementById('original-url');
    const copyBtn = document.getElementById('copy-btn');
    const urlListSection = document.getElementById('url-list-section');
    const urlListContent = document.getElementById('url-list-content');
    const statsDiv = document.getElementById('stats');

    let listVisible = false;

    function showError(message) {
      errorMsg.textContent = message;
      errorMsg.classList.add('show');
      resultCard.classList.remove('show');
      setTimeout(() => errorMsg.classList.remove('show'), 5000);
    }

    function showResult(data) {
      shortLink.href = data.short_url;
      shortLink.textContent = data.short_url;
      originalUrl.textContent = '→ ' + data.original_url;
      resultCard.classList.add('show');
      errorMsg.classList.remove('show');
      copyBtn.textContent = 'Copy';
      copyBtn.classList.remove('copied');
    }

    async function shorten() {
      const url = urlInput.value.trim();
      if (!url) {
        showError('Please enter a URL');
        return;
      }

      const body = { url: url };
      const customCode = codeInput.value.trim();
      if (customCode) body.custom_code = customCode;

      loading.style.display = 'block';
      resultCard.classList.remove('show');
      errorMsg.classList.remove('show');

      try {
        const res = await fetch('/shorten', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify(body)
        });

        const data = await res.json();
        loading.style.display = 'none';

        if (res.status >= 400) {
          showError(data.error || 'Failed to shorten URL');
          return;
        }

        showResult(data);
        urlInput.value = '';
        codeInput.value = '';

        if (listVisible) {
          loadUrls();
        }
      } catch (e) {
        loading.style.display = 'none';
        showError('Network error. Please try again.');
      }
    }

    async function loadUrls() {
      try {
        const res = await fetch('/api/urls');
        const data = await res.json();

        statsDiv.innerHTML = '';
        const totalClicks = data.urls.reduce((sum, u) => sum + u.clicks, 0);

        statsDiv.innerHTML = '<div class="stat-card">' +
          '<div class="stat-value">' + data.count + '</div>' +
          '<div class="stat-label">Total URLs</div>' +
          '</div>' +
          '<div class="stat-card">' +
          '<div class="stat-value">' + totalClicks + '</div>' +
          '<div class="stat-label">Total Clicks</div>' +
          '</div>';

        urlListContent.innerHTML = '';

        if (data.urls.length === 0) {
          urlListContent.innerHTML = '<p style="text-align:center;color:var(--text-muted);padding:40px;">No URLs yet. Create your first shortened URL above!</p>';
          return;
        }

        data.urls.sort((a, b) => new Date(b.created_at) - new Date(a.created_at));

        data.urls.forEach(url => {
          const date = new Date(url.created_at).toLocaleString();
          const item = document.createElement('div');
          item.className = 'url-item';
          item.innerHTML = '<div class="url-item-header">' +
            '<span class="url-item-code">/' + url.short_code + '</span>' +
            '<span class="url-item-clicks">' + url.clicks + ' clicks</span>' +
            '</div>' +
            '<div class="url-item-original">' + url.original_url + '</div>' +
            '<div class="url-item-date">Created: ' + date + '</div>';
          urlListContent.appendChild(item);
        });
      } catch (e) {
        showError('Failed to load URLs');
      }
    }

    function toggleList() {
      listVisible = !listVisible;
      if (listVisible) {
        urlListSection.classList.add('show');
        toggleListBtn.textContent = 'Hide URLs';
        loadUrls();
      } else {
        urlListSection.classList.remove('show');
        toggleListBtn.textContent = 'View All URLs';
      }
    }

    async function copyToClipboard() {
      try {
        await navigator.clipboard.writeText(shortLink.textContent);
        copyBtn.textContent = 'Copied!';
        copyBtn.classList.add('copied');
        setTimeout(() => {
          copyBtn.textContent = 'Copy';
          copyBtn.classList.remove('copied');
        }, 2000);
      } catch (e) {
        showError('Failed to copy to clipboard');
      }
    }

    shortenBtn.addEventListener('click', shorten);
    toggleListBtn.addEventListener('click', toggleList);
    copyBtn.addEventListener('click', copyToClipboard);

    urlInput.addEventListener('keydown', (e) => {
      if (e.key === 'Enter') shorten();
    });

    codeInput.addEventListener('keydown', (e) => {
      if (e.key === 'Enter') shorten();
    });
  </script>
</body>
</html>`
