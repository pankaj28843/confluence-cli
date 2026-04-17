package main

import "github.com/pankaj28843/confluence-cli/internal/output"

func getWriter() *output.Writer {
	return output.New(output.Options{
		JSON:     jsonOutput,
		JQ:       jqExpr,
		Template: tmpl,
		Timing:   timing,
	})
}
