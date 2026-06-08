package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/pankaj28843/confluence-cli/internal/conf"
)

func customContentCmdReal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "custom-content",
		Short: "Cloud custom content reads",
		Long: `Cloud custom-content read operations.

These commands use documented Cloud v2 custom-content, children, and version
endpoints. Create, update, and delete are mutations and remain deferred behind
explicit safety gates.

Examples:
  confluence custom-content list --type ac:example --json
  confluence custom-content page 12345 --type ac:example
  confluence custom-content view 67890 --include-version --json`,
	}
	cmd.AddCommand(customContentListCmd())
	cmd.AddCommand(customContentContainerListCmd("page", "page", true))
	cmd.AddCommand(customContentContainerListCmd("blogpost", "blogpost", true))
	cmd.AddCommand(customContentContainerListCmd("space", "space", false))
	cmd.AddCommand(customContentViewCmd())
	cmd.AddCommand(customContentChildrenCmd())
	cmd.AddCommand(customContentVersionsCmd())
	cmd.AddCommand(customContentVersionCmd())
	return cmd
}

func customContentVersionsCmd() *cobra.Command {
	return versionListReadCmd("custom-content", "custom-content", true, false)
}

func customContentVersionCmd() *cobra.Command {
	return versionDetailReadCmd("custom-content", "custom-content", false)
}

func customContentListCmd() *cobra.Command {
	var typ string
	var ids []string
	var spaceIDs []string
	var limit int
	var sort string
	var bodyFormat string
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List Cloud custom content by type",
		Long: `List Cloud custom content by type using the documented v2 global endpoint.

Examples:
  confluence custom-content list --type ac:example
  confluence custom-content list --type ac:example --space-id 100 --limit 25 --json
  confluence custom-content list --type ac:example --id 777 --body-format storage`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := newContext()
			defer cancel()
			w := getWriter()
			defer w.Finish()
			if typ == "" {
				return fmt.Errorf("--type is required")
			}
			c, err := newClient()
			if err != nil {
				return err
			}
			items, err := conf.ListCustomContent(ctx, c, conf.CustomContentListOptions{
				Type:       typ,
				IDs:        ids,
				SpaceIDs:   spaceIDs,
				Limit:      limit,
				Sort:       sort,
				BodyFormat: bodyFormat,
			})
			if err != nil {
				return err
			}
			if w.IsJSON() {
				return w.JSON(items)
			}
			printCustomContent(w, items)
			return nil
		},
	}
	cmd.Flags().StringVar(&typ, "type", "", "Required custom content type")
	cmd.Flags().StringSliceVar(&ids, "id", nil, "Custom content id filter; repeatable or comma-separated")
	cmd.Flags().StringSliceVar(&spaceIDs, "space-id", nil, "Space id filter; repeatable or comma-separated")
	cmd.Flags().IntVar(&limit, "limit", 25, "Max results (hard cap 200)")
	cmd.Flags().StringVar(&sort, "sort", "", "Cloud sort expression")
	cmd.Flags().StringVar(&bodyFormat, "body-format", "", "Cloud body representation to include, e.g. storage")
	return cmd
}

func customContentContainerListCmd(use, containerType string, includeSort bool) *cobra.Command {
	var typ string
	var limit int
	var sort string
	var bodyFormat string
	cmd := &cobra.Command{
		Use:   use + " <id>",
		Short: "List Cloud custom content in a " + containerType,
		Long: fmt.Sprintf(`List Cloud custom content by type inside one %s.

Examples:
  confluence custom-content %s 12345 --type ac:example
  confluence custom-content %s 12345 --type ac:example --body-format storage --json`, containerType, use, use),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := newContext()
			defer cancel()
			w := getWriter()
			defer w.Finish()
			if typ == "" {
				return fmt.Errorf("--type is required")
			}
			c, err := newClient()
			if err != nil {
				return err
			}
			items, err := conf.ListCustomContent(ctx, c, conf.CustomContentListOptions{
				Type:          typ,
				ContainerType: containerType,
				ContainerID:   args[0],
				Limit:         limit,
				Sort:          sort,
				BodyFormat:    bodyFormat,
			})
			if err != nil {
				return err
			}
			if w.IsJSON() {
				return w.JSON(items)
			}
			printCustomContent(w, items)
			return nil
		},
	}
	cmd.Flags().StringVar(&typ, "type", "", "Required custom content type")
	cmd.Flags().IntVar(&limit, "limit", 25, "Max results (hard cap 200)")
	if includeSort {
		cmd.Flags().StringVar(&sort, "sort", "", "Cloud sort expression")
	}
	cmd.Flags().StringVar(&bodyFormat, "body-format", "", "Cloud body representation to include, e.g. storage")
	return cmd
}

