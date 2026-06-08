# confluence comment

Comments (list/versions/version/add/update/delete)

## Synopsis

Comment operations.

## Examples

confluence comment list --page 12345
  confluence comment add --page 12345 --body "<p>Looks good.</p>"
  confluence comment update 998877 --body-file comment.html
  confluence comment delete 998877 --force
  confluence comment list --page 12345 --locations footer,inline
  confluence comment list --page 12345 --json --limit 50

## Usage

```text
confluence comment
```

## Commands

- `add` - Add a footer comment
- `delete` - Delete a footer comment
- `list` - List comments (footer, inline, resolved by default)
- `update` - Update a footer comment body
- `version` - Show one comment version
- `versions` - List comment versions

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
