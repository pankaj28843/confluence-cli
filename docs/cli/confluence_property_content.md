# confluence property content

Content properties (list, get, set, delete)

## Synopsis

Content property operations for a page/content id.

## Examples

confluence property content list --page 12345
  confluence property content set --page 12345 --key release --value '{"ready":true}'

## Usage

```text
confluence property content
```

## Commands

- `delete` - Delete one content property by key
- `get` - Get one content property by key
- `list` - List content properties
- `set` - Create or update one content property

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
