# confluence page ancestors

List ancestor pages of a content id

## Synopsis

Walk parents up to the space root.

## Examples

confluence page ancestors 12345
  confluence page ancestors 12345 --json

## Usage

```text
confluence page ancestors <id>
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
