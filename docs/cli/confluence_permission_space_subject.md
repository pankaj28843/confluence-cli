# confluence permission space subject

List Server/Data Center space permissions for one subject

## Synopsis

List Server/Data Center space permissions for one subject.

Select exactly one subject with --anonymous, --group, or --user-key.

## Examples

confluence permission space subject --space ENG --anonymous
  confluence permission space subject --space ENG --group confluence-users --json
  confluence permission space subject --space ENG --user-key ada --limit 100

## Usage

```text
confluence permission space subject [flags]
```

## Options

```text
      --anonymous         Read anonymous subject permissions
      --group string      Read permissions for this group name
      --limit int         Max permission assignments (hard cap 200) (default 25)
      --space string      Space key (required)
      --user-key string   Read permissions for this Server/Data Center user key
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
