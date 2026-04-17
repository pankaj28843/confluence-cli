package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/spf13/cobra"
)

var (
	jsonOutput bool
	jqExpr     string
	tmpl       string
	timing     bool
	debug      bool

	version   = "dev"
	buildTime = "unknown"
	commit    = "unknown"
)

const exitCodeConfig = 2

type configError struct{ error }

func newConfigError(err error) error {
	if err == nil {
		return nil
	}
	return configError{err}
}

func main() {
	root := &cobra.Command{
		Use:   "confluence",
		Short: "Atlassian Confluence CLI (Server/DC + Cloud)",
		Long: fmt.Sprintf(`confluence - Atlassian Confluence CLI
Version: %s (built %s, commit %s)

Query and act on Atlassian Confluence from the terminal. Supports both
Confluence Server / Data Center (Bearer PAT, /rest/api/...) and Confluence
Cloud (Basic email:token, /wiki/rest/api/ + /wiki/api/v2/).

Requires environment variables:
  CONFLUENCE_URL                Base URL, e.g. https://wiki.example.com
                                (Server/DC) or https://example.atlassian.net/wiki (Cloud)
  Server/DC:
    CONFLUENCE_PAT              Personal access token (alias:
                                CONFLUENCE_PERSONAL_ACCESS_TOKEN)
  Cloud:
    CONFLUENCE_EMAIL            Atlassian account email
    CONFLUENCE_API_TOKEN        API token from id.atlassian.com
  CONFLUENCE_FLAVOR             Optional override: server|cloud
  CONFLUENCE_DEFAULT_SPACE      Optional default for --space

Workflow:
  confluence doctor                                         Health + auth + flavor probe
  confluence space list                                     Spaces in the site
  confluence page search --cql "type=page AND space=ENG"    CQL-powered page search
  confluence page view 12345 --markdown                     Fetch + render a page
  confluence page children 12345 --recursive --json         Page-tree traversal
  confluence page comments 12345                            Inline + footer + resolved
  confluence label add 12345 --label needs-review           Low-risk mutation
  confluence search all "release process" --json            Unified fan-out (code/space/user/attachment)
  confluence api /rest/api/user/current                     Raw REST passthrough`, version, buildTime, commit),
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	root.PersistentFlags().BoolVar(&jsonOutput, "json", false, "Output as JSON (machine-readable)")
	root.PersistentFlags().StringVar(&jqExpr, "jq", "", "Filter JSON output through a jq expression (requires jq on PATH)")
	root.PersistentFlags().StringVar(&tmpl, "template", "", "Render JSON output through a Go text/template")
	root.PersistentFlags().BoolVar(&timing, "timing", false, "Show execution time on stderr")
	root.PersistentFlags().BoolVar(&debug, "debug", false, "Log HTTP requests to stderr (Authorization header redacted)")
	root.Version = fmt.Sprintf("%s (built %s, commit %s)", version, buildTime, commit)

	root.AddCommand(versionCmd())
	root.AddCommand(doctorCmd())
	root.AddCommand(spaceCmd())
	root.AddCommand(pageCmd())
	root.AddCommand(attachmentCmd())
	root.AddCommand(labelCmd())
	root.AddCommand(commentCmd())
	root.AddCommand(userCmd())
	root.AddCommand(groupCmd())
	root.AddCommand(watcherCmd())
	root.AddCommand(restrictionCmd())
	root.AddCommand(searchCmd())
	root.AddCommand(apiCmd())

	if err := root.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		if _, ok := err.(configError); ok {
			os.Exit(exitCodeConfig)
		}
		os.Exit(1)
	}
}

func newContext() (context.Context, context.CancelFunc) {
	return signal.NotifyContext(context.Background(), os.Interrupt)
}
