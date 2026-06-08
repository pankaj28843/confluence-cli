# confluence content-state current

Show current Cloud content state

## Synopsis

Show the content state attached to the draft, current, or archived
version of a Cloud content item.

## Examples

confluence content-state current 12345
  confluence content-state current 12345 --status draft --json

## Usage

```text
confluence content-state current <content-id> [flags]
```

## Options

```text
      --status string   Cloud content status: current, draft, or archived (default "current")
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
