# confluence task view

Show one Cloud content task

## Synopsis

Show one Confluence Cloud content task by id.

## Examples

confluence task view 42
  confluence task view 42 --body-format storage --json

## Usage

```text
confluence task view <id> [flags]
```

## Options

```text
      --body-format string   Body format: storage | atlas_doc_format
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
