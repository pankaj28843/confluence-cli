# confluence comment delete

Delete a footer comment

## Synopsis

Delete a footer comment permanently. A confirmation prompt is shown unless
--force is supplied.

## Examples

confluence comment delete 998877
  confluence comment delete 998877 --force
  confluence comment delete 998877 --force --json

## Usage

```text
confluence comment delete <id> [flags]
```

## Options

```text
      --force   Do not prompt for confirmation
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
