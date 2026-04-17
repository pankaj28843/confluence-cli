package main

import "github.com/spf13/cobra"

func versionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version, build time, and commit",
		Long: `Print the confluence CLI version, build timestamp, and git commit.

Examples:
  confluence version
  confluence version --json`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			w := getWriter()
			defer w.Finish()
			info := map[string]string{"version": version, "buildTime": buildTime, "commit": commit}
			if w.IsJSON() {
				return w.JSON(info)
			}
			w.Text("confluence %s (commit: %s, built: %s)\n", version, commit, buildTime)
			return nil
		},
	}
}
