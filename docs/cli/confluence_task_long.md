# confluence task long

Server/Data Center long-running tasks

## Synopsis

Server/Data Center long-running task operations.

## Examples

confluence task long list --limit 10
  confluence task long view 123456 --expand messages

## Usage

```text
confluence task long
```

## Commands

- `list` - List Server/Data Center long-running tasks
- `view` - Show one Server/Data Center long-running task

## Inherited Options

```text
      --debug             Log HTTP requests to stderr (Authorization header redacted)
      --jq string         Filter JSON output through a jq expression (requires jq on PATH)
      --json              Output as JSON (machine-readable)
      --template string   Render JSON output through a Go text/template
      --timing            Show execution time on stderr
```

## See Also

- [confluence task](confluence_task.md)
