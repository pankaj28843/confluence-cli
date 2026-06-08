package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/pankaj28843/confluence-cli/internal/conf"
)

func macroCmdReal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "macro",
		Short: "Macro body utilities",
		Long: `Macro body utility operations.

Examples:
  confluence macro body --page 12345 --version 2 --macro-id 50884bd9-0cb8-41d5-98be-f80943c14f96
  confluence macro body --page 12345 --version 2 --hash abc123 --json`,
	}
	cmd.AddCommand(macroBodyCmd())
	return cmd
}

func macroBodyCmd() *cobra.Command {
	var pageID, macroID, hash string
	var version int
	cmd := &cobra.Command{
		Use:   "body",
		Short: "Fetch one macro body",
		Long: `Fetch the body of one macro from a specific content version.

Cloud and Server/Data Center support --macro-id. Server/Data Center also
supports the documented deprecated --hash lookup.

Examples:
  confluence macro body --page 12345 --version 2 --macro-id 50884bd9-0cb8-41d5-98be-f80943c14f96
  confluence macro body --page 12345 --version 2 --macro-id my-macro --json
  confluence macro body --page 12345 --version 2 --hash abc123 --json`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := newContext()
			defer cancel()
			w := getWriter()
			defer w.Finish()
			if pageID == "" {
				return fmt.Errorf("--page is required")
			}
			if version <= 0 {
				return fmt.Errorf("--version must be greater than zero")
			}
			if (macroID == "") == (hash == "") {
				return fmt.Errorf("exactly one of --macro-id or --hash is required")
			}
			c, err := newClient()
			if err != nil {
				return err
			}
			got, err := conf.GetMacroBody(ctx, c, conf.MacroLookup{
				ContentID: pageID,
				Version:   version,
				MacroID:   macroID,
				Hash:      hash,
			})
			if err != nil {
				return err
			}
			if w.IsJSON() {
				return w.JSON(got)
			}
			printMacroInstance(w, *got)
			return nil
		},
	}
	cmd.Flags().StringVar(&pageID, "page", "", "Page or content id containing the macro (required)")
	cmd.Flags().IntVar(&version, "version", 0, "Content version containing the macro (required)")
	cmd.Flags().StringVar(&macroID, "macro-id", "", "Macro id")
	cmd.Flags().StringVar(&hash, "hash", "", "Deprecated Server/Data Center macro body hash")
	return cmd
}

type macroTextWriter interface {
	Text(format string, args ...interface{})
}

func printMacroInstance(w macroTextWriter, got conf.MacroInstance) {
	if got.Body != "" {
		w.Text("%s\n", got.Body)
		return
	}
	if got.Name != "" {
		w.Text("%s\n", got.Name)
	}
}
