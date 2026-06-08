# confluence permission space

Space permission reads

## Synopsis

Space permission reads.

## Examples

confluence permission space list --space ENG
  confluence permission space available --json
  confluence permission space subject --space ENG --anonymous

## Usage

```text
confluence permission space
```

## Commands

- `available` - List available Cloud space permissions
- `list` - List permissions assigned on a space
- `subject` - List Server/Data Center space permissions for one subject

## Inherited Options

```text
      --debug             Log HTTP requests to stderr (Authorization header redacted)
      --jq string         Filter JSON output through a jq expression (requires jq on PATH)
      --json              Output as JSON (machine-readable)
      --template string   Render JSON output through a Go text/template
      --timing            Show execution time on stderr
```

## See Also

- [confluence permission](confluence_permission.md)
