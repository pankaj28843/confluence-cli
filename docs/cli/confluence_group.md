# confluence group

Groups (list, view, picker, members, hierarchy)

## Synopsis

Group operations.

## Examples

confluence group list
  confluence group view engineering
  confluence group view --id 11111111-2222-3333-4444-555555555555
  confluence group picker eng
  confluence group members engineering

## Usage

```text
confluence group
```

## Commands

- `ancestors` - List Server/DC ancestor groups
- `children` - List Server/DC child groups
- `list` - List groups
- `members` - List members of a group
- `parents` - List Server/DC parent groups
- `picker` - Search Cloud groups by picker query
- `view` - Show a group

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
