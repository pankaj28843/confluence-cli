# confluence task long list

List Server/Data Center long-running tasks

## Synopsis

List Server/Data Center long-running tasks.

## Examples

confluence task long list
  confluence task long list --expand messages --limit 10 --json

## Usage

```text
confluence task long list [flags]
```

## Options

```text
      --expand string   Expand parameter, for example messages
      --limit int       Max long tasks (hard cap 200) (default 25)
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
