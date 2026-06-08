# confluence settings lookandfeel

Show Cloud look-and-feel settings

## Synopsis

Show global or space-specific Confluence Cloud look-and-feel settings.

## Examples

confluence settings lookandfeel
  confluence settings lookandfeel --space ENG --json

## Usage

```text
confluence settings lookandfeel [flags]
```

## Options

```text
      --space string   Space key; omit for global look-and-feel settings
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

- [confluence settings](confluence_settings.md)
