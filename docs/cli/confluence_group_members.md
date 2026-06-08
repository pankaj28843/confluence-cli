# confluence group members

List members of a group

## Synopsis

List group members.

## Examples

confluence group members engineering
  confluence group members engineering --json --limit 200

## Usage

```text
confluence group members <name> [flags]
```

## Options

```text
      --limit int   Max members (hard cap 200) (default 100)
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

- [confluence group](confluence_group.md)
