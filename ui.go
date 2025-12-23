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
<html>
<head>
  <meta charset="utf-8">
  <title>URL Shortener</title>
  <style>
    body{font-family:system-ui,Segoe UI,Roboto,Helvetica,Arial;max-width:720px;margin:40px auto;padding:0 16px}
    input,button{font-size:16px;padding:8px}
    .row{display:flex;gap:8px;margin-bottom:8px}
    .result{margin-top:12px;word-break:break-all}
  </style>
</head>
<body>
  <h1>URL Shortener</h1>
  <div>
    <div class="row">
      <input id="url" type="url" placeholder="https://example.com" style="flex:1" />
      <input id="code" type="text" placeholder="custom code (optional)" />
      <button id="shorten">Shorten</button>
    </div>
    <div>
      <button id="list">Show stored URLs</button>
    </div>
    <div class="result" id="result"></div>
    <pre id="listout" style="display:none;background:#f7f7f7;padding:12px;border-radius:6px"></pre>
  </div>
  <script>
    async function shorten() {
      const url = document.getElementById('url').value;
      const code = document.getElementById('code').value;
      if (!url) { alert('Enter a URL'); return; }
      const body = { url: url };
      if (code) body.custom_code = code;
      const res = await fetch('/shorten', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(body)
      });
      const text = await res.text();
      try {
        const j = JSON.parse(text);
        if (res.status >= 400) {
          document.getElementById('result').innerText = 'Error: ' + (j.error || text);
          return;
        }
        const a = document.createElement('a');
        a.href = j.short_url;
        a.textContent = j.short_url;
        a.target = '_blank';
        const cont = document.getElementById('result');
        cont.innerHTML = '';
        cont.appendChild(a);
        const meta = document.createElement('div');
        meta.textContent = ' → ' + j.original_url;
        cont.appendChild(meta);
      } catch (e) {
        document.getElementById('result').innerText = text;
      }
    }
    async function listUrls() {
      const res = await fetch('/api/urls');
      const j = await res.json();
      const out = document.getElementById('listout');
      out.style.display = 'block';
      out.textContent = JSON.stringify(j, null, 2);
    }
    document.getElementById('shorten').addEventListener('click', shorten);
    document.getElementById('url').addEventListener('keydown', (e)=>{ if(e.key==='Enter') shorten(); });
    document.getElementById('list').addEventListener('click', listUrls);
  </script>
</body>
</html>`