func customContentViewCmd() *cobra.Command {
	var bodyFormat string
	var version int
	var includeLabels bool
	var includeProperties bool
	var includeOperations bool
	var includeVersions bool
	var includeVersion bool
	var includeCollaborators bool
	cmd := &cobra.Command{
		Use:   "view <id>",
		Short: "Show one Cloud custom-content record",
		Long: `Show one Cloud custom-content record by id.

Examples:
  confluence custom-content view 12345
  confluence custom-content view 12345 --body-format storage --include-version --json
  confluence custom-content view 12345 --include-labels --include-properties`,
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
			item, err := conf.GetCustomContent(ctx, c, args[0], conf.CustomContentGetOptions{
				BodyFormat:           bodyFormat,
				Version:              version,
				IncludeLabels:        includeLabels,
				IncludeProperties:    includeProperties,
				IncludeOperations:    includeOperations,
				IncludeVersions:      includeVersions,
				IncludeVersion:       includeVersion,
				IncludeCollaborators: includeCollaborators,
			})
			if err != nil {
				return err
			}
			if w.IsJSON() {
				return w.JSON(item)
			}
			printCustomContent(w, []conf.Content{*item})
			return nil
		},
	}
	cmd.Flags().StringVar(&bodyFormat, "body-format", "", "Cloud body representation to include, e.g. storage")
	cmd.Flags().IntVar(&version, "version", 0, "Cloud content version to retrieve")
	cmd.Flags().BoolVar(&includeLabels, "include-labels", false, "Include labels in the Cloud response")
	cmd.Flags().BoolVar(&includeProperties, "include-properties", false, "Include properties in the Cloud response")
	cmd.Flags().BoolVar(&includeOperations, "include-operations", false, "Include permitted operations in the Cloud response")
	cmd.Flags().BoolVar(&includeVersions, "include-versions", false, "Include versions in the Cloud response")
	cmd.Flags().BoolVar(&includeVersion, "include-version", false, "Include current version in the Cloud response")
	cmd.Flags().BoolVar(&includeCollaborators, "include-collaborators", false, "Include collaborators in the Cloud response")
	return cmd
}

func customContentChildrenCmd() *cobra.Command {
	var limit int
	var sort string
	var types []string
	cmd := &cobra.Command{
		Use:   "children <id>",
		Short: "List child Cloud custom content",
		Long: `List child custom content under one Cloud custom-content id.

Examples:
  confluence custom-content children 12345
  confluence custom-content children 12345 --type ac:example --json
  confluence custom-content children 12345 --sort title --limit 25`,
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
			children, err := conf.ListCustomContentChildren(ctx, c, args[0], conf.DirectChildrenOptions{
				Limit: limit,
				Types: types,
				Sort:  sort,
			})
			if err != nil {
				return err
			}
			if w.IsJSON() {
				return w.JSON(children)
			}
			printCustomContent(w, children)
			return nil
		},
	}
	cmd.Flags().IntVar(&limit, "limit", 25, "Max children (hard cap 200)")
	cmd.Flags().StringVar(&sort, "sort", "", "Cloud sort expression")
	cmd.Flags().StringSliceVar(&types, "type", nil, "Client-side custom content type filter; repeatable or comma-separated")
	return cmd
}

type customContentTextWriter interface {
	Text(format string, args ...any)
}

func printCustomContent(w customContentTextWriter, items []conf.Content) {
	for _, item := range items {
		w.Text("%s\t%s\t%s\t%s", item.ID, firstNonEmpty(item.Type, "-"), firstNonEmpty(item.Status, "-"), firstNonEmpty(item.Title, "-"))
		if item.SpaceID != "" {
			w.Text("\tspace=%s", item.SpaceID)
		}
		if item.PageID != "" {
			w.Text("\tpage=%s", item.PageID)
		}
		if item.BlogPostID != "" {
			w.Text("\tblogpost=%s", item.BlogPostID)
		}
		if item.CustomContentID != "" {
			w.Text("\tparent=%s", item.CustomContentID)
		}
		w.Text("\n")
	}
}
