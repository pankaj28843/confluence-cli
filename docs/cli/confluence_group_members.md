# confluence group members

List members of a group

## Synopsis

List group members. On Server/DC pass the group name argument. On Cloud pass --id.

## Examples

confluence group members engineering
  confluence group members engineering --json --limit 200
  confluence group members --id 11111111-2222-3333-4444-555555555555 --expand personalSpace --json

## Usage

```text
confluence group members [name] [flags]
```

## Options

```text
      --expand strings   Expand value; repeatable or comma-separated
      --id string        Cloud group id
      --limit int        Max members (hard cap 200) (default 25)
      --total-size       Ask Cloud to include total size metadata
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
