# confluence custom-content

Cloud custom content reads

## Synopsis

Cloud custom-content read operations.

These commands use documented Cloud v2 custom-content, children, and version
endpoints. Create, update, and delete are mutations and remain deferred behind
explicit safety gates.

## Examples

confluence custom-content list --type ac:example --json
  confluence custom-content page 12345 --type ac:example
  confluence custom-content view 67890 --include-version --json

## Usage

```text
confluence custom-content
```

## Commands

- `blogpost` - List Cloud custom content in a blogpost
- `children` - List child Cloud custom content
- `list` - List Cloud custom content by type
- `page` - List Cloud custom content in a page
- `space` - List Cloud custom content in a space
- `version` - Show one custom-content version
- `versions` - List custom-content versions
- `view` - Show one Cloud custom-content record

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
