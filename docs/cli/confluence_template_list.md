# confluence template list

List Cloud content templates

## Synopsis

List Confluence Cloud content templates. Use --space to list templates in
one space, or omit it for global templates.

## Examples

confluence template list
  confluence template list --space ENG --limit 10 --json
  confluence template list --expand body.storage --expand space

## Usage

```text
confluence template list [flags]
```

## Options

```text
      --expand strings   Expand value; repeatable or comma-separated
      --limit int        Max templates (hard cap 200) (default 25)
      --space string     Space key; omit for global templates
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

- [confluence template](confluence_template.md)
