package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func apiCmdReal() *cobra.Command {
	var method string
	var params []string
	var dataFlag string
	var headers []string

	cmd := &cobra.Command{
		Use:   "api <path>",
		Short: "Call any Confluence REST endpoint (escape hatch)",
		Long: `Issue a raw REST call to Confluence. Auto-routes by path — start your
path with /rest/api/... for v1 or /wiki/api/v2/... for Cloud v2. The path is
appended to CONFLUENCE_URL verbatim (the client does NOT inject /wiki).

Body:
  --data '@file.json'    reads from disk
  --data '@-'            reads from stdin
  --data '<literal>'     inline JSON

Auth (Bearer for Server/DC, Basic for Cloud) is attached automatically and
never logged.

Examples:
  confluence api /rest/api/user/current
  confluence api /rest/api/space --param 'limit=10' --jq '.results[].key'
  confluence api /rest/api/content/search --param 'cql=type=page' --param 'limit=5'
  confluence api /wiki/api/v2/spaces --param 'limit=5'         # Cloud v2
  echo '{"foo":"bar"}' | confluence api /rest/api/some-write -X POST --data '@-'`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := newContext()
			defer cancel()
			w := getWriter()
			defer w.Finish()

			path := args[0]
			if !strings.HasPrefix(path, "/") {
				path = "/" + path
			}

			qs, err := parseParams(params)
			if err != nil {
				return err
			}
			c, err := newClient()
			if err != nil {
				return err
			}

			var body any
			if dataFlag != "" {
				raw, err := readDataArg(dataFlag)
				if err != nil {
					return err
				}
				if len(raw) > 0 {
					var parsed any
					if err := json.Unmarshal(raw, &parsed); err != nil {
						return fmt.Errorf("--data is not valid JSON: %w", err)
					}
					body = parsed
				}
			}

			_ = headers // reserved for future per-call header overrides.

			var (
				data []byte
				rerr error
			)
			switch strings.ToUpper(method) {
			case "GET":
				data, _, rerr = c.Get(ctx, path, qs)
			case "POST":
				data, _, rerr = c.Post(ctx, path, qs, body)
			case "PUT":
				data, _, rerr = c.Put(ctx, path, qs, body)
			case "DELETE":
				data, _, rerr = c.Delete(ctx, path, qs)
			default:
				return fmt.Errorf("unsupported --method %q (GET, POST, PUT, DELETE supported)", method)
			}
			if rerr != nil {
				return rerr
			}

			if w.IsJSON() {
				var parsed any
				if err := json.Unmarshal(data, &parsed); err != nil {
					return w.JSON(map[string]string{"raw": string(data)})
				}
				return w.JSON(parsed)
			}
			_, _ = w.Out.Write(data)
			if len(data) > 0 && data[len(data)-1] != '\n' {
				_, _ = io.WriteString(w.Out, "\n")
			}
			return nil
		},
	}
	cmd.Flags().StringVarP(&method, "method", "X", "GET", "HTTP method (GET, POST, PUT, DELETE)")
	cmd.Flags().StringArrayVar(&params, "param", nil, "Query parameter k=v (may be repeated)")
	cmd.Flags().StringVar(&dataFlag, "data", "", "JSON body; prefix with @ for file, @- for stdin")
	cmd.Flags().StringArrayVarP(&headers, "header", "H", nil, "Extra header (reserved)")
	return cmd
}

func parseParams(pairs []string) (url.Values, error) {
	v := url.Values{}
	for _, p := range pairs {
		idx := strings.Index(p, "=")
		if idx <= 0 {
			return nil, fmt.Errorf("--param must be k=v, got %q", p)
		}
		v.Add(p[:idx], p[idx+1:])
	}
	return v, nil
}

func readDataArg(s string) ([]byte, error) {
	if !strings.HasPrefix(s, "@") {
		return []byte(s), nil
	}
	if s == "@-" {
		return io.ReadAll(os.Stdin)
	}
	return os.ReadFile(s[1:])
}
