# 🔒 RedactLog

**Privacy-preserving PII redaction that runs entirely on-device, at the log line.**

Logs leak sensitive data — emails, SSNs, credit cards, API tokens, phone numbers, IPs. RedactLog scrubs PII from log streams *before* they reach storage, using a single-scan multi-pattern engine. All processing is local; no data leaves the machine.

🔗 **Live demo:** https://redactlog.onrender.com

---

## The problem
Application logs routinely capture PII and secrets. Once written to centralized log storage, that sensitive data is exposed to anyone with log access — a real compliance and breach risk. RedactLog redacts PII inline, on-device, so secrets never hit the log store in the first place.

## What it does
- **6 PII types** — SSN, credit card, email, API token (`sk_/pk_/tok_`), phone, IP.
- **3 redaction modes:**
  - `mask` → `[REDACTED:EMAIL]`
  - `hash` → `[EMAIL:a1b2c3d4]` (stable hash — correlate without exposing the value)
  - `drop` → removed entirely
- **Single-scan engine** — all patterns compiled into one regex with named groups, so each line is scanned once (not once per pattern).
- **Streaming CLI** — reads stdin line-by-line; drop it into any log pipeline:
  `app | redactlog | logstore`.
- **On-device** — no network calls, no external services. Privacy by design.

## Demo
```bash
$ echo "user jane@acme.com from 192.168.1.1 paid with card 4111 1111 1111 1111, ssn 555-12-3456" | redactlog user [REDACTED:EMAIL] from [REDACTED:IP] paid with card [REDACTED:CREDIT_CARD], ssn [REDACTED:SSN] 
redacted 4 PII items
```

## Correctness
- 100% recall on a synthetic PII test set covering all 6 types.
- A recall test guards against regressions — a missed pattern in a redaction tool is a data leak, so recall is measured, not assumed.

## Performance
Benchmarked on Apple M3 Pro:
- ~12.6 MB/s on PII-dense input; ~7 MB/s on realistic logs.
- Optimization reduced allocations from 29 → 11 per line by collapsing six regex passes into a single combined-regex scan.
- Throughput is bounded by Go's RE2 regex engine, which guarantees linear-time matching (no catastrophic backtracking / ReDoS) — a deliberate safety trade-off for processing untrusted log input.

## Run it
```bash
git clone https://github.com/siriscent7/redactlog.git
cd redactlog

# pipe logs through it
echo "email a@b.com ssn 111-22-3333" | go run main.go

# choose a mode
echo "login a@b.com a@b.com" | go run main.go -mode hash
echo "secret sk_live_abcdefghij1234567890" | go run main.go -mode drop

# or build a binary
go build -o redactlog .
tail -f /var/log/app.log | ./redactlog
```

## Tests & benchmark
```bash
go test ./redactor -v                    # 7 tests incl. recall
go test ./redactor -bench=. -benchmem    # throughput + allocations
```

## Architecture
```
log stream (stdin)
   │  line by line
   ▼
combined regex (one scan, named groups per PII type)
   │
   ▼
mode handler ── mask | hash | drop
   │
   ▼
clean log (stdout)
```

## Tech stack
Go · RE2 regex · streaming I/O

## Limitations
- Regex-based detection only catches anticipated formats. Novel PII shapes (free-text names, unusual ID formats) are missed — the fundamental limit of pattern matching.
- Throughput is bounded by RE2; not optimized for multi-GB/s pipelines.

## Future Work
- NER/ML tier — a second pass with named-entity recognition to catch free-text PII (names, addresses) that regex can't.
- Specialized scanners — Aho-Corasick or hyperscan bindings for higher throughput on fixed patterns.
- Concurrent processing — fan out lines across goroutines for multi-core throughput.
- Configurable patterns — user-defined PII rules via a config file.