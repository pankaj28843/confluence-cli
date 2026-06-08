# confluence property space

Space properties (list, get, set, delete)

## Synopsis

Space property operations for a space key.

## Examples

confluence property space list --space ENG
  confluence property space set --space ENG --key retention --value '{"days":30}'

## Usage

```text
confluence property space
```

## Commands

- `delete` - Delete one space property by key
- `get` - Get one space property by key
- `list` - List space properties
- `set` - Create or update one space property

## Inherited Options

```text
      --debug             Log HTTP requests to stderr (Authorization header redacted)
      --jq string         Filter JSON output through a jq expression (requires jq on PATH)
      --json              Output as JSON (machine-readable)
      --template string   Render JSON output through a Go text/template
      --timing            Show execution time on stderr
```

## See Also

- [confluence property](confluence_property.md)
