# confluence task long view

Show one Server/Data Center long-running task

## Synopsis

Show one Server/Data Center long-running task by id.

## Examples

confluence task long view 123456
  confluence task long view 123456 --expand messages --json

## Usage

```text
confluence task long view <id> [flags]
```

## Options

```text
      --expand string   Expand parameter, for example messages
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

- [confluence task long](confluence_task_long.md)
