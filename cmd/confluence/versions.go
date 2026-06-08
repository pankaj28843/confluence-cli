package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/pankaj28843/confluence-cli/internal/conf"
)

func versionListReadCmd(resource, target string, bodyFormatFlag, commentLocation bool) *cobra.Command {
	var limit int
	var sort string
	var bodyFormat string
	var location string
	extraExample := fmt.Sprintf("  confluence %s versions 12345 --sort -modified-date", resource)
	if bodyFormatFlag {
		extraExample = fmt.Sprintf("  confluence %s versions 12345 --body-format storage --sort -modified-date", resource)
	}
	cmd := &cobra.Command{
		Use:   "versions <id>",
		Short: "List " + resource + " versions",
		Long: fmt.Sprintf(`List version records for one %s.

Cloud uses the documented v2 Version endpoints. Server/Data Center uses the
documented content version route where the target is content-like.

Examples:
  confluence %s versions 12345
  confluence %s versions 12345 --limit 25 --json
%s`, resource, resource, resource, extraExample),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := newContext()
			defer cancel()
			w := getWriter()
			defer w.Finish()
			c, err := newClient()
			if err != nil {
				return err
			}
			selectedTarget := target
			if commentLocation {
				selectedTarget, err = versionCommentTarget(location)
				if err != nil {
					return err
				}
			}
			versions, err := conf.ListVersions(ctx, c, selectedTarget, args[0], conf.VersionListOptions{
				Limit:      limit,
				Sort:       sort,
				BodyFormat: bodyFormat,
			})
			if err != nil {
				return err
			}
			if w.IsJSON() {
				return w.JSON(versions)
			}
			printVersions(w, versions)
			return nil
		},
	}
	cmd.Flags().IntVar(&limit, "limit", 25, "Max versions (hard cap 200)")
	cmd.Flags().StringVar(&sort, "sort", "", "Cloud sort expression")
	if bodyFormatFlag {
		cmd.Flags().StringVar(&bodyFormat, "body-format", "", "Cloud body representation to include, e.g. storage")
	}
	if commentLocation {
		cmd.Flags().StringVar(&location, "location", "footer", "Cloud comment location: footer or inline")
	}
	return cmd
}

func versionDetailReadCmd(resource, target string, commentLocation bool) *cobra.Command {
	var location string
	cmd := &cobra.Command{
		Use:   "version <id> <number>",
		Short: "Show one " + resource + " version",
		Long: fmt.Sprintf(`Show one version record for one %s.

Cloud uses the documented v2 Version detail endpoints. Server/Data Center uses
the documented content version detail route where the target is content-like.

Examples:
  confluence %s version 12345 2
  confluence %s version 12345 2 --json`, resource, resource, resource),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := newContext()
			defer cancel()
			w := getWriter()
			defer w.Finish()
			versionNumber, err := strconv.Atoi(args[1])
			if err != nil || versionNumber <= 0 {
				return fmt.Errorf("invalid version number %q", args[1])
			}
			c, err := newClient()
			if err != nil {
				return err
			}
			selectedTarget := target
			if commentLocation {
				selectedTarget, err = versionCommentTarget(location)
				if err != nil {
					return err
				}
			}
			version, err := conf.GetVersion(ctx, c, selectedTarget, args[0], versionNumber)
			if err != nil {
				return err
			}
			if w.IsJSON() {
				return w.JSON(version)
			}
			printVersions(w, []conf.Version{*version})
			return nil
		},
	}
	if commentLocation {
		cmd.Flags().StringVar(&location, "location", "footer", "Cloud comment location: footer or inline")
	}
	return cmd
}

func versionCommentTarget(location string) (string, error) {
	switch strings.ToLower(strings.TrimSpace(location)) {
	case "", "footer":
		return "footer-comment", nil
	case "inline":
		return "inline-comment", nil
	default:
		return "", fmt.Errorf("unsupported comment version location %q", location)
	}
}

type versionTextWriter interface {
	Text(format string, args ...any)
}

func printVersions(w versionTextWriter, versions []conf.Version) {
	for _, version := range versions {
		when := firstNonEmpty(version.CreatedAt, version.When, "-")
		author := firstNonEmpty(version.AuthorID, version.By.DisplayName, version.By.PublicName, version.By.Username, version.By.AccountID, "-")
		w.Text("%d\t%s\t%s", version.Number, when, author)
		if version.Message != "" {
			w.Text("\t%s", version.Message)
		}
		w.Text("\n")
	}
}
