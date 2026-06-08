# confluence task

Tasks (list, view, complete, reopen, long)

## Synopsis

Task operations.

Cloud content tasks are available through list/view/complete/reopen.
Server/Data Center long-running tasks are available under task long.

## Examples

confluence task list --page 12345 --status incomplete
  confluence task view 42 --body-format storage
  confluence task long list --limit 10

## Usage

```text
confluence task
```

## Commands

- `complete` - Mark one Cloud content task complete
- `list` - List Cloud content tasks
- `long` - Server/Data Center long-running tasks
- `reopen` - Mark one Cloud content task incomplete
- `view` - Show one Cloud content task

## Inherited Options

```text
      --debug             Log HTTP requests to stderr (Authorization header redacted)
      --jq string         Filter JSON output through a jq expression (requires jq on PATH)
      --json              Output as JSON (machine-readable)
      --template string   Render JSON output through a Go text/template
      --timing            Show execution time on stderr
```

## See Also

- [confluence](confluence.md)
