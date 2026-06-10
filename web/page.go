package web

const indexHTML = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>RedactLog</title>
<style>
  body { font-family:-apple-system,system-ui,sans-serif; background:#0f172a;
         color:#e2e8f0; margin:0; padding:40px; line-height:1.6; }
  .wrap { max-width:780px; margin:0 auto; }
  h1 { font-size:1.9rem; margin-bottom:2px; }
  .tag { color:#94a3b8; margin-bottom:24px; }
  label { font-size:0.85rem; color:#94a3b8; display:block; margin:14px 0 6px; }
  textarea { width:100%; box-sizing:border-box; background:#1e293b; color:#e2e8f0;
             border:1px solid #334155; border-radius:8px; padding:12px;
             font-family:ui-monospace,monospace; font-size:0.9rem; min-height:120px; }
  .row { display:flex; gap:12px; align-items:center; margin:14px 0; }
  select, button { background:#1e293b; color:#e2e8f0; border:1px solid #334155;
                   border-radius:8px; padding:9px 14px; font-size:0.9rem; cursor:pointer; }
  button { background:#16a34a; border:none; font-weight:600; }
  button:hover { background:#15803d; }
  .out { background:#0b1220; border:1px solid #334155; border-radius:8px;
         padding:14px; font-family:ui-monospace,monospace; font-size:0.9rem;
         white-space:pre-wrap; word-break:break-word; min-height:60px; }
  .count { color:#16a34a; font-size:0.85rem; margin-top:8px; }
  .ex { color:#7dd3fc; cursor:pointer; font-size:0.85rem; }
  a { color:#7dd3fc; }
</style>
</head>
<body>
  <div class="wrap">
    <h1>🔒 RedactLog</h1>
    <p class="tag">Privacy-preserving PII redaction — runs on-device, scrubs
       emails, SSNs, cards, tokens, phones &amp; IPs from logs.</p>

    <label>Log text (paste anything with PII)</label>
    <textarea id="in">user jane@acme.com from 192.168.1.1 paid with card 4111 1111 1111 1111, ssn 555-12-3456, token sk_live_abcdefghij1234567890</textarea>
    <span class="ex" onclick="loadExample()">↺ reset example</span>

    <div class="row">
      <label style="margin:0">Mode</label>
      <select id="mode">
        <option value="mask">mask — [REDACTED:TYPE]</option>
        <option value="hash">hash — stable hash</option>
        <option value="drop">drop — remove</option>
      </select>
      <button onclick="redact()">Redact</button>
    </div>

    <label>Redacted output</label>
    <div class="out" id="out">—</div>
    <div class="count" id="count"></div>

    <p style="margin-top:28px; color:#64748b; font-size:0.85rem;">
      Source: <a href="https://github.com/siriscent7/redactlog">github.com/siriscent7/redactlog</a>
      · also runs as a streaming CLI: <code>app | redactlog</code>
    </p>
  </div>

<script>
const EXAMPLE = "user jane@acme.com from 192.168.1.1 paid with card 4111 1111 1111 1111, ssn 555-12-3456, token sk_live_abcdefghij1234567890";
function loadExample(){ document.getElementById('in').value = EXAMPLE; }
async function redact(){
  const text = document.getElementById('in').value;
  const mode = document.getElementById('mode').value;
  const res = await fetch('/api/redact', {
    method:'POST', headers:{'Content-Type':'application/json'},
    body: JSON.stringify({text, mode})
  });
  const data = await res.json();
  document.getElementById('out').textContent = data.redacted || '—';
  document.getElementById('count').textContent = data.count + ' PII item(s) redacted';
}
redact();
</script>
</body>
</html>`
