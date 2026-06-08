# confluence search all

Unified search — pages + spaces + users + attachments in parallel

## Synopsis

Fan-out page + space + user + attachment search in parallel, merge via
reciprocal-rank fusion (k=60). Partial failures log a warning to stderr and
the command continues with the healthy branches.

## Output schema (stable)

[{"kind":"page|space|user|attachment","title":"...","url":"...",
    "snippet":"...","score":0.0123,"source":{...raw...}}, ...]

## Examples

confluence search all "release process" --limit 20 --json
  confluence search all "deploy" --space ENG --branch-timeout 800ms
  confluence search all "<q>" --jq '.[] | {kind, title, space}'

## Usage

```text
confluence search all <query> [flags]
```

## Options

```text
      --branch-timeout duration   Cancel a branch after this duration (0 disables — parent context only) (default 1.2s)
      --limit int                 Max results per branch (hard cap 200) (default 25)
      --space string              Optional space key filter for page/attachment branches
```

## Inherited Options

```text
      --debug             Log HTTP requests to stderr (Authorization header redacted)
      --jq string         Filter JSON output through a jq expression (requires jq on PATH)
      --json              Output as JSON (machine-readable)
      --template string   Render JSON output through a Go text/template
      --timing            Show execution time on stderr
```

## See Also

- [confluence search](confluence_search.md)
