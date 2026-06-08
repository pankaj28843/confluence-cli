# confluence page direct-children

List direct mixed children of a page

## Synopsis

List direct mixed content-tree children under a page or content id.

Cloud returns page, database, embed, folder, and whiteboard children through the
v2 direct-children route. Server/Data Center returns expanded direct child
content types through the documented child content route.

## Examples

confluence page direct-children 12345
  confluence page direct-children 12345 --type page --type database --json
  confluence page direct-children 12345 --type page,comment --limit 100

## Usage

```text
confluence page direct-children <id> [flags]
```

## Options

```text
      --limit int      Max children (hard cap 200) (default 50)
      --type strings   Content type filter; repeatable or comma-separated
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
