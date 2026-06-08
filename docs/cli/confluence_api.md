# confluence api

Call any Confluence REST endpoint (escape hatch)

## Synopsis

Issue a raw REST call to Confluence. Auto-routes by path — start your
path with /rest/api/... for v1 or /api/v2/... for Cloud v2. The path is
appended to CONFLUENCE_URL verbatim. For Cloud, CONFLUENCE_URL already ends in
/wiki, so do not add another /wiki prefix in the path.

## Body

--data '@file.json'    reads from disk
  --data '@-'            reads from stdin
  --data '<literal>'     inline JSON

Auth (Bearer for Server/DC, Basic for Cloud) is attached automatically and
never logged.

## Examples

confluence api /rest/api/user/current
  confluence api /rest/api/space --param 'limit=10' --jq '.results[].key'
  confluence api /rest/api/content/search --param 'cql=type=page' --param 'limit=5'
  confluence api /api/v2/spaces --param 'limit=5'              # Cloud v2
  echo '{"foo":"bar"}' | confluence api /rest/api/some-write -X POST --data '@-'

## Usage

```text
confluence api <path> [flags]
```

## Options

```text
      --data string          JSON body; prefix with @ for file, @- for stdin
  -H, --header stringArray   Extra header (reserved)
  -X, --method string        HTTP method (GET, POST, PUT, DELETE) (default "GET")
      --param stringArray    Query parameter k=v (may be repeated)
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

- [confluence](confluence.md)
