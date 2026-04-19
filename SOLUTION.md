# Niche Full-Stack Exercise — Solution

## Table of contents

- [How to run](#how-to-run)
- [Front-end — key decisions](#front-end--key-decisions)
- [Back-end — key decisions](#back-end--key-decisions)
- [API contract](#api-contract)
  - [`GET /autocomplete`](#get-autocomplete)
  - [`GET /reviews`](#get-reviews)
  - [What I'd improve in the contract](#what-id-improve-in-the-contract)
- [AI workflow — what I did well](#ai-workflow--what-i-did-well)
- [Example prompts (Go tests + Go code)](#example-prompts-go-tests--go-code)
  - [1. Data-engineer — profile the CSV before writing code](#1-data-engineer--profile-the-csv-before-writing-code)
  - [2. Go — table-driven test for the CSV loader](#2-go--table-driven-test-for-the-csv-loader)
  - [3. Go — idiomatic error wrapping on Close](#3-go--idiomatic-error-wrapping-on-close)
  - [4. Front-end — race-safe autocomplete](#4-front-end--race-safe-autocomplete)
- [Trade-offs and what I'd do with more time](#trade-offs-and-what-id-do-with-more-time)
- [Where AI helped most / least](#where-ai-helped-most--least)

---

## How to run

**Back-end** (Go):

```
cd back-end
go mod tidy
go run .               # listens on :8080
go test ./... -race    # tests
```

**Front-end** (Node):

```
cd front-end
npm install
npm run dev            # http://localhost:3000
npm run build          # production bundle
```

Environment: `front-end/.env.local` sets `NEXT_PUBLIC_API_BASE=http://localhost:8080`. Override `API_BASE` if the Next server should reach the Go service over a different hostname than the browser does.

---

## Front-end — key decisions

- **Next.js App Router with a real route per college (`/colleges/[slug]`).** Selecting an autocomplete item navigates to a URL — links are shareable, back/forward works, and `generateMetadata` gives each page a proper `<title>`. A modal or inline panel would have been simpler but would have thrown all that away.
- **Server Component for the college page, Client Component only for search.** `CollegePage` is `async` and calls `fetchReviews` on the server with Next's `{ next: { revalidate: 60 } }` cache hint. The search box in the header is the only island of interactivity — smaller JS payload, cacheable review pages.
- **Debounced autocomplete with race-safe requests.** `useDebouncedValue` (200 ms) plus `AbortController` plus a `latestRequestId` ref in `Search.tsx`. Out-of-order responses get discarded instead of clobbering the list.
- **Combobox ARIA pattern over a bespoke dropdown.** `role="combobox"`, `aria-controls`, `aria-activedescendant`, `role="listbox"`/`option`, `aria-selected`. Keyboard (↑/↓ wrap, Enter selects, Escape closes) and screen-reader users both work.
- **CSS Modules, no runtime CSS-in-JS.** Component-scoped class names without shipping a styling runtime. Typography uses `next/font` with three CSS variables (`Inter`, `Fraunces`, `JetBrains_Mono`) to expose a tiny design system.
- **A distinct `NotFoundError` in the API layer.** `CollegePage` calls `notFound()` for 404s but re-throws anything else to the nearest `error.tsx`. No silent swallowing of real bugs.
- **Separate `API_BASE` (server) vs `NEXT_PUBLIC_API_BASE` (client).** Lets the server talk to the Go service over an internal hostname in production without leaking it to the client bundle.
- **Typed API surface (`College`, `ReviewsResponse`) in one file.** The shape of the Go JSON response is the contract, imported by every caller.

## Back-end — key decisions

- **Single CSV load at startup, two data structures.** `Reviews map[string][]string` for O(1) lookup by URL, plus a pre-sorted `Colleges []College` slice for autocomplete. No index rebuild per request.
- **Prefix-match autocomplete, case-insensitive, default limit 20.** `strings.HasPrefix` on pre-lowercased names, `?limit=` is configurable. Simple and predictable; I'd reach for a trie only if N went into the millions.
- **Dash-collapsing on URL normalization.** `academy----of-art-university` and `academy--of-art-university` both collapse to `academy-of-art-university` so rows with malformed slugs dedupe instead of creating ghost colleges. Applied both on load and on inbound `?url=` lookup.
- **Skip malformed rows, don't fail the load.** Rows with missing `COLLEGE_NAME` or `COLLEGE_URL` are logged (with the UUID) and dropped, instead of blowing up the service on one bad row.
- **Strict header validation.** `loadReviewsFile` refuses to start if the CSV header doesn't match the four expected columns. Catches schema drift immediately rather than silently misreading columns.
- **Empty results serialize as `[]`, not `null`.** `handleAutocomplete` explicitly swaps `nil` for `[]College{}` before encoding, so the front-end never has to write `?? []`.
- **`errors.Join` on deferred close errors plus `%w` wrapping everywhere.** Close-after-read errors don't get lost, and every boundary uses `fmt.Errorf("... %w", err)` so callers can `errors.Is`/`errors.As` up the stack.
- **Stdlib only — `net/http` + `http.ServeMux`.** The new `GET /path` pattern covers what we need. Bringing in chi/gin for two endpoints would have been ceremony.
- **`TestMain` redirects `log.Output` to `io.Discard`.** Keeps `go test -v` clean of "skipping row…" lines from the CSV-loader tests without losing logs in production.
- **Table-driven tests with `httptest.NewRecorder` for handlers.** One table for the loader (header, encoding, dedup, normalization, malformed rows), one for the HTTP layer (routing, 4xx, null-vs-`[]`). Adding a new edge case is a one-line diff.

---

## API contract

Two endpoints. Both return JSON, set `Content-Type: application/json` and `Access-Control-Allow-Origin: *`. The Go types are the source of truth:

```go
type College struct {
    Name string `json:"name"`
    URL  string `json:"url"`
}
```

### `GET /autocomplete`

**Query params**

| Param   | Type   | Required | Default | Notes                                                |
|---------|--------|----------|---------|------------------------------------------------------|
| `q`     | string | no       | `""`    | Prefix to match against `name`, case-insensitive.    |
| `limit` | int    | no       | `20`    | Must be a positive integer if present.               |

**Responses**

- `200 OK` — `[]College`. Empty/missing/no-match `q` returns `[]`, **not** `null`. Order is alphabetical, case-insensitive by `name`.
- `400 Bad Request` — `limit` is non-numeric, zero, or negative. Body: `{"error":"limit must be a positive integer"}`.

**Example**

```
GET /autocomplete?q=alp&limit=2

200 OK
[
  { "name": "Alpha University", "url": "alpha" },
  { "name": "alpha State",      "url": "alpha-state" }
]
```

### `GET /reviews`

**Query params**

| Param | Type   | Required | Notes                                                       |
|-------|--------|----------|-------------------------------------------------------------|
| `url` | string | yes      | College slug. Repeated dashes are collapsed before lookup.  |

**Responses**

- `200 OK` — `{ "college": College, "reviews": string[] }`. `reviews` is never null; it may be `[]`.
- `400 Bad Request` — `url` is missing. Body: `{"error":"missing url query param"}`.
- `404 Not Found` — no college matches that slug. Body: `{"error":"college not found"}`.

**Example**

```
GET /reviews?url=alpha

200 OK
{
  "college": { "name": "Alpha University", "url": "alpha" },
  "reviews": [
    "Great school.",
    "Loved the campus."
  ]
}
```

### What I'd improve in the contract

- **Typed error envelope.** Error bodies are ad-hoc JSON (`{"error":"..."}`). A shared shape — `{ "error": { "code": "LIMIT_INVALID", "message": "..." } }` — lets the client branch on `code` instead of string-matching, and survives copy changes.
- **Paginate `/autocomplete` and `/reviews`.** `limit` caps the slice but there's no `offset`/`cursor`, and `/reviews` returns the entire list. For colleges with hundreds of reviews, the payload grows unboundedly. Cursor-based pagination (`?cursor=<opaque>&limit=20`) plus a `nextCursor` field in the response is the fix.
- **Richer `College` and review shape.** `College` is just `{name, url}`. Production would want `id` (stable UUID), `location`, `logoUrl`, and at least a review count. Reviews are bare strings — they should carry at minimum `id`, `rating`, `createdAt`, maybe `author`. That's what a real CSV schema would give us; the current exercise data is already narrow.
- **POST-style filters instead of overloaded query strings.** Once filters grow (by state, size, rating range), `GET` with a dozen query params gets unreadable. `POST /search` with a JSON body is clearer — and easier to version.
- **Versioning.** Everything is unversioned (`/autocomplete`, `/reviews`). Prefixing with `/v1/` is cheap now and expensive later.
- **Consistent slug handling at the edge.** The server accepts `academy----of-art-university` and normalizes internally; I'd rather return `308 Permanent Redirect` to the canonical slug, so caches and clients converge on one URL per college.
- **ETag / `Cache-Control`.** Review data is read-mostly. Setting `ETag` + `Cache-Control: public, max-age=60, stale-while-revalidate=300` would offload a huge fraction of traffic to the CDN and browser cache — right now every request hits Go.
- **CORS should not be `*` in production.** Fine for the exercise; real deployments restrict to the front-end origin, and set `Vary: Origin` so caches don't leak responses across origins.
- **Observability.** No request logging, no metrics, no trace IDs. A tiny middleware chain (request-id, structured `slog` access log, Prometheus counters on `{endpoint, status}`) would make the service debuggable without code changes.
- **Rate limiting.** `/autocomplete` fires on every keystroke; a sloppy client can hammer it. Token-bucket per IP (e.g. 20 rps) is table stakes.
- **OpenAPI spec.** A hand-written `openapi.yaml` (or generated from struct tags) would let the front-end generate a typed client and keep the contract honest as it evolves.

---

## AI workflow — what I did well

Framed around **what the git history actually shows**, not claims:

- **Small, single-purpose commits with conventional messages.** `c00e58b fix(reviews): skip rows with missing name or url`, `7c54ac6 fix(reviews): collapse repeated dashes in slugs`, `cecab4a feat(autocomplete): add GET /autocomplete with prefix search`. Each commit is reviewable on its own and reverts cleanly.
- **Tests landed with the code, not after.** The first feature commit (`597a37c feat(reviews): add csv loader with tests`) shipped with its test file; every subsequent fix added a case to the existing table. No "tests later" commits.
- **AI output is a first draft, not the answer.** Every suggestion went through review: re-read the diff, run `go test ./... -race` and `npm run build`, then commit. When the model reached for an unidiomatic pattern (bespoke dropdown state instead of combobox ARIA, a one-off closure instead of a hook), I pushed back and asked for the idiomatic version.
- **Scope discipline.** When AI proposed incidental refactors (tidying files I hadn't touched, renaming variables I hadn't asked about), I declined and kept the change to the asked scope. It's why the diffs are readable.
- **Explicit decisions over implicit defaults.** When the model wanted to silently swallow a malformed CSV row, I had it `log.Printf` the UUID instead. When it wanted to return `null` on an empty autocomplete, I had it return `[]`. Each small call closes a foot-gun on the consumer side.
- **Tight prompts, not essays.** I named the file, the function, the failure mode, the assertion style, and the exact error substrings I wanted — so the model produced something mergeable instead of generic boilerplate. See below.

---

## Example prompts (Go tests + Go code)

Four prompts that did real work on this exercise. Each is followed by one line on why it worked. Prompt 1 is the one I ran **before writing any Go code** — it's what surfaced the data quality issues that later became the `fix(reviews): skip rows with missing name or url` and `fix(reviews): collapse repeated dashes in slugs` commits.

### 1. Data-engineer — profile the CSV before writing code

```
Act as a data engineer. Before I write any Go loader code, I want a profile
of back-end/data/niche_reviews.csv so I know what edge cases the loader
must handle.

Run a quick inspection (cat, awk, python — whatever is shortest) and report:

 1. Row count, column count, and whether every row has exactly 4 columns.
 2. Header: confirm columns are COLLEGE_UUID, COLLEGE_NAME, COLLEGE_URL,
    REVIEW_TEXT in that order. Flag any deviation.
 3. Null / empty / whitespace-only cells per column. Give counts.
 4. Duplicate COLLEGE_UUIDs. Do the same (name, url) pair ever appear under
    different UUIDs, or vice versa?
 5. COLLEGE_URL hygiene: any URLs with leading/trailing dashes, repeated
    dashes ("--", "---", "----"), whitespace, uppercase, or non-ASCII?
    Show the top 10 weirdest examples verbatim.
 6. Do colleges with the same NAME ever map to different URL slugs? If so,
    cluster them and show 5 examples.
 7. REVIEW_TEXT: min / median / max length, presence of embedded newlines,
    embedded commas, escaped quotes. Any empty reviews?
 8. Character encoding surprises (BOM, non-UTF8 bytes, smart quotes).

For each anomaly, tell me:
 - what it is,
 - how many rows are affected,
 - the concrete decision the loader needs to make (skip? normalize? error?),
 - which test case it turns into.

Do not propose Go code yet. I want the *data surface* first.
```

*What it caught:* the repeated-dash URL collisions (`academy----of-art-university` vs `academy--of-art-university` vs `academy-of-art-university` all pointing at the same college), a handful of rows with empty names/URLs, and embedded newlines inside quoted review text. Each finding became a one-line test case in the table and a documented loader rule — nothing discovered later by accident in `go test`.

*Why it worked:* framing the model as a **data engineer** with a fixed reporting template kept it from jumping straight to code. The numbered checklist forced it to actually run commands (not guess), and the "no Go yet" rule meant the output was a decision list, not a draft implementation.

### 2. Go — table-driven test for the CSV loader

```
Write a table-driven test in Go for loadReviewsFile in back-end/reviews.go.
Use t.Parallel at both the outer and subtest level. Use t.TempDir to write
CSV fixtures. Cover these cases:
  - happy path with two colleges, second one has two reviews (dedup check)
  - embedded comma and escaped double-quote inside REVIEW_TEXT
  - multi-line quoted review (newline inside quotes)
  - valid header, zero rows
  - wrong header -> error substring "invalid header"
  - short row (3 cols) -> error substring "read row"
  - rows with missing name, missing url, both missing, whitespace-only
    name/url are all skipped and do not pollute later valid rows
  - repeated dashes in the URL ("academy----of-art-university") normalize
    to single dashes and dedupe with the canonical form
  - colleges sorted case-insensitive by name

Use reflect.DeepEqual for the ReviewsData comparison, wantErr as a substring
match, and a writeCSV(t, content) helper that returns the tempfile path.
Do not mock the filesystem — write a real CSV to t.TempDir().
```

*Why it worked:* named helpers, patterns, and the assertion style. The model had nowhere to drift to, so the output was mergeable on the first pass.

### 3. Go — idiomatic error wrapping on Close

```
In loadReviewsFile I want defer file.Close() to stop silently dropping close
errors. Rewrite using a named return (err), and in the deferred func do
errors.Join(err, fmt.Errorf("close reviews csv: %w", closeErr)) when
closeErr is non-nil. Keep the rest of the function signature and behavior
identical.
```

*Why it worked:* small, surgical, single-purpose. No "rewrite this function" ambiguity.

### 4. Front-end — race-safe autocomplete

```
In front-end/src/components/Search.tsx the autocomplete fires a fetch per
debounced keystroke. I need out-of-order responses to be discarded. Use
AbortController for cancellation AND a latestRequestId ref as a secondary
guard (compare in the .then handler, bail if stale). Debounce input with a
useDebouncedValue hook (value + delayMs -> debounced value, cleared on
unmount). Don't change Autocomplete's props — Search should still pass
{ query, onQueryChange, items, isSettled }.
```

*Why it worked:* named the failure mode (out-of-order responses), the two defenses, and fixed the component contract — so the diff stayed small.

---

## Trade-offs and what I'd do with more time

- **CSV is re-read at startup and held in memory.** Fine for this dataset. For a much larger corpus I'd stream from disk, swap the sorted slice for a trie, and move to a Radix/FST for prefix search.
- **Autocomplete is prefix-only.** A real product wants fuzzy/substring matching and ranking by popularity; I'd bring in a proper index (e.g., Bleve) before rolling my own scorer.
- **No cross-service integration test.** The handler tests exercise the router and the front-end is verified in the browser; a Playwright smoke test hitting both services would be the next add.
- **No rate limiting or request metrics.** Two CORS-wide-open endpoints are fine for the exercise; production would want a middleware stack (logging, rate limit, tracing).

## Where AI helped most / least

**Most:** writing the table-driven test scaffolds, the combobox ARIA wiring, and the boilerplate around `generateMetadata` / `notFound()`. These are patterns with one "right" shape — a well-scoped prompt produces mergeable code fast.

**Least:** data-shape decisions. The choice to collapse repeated dashes, return `[]` over `null`, and split `/colleges/[slug]` as a real route all came from me reading the CSV and thinking about the consumer, not from prompting. The model will pick a plausible default if you let it — the value I added was not letting it.
