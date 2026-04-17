# confluence — an open-source CLI for Atlassian Confluence

Inspired by `gh` — a terminal-native, JSON-first CLI for Atlassian
**Confluence Server / Data Center** and **Confluence Cloud**. Designed
to compose with `gh` (GitHub), `docsearch` (TechDocs), and any other
shell tools you already use with Claude Code / LLM agents.

```
 $ confluence search all "release process" --json --limit 20
   | jq '.[] | {kind, title, space, url}'

 Fan-out: page + space + user + attachment search in parallel,
 merged via reciprocal-rank fusion. Sub-2s warm.
```

- Single Go binary, `cobra`-only runtime dep (+ `goquery` for the
  storage-format → Markdown converter).
- Supports both Confluence flavors: Server/DC with Bearer PAT, Cloud
  with Basic `email:api_token`.
- `--json` / `--jq` / `--template` / `--timing` / `--debug` on every
  command.
- `confluence api <path>` escape hatch for any endpoint not modelled.
- MIT licensed.

## Install

```bash
# Option 1: go install (requires Go 1.25+)
go install github.com/pankaj28843/confluence-cli/cmd/confluence@latest

# Option 2: clone + make
git clone https://github.com/pankaj28843/confluence-cli ~/Code/confluence-cli
cd ~/Code/confluence-cli
make install          # installs to $HOME/.local/bin/confluence
confluence version
```

## Configure

### Server / Data Center (Bearer PAT)

Create a personal access token at
`{your-confluence}/plugins/personalaccesstokens/usertokens.action` (DC 7.9+).

```bash
export CONFLUENCE_URL=https://wiki.example.com
export CONFLUENCE_PAT=<your-pat>
```

### Cloud (Basic email + API token)

