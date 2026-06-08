# confluence watcher content

Watchers subscribed to a content id (raw passthrough)

## Synopsis

Show watcher/notification records for a content id. The shape varies
between Confluence versions, so output is passed through verbatim.

## Examples

confluence watcher content --page 12345 --json

## Usage

```text
confluence watcher content [flags]
```

## Options

```text
      --page string   Content id (required)
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

- [confluence watcher](confluence_watcher.md)
