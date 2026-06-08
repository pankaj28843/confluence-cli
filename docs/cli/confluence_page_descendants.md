# confluence page descendants

List descendant pages/content under a page

## Synopsis

List descendants under a page. Cloud uses the v2 descendants endpoint.
Server/Data Center walks documented child page routes recursively.

## Examples

confluence page descendants 12345
  confluence page descendants 12345 --depth 2 --limit 100
  confluence page descendants 12345 --json

## Usage

```text
confluence page descendants <id> [flags]
```

## Options

```text
      --depth int   Max descendant depth; 0 means no explicit cap
      --limit int   Max descendants (hard cap 200) (default 50)
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
