# confluence group view

Show a group

## Synopsis

Show a group. On Server/DC pass the group name argument. On Cloud pass --id.

## Examples

confluence group view engineering
  confluence group view engineering --expand members
  confluence group view --id 11111111-2222-3333-4444-555555555555 --json

## Usage

```text
confluence group view [name] [flags]
```

## Options

```text
      --expand string   Server/DC expand value
      --id string       Cloud group id
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

- [confluence group](confluence_group.md)
