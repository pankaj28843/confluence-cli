# confluence body convert

Convert a content body representation

## Synopsis

Convert a Confluence content body between documented representations.

Cloud uses the asynchronous content-body conversion API and polls for the
result by default. Use --poll-attempts 0 to return only the Cloud async id.

## Examples

confluence body convert --from storage --to view --value '<p>Hello</p>'
  confluence body convert --from storage --to export_view --value @body.xml --json
  confluence body convert --to view --value @- --expand webresource.uris.css --expand webresource.uris.js

## Usage

```text
confluence body convert [flags]
```

## Options

```text
      --content-context string   Cloud content id context for permission-sensitive conversion
      --embedded-render string   Cloud embeddedContentRender value, for example current
      --expand strings           Expand value; repeatable or comma-separated
      --from string              Source representation (default "storage")
      --no-cache                 Cloud only: queue conversion with allowCache=false
      --poll-attempts int        Cloud poll attempts after queueing; 0 returns async id only (default 10)
      --poll-interval duration   Cloud poll interval (default 500ms)
      --space-context string     Cloud space key context for permission-sensitive conversion
      --to string                Target representation (default "view")
      --value string             Body value, @file, or @- (required)
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

- [confluence body](confluence_body.md)
