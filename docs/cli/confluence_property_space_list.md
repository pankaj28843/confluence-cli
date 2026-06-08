# confluence property space list

List space properties

## Synopsis

List properties for a space key.

## Examples

confluence property space list --space ENG
  confluence property space list --space ENG --key retention --json
  confluence property space list --space ENG --limit 100

## Usage

```text
confluence property space list [flags]
```

## Options

```text
      --key string     Optional property key filter
      --limit int      Max properties (hard cap 200) (default 25)
      --space string   Space key (required)
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

- [confluence property space](confluence_property_space.md)
