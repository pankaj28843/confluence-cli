# confluence task list

List Cloud content tasks

## Synopsis

List Confluence Cloud content tasks. Server/Data Center long-running tasks
are under task long list.

## Examples

confluence task list --status incomplete --limit 50
  confluence task list --page 12345 --body-format storage --json
  confluence task list --assigned-to 557058:abc --include-blank

## Usage

```text
confluence task list [flags]
```

## Options

```text
      --assigned-to strings    Assignee account id filter; repeatable or comma-separated
      --blogpost string        Blog post id filter
      --body-format string     Body format: storage | atlas_doc_format
      --completed-by strings   Completer account id filter; repeatable or comma-separated
      --created-by strings     Creator account id filter; repeatable or comma-separated
      --include-blank          Include blank tasks
      --limit int              Max tasks (hard cap 200) (default 25)
      --page string            Page id filter
      --space-id strings       Cloud space id filter; repeatable or comma-separated
      --status string          Filter: complete | incomplete
      --task-id strings        Task id filter; repeatable or comma-separated
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

- [confluence task](confluence_task.md)
