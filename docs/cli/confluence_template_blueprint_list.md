# confluence template blueprint list

List Cloud blueprint templates

## Synopsis

List Confluence Cloud blueprint templates. Use --space to list blueprints in
one space, or omit it for global blueprints.

## Examples

confluence template blueprint list
  confluence template blueprint list --space ENG --limit 10 --json
  confluence template blueprint list --expand body.storage

## Usage

```text
confluence template blueprint list [flags]
```

## Options

```text
      --expand strings   Expand value; repeatable or comma-separated
      --limit int        Max templates (hard cap 200) (default 25)
      --space string     Space key; omit for global blueprints
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

- [confluence template blueprint](confluence_template_blueprint.md)
