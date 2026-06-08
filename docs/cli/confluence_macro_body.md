# confluence macro body

Fetch one macro body

## Synopsis

Fetch the body of one macro from a specific content version.

Cloud and Server/Data Center support --macro-id. Server/Data Center also
supports the documented deprecated --hash lookup.

## Examples

confluence macro body --page 12345 --version 2 --macro-id 50884bd9-0cb8-41d5-98be-f80943c14f96
  confluence macro body --page 12345 --version 2 --macro-id my-macro --json
  confluence macro body --page 12345 --version 2 --hash abc123 --json

## Usage

```text
confluence macro body [flags]
```

## Options

```text
      --hash string       Deprecated Server/Data Center macro body hash
      --macro-id string   Macro id
      --page string       Page or content id containing the macro (required)
      --version int       Content version containing the macro (required)
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

- [confluence macro](confluence_macro.md)
