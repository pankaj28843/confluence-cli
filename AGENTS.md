# AGENTS.md — confluence-cli

Notes for AI coding agents working on this repo.

## What this is

`confluence` is a Go CLI that wraps **Atlassian Confluence** REST APIs —
both **Server / Data Center** (`/rest/api/...`, Bearer PAT) and **Cloud**
(`/wiki/rest/api/...` v1 and `/wiki/api/v2/...` v2, Basic email+API
token). It is open-source, MIT-licensed, and intended to compose with
`gh`, `docsearch`, and other shell tools.

## Non-negotiables

- **Auth is flavor-aware**:
  - Server/DC → `Authorization: Bearer <PAT>`
  - Cloud    → `Authorization: Basic base64(email + ":" + api_token)`
  Never send Bearer on Cloud; never send Basic email:token on Server/DC.
- **Only runtime deps**: `spf13/cobra` (+ indirect `pflag`/`mousetrap`) and
  `PuerkitoBio/goquery` (for the Markdown converter). Don't add more
  without a written justification.
- **Authorization is never logged**. The debug round-tripper uses
  `req.URL.Redacted()` and skips the `Authorization` header.
- **CA bundles**: respect `SSL_CERT_FILE` / `REQUESTS_CA_BUNDLE` only.
  Do NOT auto-load vendor-specific CA paths.
- **OSS safety**: no company-specific hostnames, internal project codes,
  or employer names anywhere in the repo. `wiki.example.com` and
  `example.atlassian.net` are the sanctioned example domains.

## Layout

```
cmd/confluence/     cobra entry + one file per command group, plus
                    stubs.go that wires live groups to their *Real() impls
internal/client/    HTTP + dual-flavor auth + retry + debug
internal/conf/      Confluence resource helpers + storage→Markdown
internal/output/    --json / --jq / --template / --timing writer
scripts/pre-commit  gofmt, vet, test (install via `make setup`)
```

## House style

- `Workflow:` block on the root; `Examples:` block on every leaf verb.
- Every command accepts `--json`, `--jq`, `--template`, `--timing`,
  `--debug` (inherited from root).
- Errors from env/auth use `newConfigError(err)` → exit 2. Everything else
  is exit 1.
- Use `signal.NotifyContext(context.Background(), os.Interrupt)` via
  `newContext()` — never a bare `context.Background()` inside a `RunE`.
- Default `--limit` is 25; hard cap is 200.
- Use the `Content` struct from `internal/conf/content.go` for
  page/blogpost/comment/attachment — do NOT mint new parallel structs.

## Tests

```bash
make test   # -count=1 -timeout 60s
make race   # -race -count=1
```

Tests mock HTTP via `httptest.NewServer`. Both flavors (Server/DC +
Cloud) must be exercised for any new client-layer change — see
`internal/client/client_test.go` for the pattern.

## Writing new verbs

1. Add a helper in `internal/conf/<resource>.go` with a unit test.
2. Add the cobra leaf in `cmd/confluence/<group>.go` as `*Real()`.
3. Wire the `Real` function into `stubs.go#<group>Cmd`.
4. Wire `--json` via `getWriter().IsJSON() -> w.JSON(resp)`.

## Confluence-specific gotchas

- **Pagination**: Server/DC + Cloud v1 use `?start=N&limit=N` and return
  `_links.next`; Cloud v2 uses cursor pagination in `_links.next`.
- **`expand` parameter**: `body.storage,version,space,ancestors,metadata.labels`
  is the default surface. The Markdown converter reads `body.storage.value`.
- **Comment locations**: `location` is repeatable — use `url.Values.Add()`.
- **Attachment upload**: multipart/form-data, field named `file`, header
  `X-Atlassian-Token: no-check` (handled by `UploadAttachment`).
- **Cloud base URL**: `/wiki` is already part of `CONFLUENCE_URL`; don't
  double-prefix when building Cloud paths.
- **CQL injection**: CQL queries are user input — we pass them through
  verbatim in `api`. Typed helpers (`page search --space ENG`) build
  quoted CQL.

## DO NOT

- Do not commit an API token, PAT, password, or email+token pair.
- Do not log the `Authorization` header.
- Do not hand-build `_links.next` URLs — follow the server's response.
- Do not introduce employer-specific or internal-hostname strings. All
  example URLs must use `wiki.example.com` (Server/DC) or
  `example.atlassian.net` (Cloud). The pre-push grep scans for common
  leaks and blocks the push.
- Do not ship `page create` / `page delete` / `page move` in v1.
- Do not replace `cobra` with another framework.
