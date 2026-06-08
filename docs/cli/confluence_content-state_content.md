# confluence content-state content

List Cloud content with a given state

## Synopsis

List Cloud content in one space that has the provided content state.

## Examples

confluence content-state content ENG --state-id 1
  confluence content-state content ENG --state-id 1 --expand space --expand version --json
  confluence content-state content ENG --state-id 0 --start 25 --limit 25

## Usage

```text
confluence content-state content <space-key> [flags]
```

## Options

```text
      --expand strings   Expand value; repeatable or comma-separated
      --limit int        Max content rows (Cloud endpoint hard cap 100) (default 25)
      --start int        Zero-based result offset
      --state-id int     Content state id; required, and 0 is valid for Cloud default states
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

- [confluence content-state](confluence_content-state.md)
