# confluence permission

Space permission reads

## Synopsis

Space permission helpers.

Inspect documented space permission assignments across Confluence Cloud and
Server/Data Center. Subject-specific reads are Server/Data Center only.

## Examples

confluence permission space list --space ENG --json
  confluence permission space available --limit 100 --json
  confluence permission space subject --space ENG --group confluence-users --json

## Usage

```text
confluence permission
```

## Commands

- `space` - Space permission reads

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
