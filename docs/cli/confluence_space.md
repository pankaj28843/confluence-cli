# confluence space

Spaces (list, view)

## Synopsis

Space operations.

## Examples

confluence space list
  confluence space list --type global --status current --json
  confluence space view ENG

## Usage

```text
confluence space
```

## Commands

- `list` - List spaces in the site
- `view` - Show one space

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
