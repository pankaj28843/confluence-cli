# confluence comment list

List comments (footer, inline, resolved by default)

## Synopsis

List comments on a content id.

## Examples

confluence comment list --page 12345
  confluence comment list --page 12345 --locations footer
  confluence comment list --page 12345 --json --limit 100

## Usage

```text
confluence comment list [flags]
```

## Options

```text
      --limit int           Max comments (hard cap 200) (default 100)
      --locations strings   Subset of footer,inline,resolved
      --page string         Content id (required)
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

- [confluence comment](confluence_comment.md)
