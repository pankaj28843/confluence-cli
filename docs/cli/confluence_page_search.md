# confluence page search

CQL-powered page search

## Synopsis

Search pages using CQL. Convenience flags build CQL for you; --cql overrides
everything. Default CQL if no flag: 'type=page'.

## Examples

confluence page search --cql "type=page AND space=ENG AND text ~ 'deploy'"
  confluence page search --space ENG --title "Release notes"
  confluence page search --cql "type=page" --limit 10 --json

## Usage

```text
confluence page search [flags]
```

## Options

```text
      --cql string     CQL expression (overrides --space/--title)
      --limit int      Max results (hard cap 200) (default 25)
      --space string   Space key filter (CQL helper)
      --title string   Title contains (CQL helper)
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

- [confluence page](confluence_page.md)