Create an API token at
[id.atlassian.com/manage-profile/security/api-tokens](https://id.atlassian.com/manage-profile/security/api-tokens).

```bash
export CONFLUENCE_URL=https://example.atlassian.net/wiki
export CONFLUENCE_EMAIL=you@example.com
export CONFLUENCE_API_TOKEN=<your-api-token>
```

### Optional

```bash
export CONFLUENCE_FLAVOR=server            # override auto-detection
export CONFLUENCE_DEFAULT_SPACE=ENG        # default for --space
export SSL_CERT_FILE=/path/to/ca.pem       # private CA for self-hosted instances
```

### Sanity check

```bash
confluence doctor
# flavor:            server
# baseUrl:           https://wiki.example.com
# authenticated as:  Alice Example (alice)
# OK — confluence is ready to use.
```

Exit codes: `0` healthy, `2` user-fixable config (bad token, missing env),
`1` unexpected error.

## Verb tree

| Group | Verbs |
|---|---|
| `doctor` / `version` | Health + identity + build info |
| `space` | `list`, `view` |
| `page` | `view --markdown`, `search --cql`, `children --recursive`, `ancestors`, `history`, `versions`, `update` |
| `attachment` | `list`, `download`, `upload` |
| `label` | `list`, `add`, `remove` |
| `comment` | `list --locations footer,inline,resolved` |
| `user` | `current`, `view --username\|--key\|--accountId`, `search` |
| `group` | `list`, `members <name>` |
| `watcher` | `content --page`, `space --space` |
| `restriction` | `list --page` |
| `search` | `content`, `spaces`, `users`, `attachments`, `all` (parallel fan-out + RRF) |
| `api` | **Generic REST passthrough — hit ANY Confluence endpoint** |

### Killer feature

```bash
confluence search all "release process" --json --limit 20
# [
#   {"kind":"page","title":"Release process","space":"ENG","score":0.0164,...},
#   {"kind":"space","title":"Release Engineering",...},
#   {"kind":"user","title":"Alice Example",...},
#   {"kind":"attachment","title":"release-notes-v1.pdf","space":"ENG",...}
# ]
```

Four goroutines fan out to the CQL search API in parallel, reciprocal-rank
fuse the results (`k=60`), and a per-branch timeout keeps the merged
response snappy even when one backend lags. Use `--branch-timeout 800ms`
to tighten, `--branch-timeout 0` to disable.

### Escape hatch: `confluence api`

Any endpoint the typed verbs don't cover is still one command away:

```bash
confluence api /rest/api/user/current
confluence api /rest/api/space --param 'limit=10' --jq '.results[].key'

# Cloud v2 (works only on Cloud flavor)
confluence api /wiki/api/v2/spaces --param 'limit=10'
confluence api /wiki/api/v2/pages --param 'space-id=123' --jq '.results[].title'

# POST
echo '{"query":"type=page"}' | \
  confluence api /rest/api/search -X POST --data '@-'

# Any method
confluence api /rest/api/content/12345 -X DELETE
```

Authorisation is attached automatically. The path you supply is appended
to `$CONFLUENCE_URL` verbatim — so `/rest/api/...` hits the v1 surface,
and `/wiki/api/v2/...` hits Cloud v2.

## Workflow examples

```bash
confluence doctor
confluence user current

confluence space list
confluence space view ENG

confluence page search --cql "type=page AND space=ENG AND text ~ 'deploy'"
confluence page view 12345 --markdown
confluence page children 12345 --recursive --json
confluence page ancestors 12345

confluence attachment list --page 12345 --json
confluence attachment download --page 12345 --name logo.png --output ./logo.png
confluence attachment upload --page 12345 --file ./report.pdf

confluence label list --page 12345
confluence label add --page 12345 --label needs-review
confluence label remove --page 12345 --label needs-review

confluence comment list --page 12345 --locations footer,inline --json

confluence page update 12345 --title "New title"
echo "<p>Updated body</p>" | confluence page update 12345 --body-format storage --body-file -

confluence search all "release" --json
```

## Compose with sibling CLIs

`confluence` is designed to pipe cleanly with `gh`, `docsearch`, and
whatever else lives on your `$PATH`.

```bash
# 1) "Find every wiki page about X, then find the PRs that mention those titles."
confluence search all "release process" --json --limit 20 \
  | jq -r '.[] | select(.kind=="page") | .title' \
  | xargs -I{} gh search prs "{}" --json

# 2) "Cross-check a wiki page's technical claims against upstream docs."
confluence page view 12345 --markdown \
  | grep -Eo 'kubectl [a-z-]+' | sort -u \
  | xargs -I{} docsearch search kubernetes "{}" --json

# 3) "For every space I watch, list its root pages."
confluence space list --json --jq '.[].key' \
  | xargs -I{} confluence page search --cql 'type=page AND space="{}" AND parent is EMPTY' --limit 3

# 4) "Download every PDF attachment from a page."
confluence attachment list --page 12345 --json \
  | jq -r '.[] | select(.extensions.mediaType | test("pdf"; "i")) | .title' \
  | xargs -I{} confluence attachment download --page 12345 --name "{}" --output "./{}"
```

## Output

Every command accepts:

- `--json` — emit JSON to stdout (indented).
- `--jq '<expr>'` — pipe JSON through `jq` (`jq` must be on `PATH`).
- `--template '<go-template>'` — render through Go `text/template`.
  Template keys are lowercase JSON field names (`.title`, `.space.key`).
- `--timing` — wall-clock on stderr.
- `--debug` — log HTTP method, URL, status to stderr. Authorization is
  redacted.

## Write operations

v1 ships a **minimum, safe write surface** — writes that let agents fix
their own output without risking destructive mistakes:

- `confluence page update` — title and/or body.
- `confluence label add` / `label remove`.
- `confluence attachment upload`.

Not in v1 (deferred to v2):
- `page create`, `page delete`, `page move` (tree-breaking).
- Whiteboards, databases, blogpost lifecycle.
- OAuth / 3LO / OAuth-bot authentication.

## Testing

```bash
make test       # -count=1 -timeout 60s
make race       # -race
make coverage   # internal/ coverage report
make setup      # install pre-commit hook
make release    # cross-compile darwin-arm64/amd64, linux-amd64, windows-amd64
```

## License

[MIT](./LICENSE). PRs welcome.

## See also

- `docsearch describe confluence-server-dev` — indexed Atlassian DC
  developer docs; great source of truth for endpoint shapes.
- [Atlassian Confluence Server REST API](https://developer.atlassian.com/server/confluence/rest/)
- [Confluence Cloud REST API v2](https://developer.atlassian.com/cloud/confluence/rest/v2/intro/)
- [Advanced Searching using CQL](https://developer.atlassian.com/server/confluence/advanced-searching-using-cql/)
