# confluence permission space list

List permissions assigned on a space

## Synopsis

List permissions assigned on a space.

On Cloud this uses the v2 space permission assignment endpoint. A space key is
resolved to its Cloud space id before reading permissions. On Server/Data Center
this uses the space permissions endpoint.

## Examples

confluence permission space list --space ENG
  confluence permission space list --space ENG --limit 100 --json

## Usage

```text
confluence permission space list [flags]
```

## Options

```text
      --limit int      Max permission assignments (hard cap 200) (default 25)
      --space string   Space key or Cloud space id (required)
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

- [confluence permission space](confluence_permission_space.md)
