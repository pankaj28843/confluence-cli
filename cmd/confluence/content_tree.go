package main

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/pankaj28843/confluence-cli/internal/conf"
)

type contentTreeCommandSpec struct {
	Use      string
	Kind     string
	Name     string
	Endpoint string
}

func databaseCmdReal() *cobra.Command {
	return contentTreeCmd(contentTreeCommandSpec{
		Use:      "database",
		Kind:     "database",
		Name:     "Database",
		Endpoint: "databases",
	})
}

func folderCmdReal() *cobra.Command {
	return contentTreeCmd(contentTreeCommandSpec{
		Use:      "folder",
		Kind:     "folder",
		Name:     "Folder",
		Endpoint: "folders",
	})
}

func whiteboardCmdReal() *cobra.Command {
	return contentTreeCmd(contentTreeCommandSpec{
		Use:      "whiteboard",
		Kind:     "whiteboard",
		Name:     "Whiteboard",
		Endpoint: "whiteboards",
	})
}

func smartLinkCmdReal() *cobra.Command {
	return contentTreeCmd(contentTreeCommandSpec{
		Use:      "smart-link",
		Kind:     "smart-link",
		Name:     "Smart Link",
		Endpoint: "embeds",
	})
}

func contentTreeCmd(spec contentTreeCommandSpec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   spec.Use,
		Short: spec.Name + " reads",
		Long: fmt.Sprintf(`Cloud %s read operations.

These commands use documented Cloud v2 %s endpoints. Create/delete operations
are mutations and remain deferred behind explicit safety gates.

Examples:
  confluence %s view 12345 --json
  confluence %s children 12345 --type page --limit 25`, spec.Name, spec.Endpoint, spec.Use, spec.Use),
	}
	cmd.AddCommand(contentTreeViewCmd(spec))
	cmd.AddCommand(contentTreeChildrenCmd(spec))
	return cmd
}

func contentTreeViewCmd(spec contentTreeCommandSpec) *cobra.Command {
	var includeCollaborators bool
	var includeDirectChildren bool
	var includeOperations bool
	var includeProperties bool
	cmd := &cobra.Command{
		Use:   "view <id>",
		Short: "Show one Cloud " + strings.ToLower(spec.Name),
		Long: fmt.Sprintf(`Show one Cloud %s by id.

Examples:
  confluence %s view 12345
  confluence %s view 12345 --include-operations --include-properties --json`, strings.ToLower(spec.Name), spec.Use, spec.Use),
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
			entity, err := conf.GetContentTreeEntity(ctx, c, spec.Kind, args[0], conf.ContentTreeEntityOptions{
				IncludeCollaborators:  includeCollaborators,
				IncludeDirectChildren: includeDirectChildren,
				IncludeOperations:     includeOperations,
				IncludeProperties:     includeProperties,
			})
			if err != nil {
				return err
			}
			if w.IsJSON() {
				return w.JSON(entity)
			}
			printContentTreeEntity(w, *entity)
			return nil
		},
	}
	cmd.Flags().BoolVar(&includeCollaborators, "include-collaborators", false, "Include collaborators in the Cloud response")
	cmd.Flags().BoolVar(&includeDirectChildren, "include-direct-children", false, "Include direct children in the Cloud response")
	cmd.Flags().BoolVar(&includeOperations, "include-operations", false, "Include permitted operations in the Cloud response")
	cmd.Flags().BoolVar(&includeProperties, "include-properties", false, "Include properties in the Cloud response")
	return cmd
}

func contentTreeChildrenCmd(spec contentTreeCommandSpec) *cobra.Command {
	var limit int
	var sort string
	var types []string
	cmd := &cobra.Command{
		Use:   "children <id>",
		Short: "List direct children of a Cloud " + strings.ToLower(spec.Name),
		Long: fmt.Sprintf(`List direct children of one Cloud %s.

The Cloud v2 endpoint returns minimal child rows for databases, Smart Link
embeds, folders, pages, and whiteboards. Use the matching view command for
full details.

Examples:
  confluence %s children 12345
  confluence %s children 12345 --type page --type database --json
  confluence %s children 12345 --sort position --limit 100`, strings.ToLower(spec.Name), spec.Use, spec.Use, spec.Use),
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
			children, err := conf.ListContentTreeDirectChildren(ctx, c, spec.Kind, args[0], conf.DirectChildrenOptions{
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
			printContentTreeChildren(w, children)
			return nil
		},
	}
	cmd.Flags().IntVar(&limit, "limit", 50, "Max children (hard cap 200)")
	cmd.Flags().StringVar(&sort, "sort", "", "Cloud sort expression")
	cmd.Flags().StringSliceVar(&types, "type", nil, "Content type filter; repeatable or comma-separated")
	return cmd
}

type contentTreeTextWriter interface {
	Text(format string, args ...any)
}

func printContentTreeEntity(w contentTreeTextWriter, entity conf.ContentTreeEntity) {
	w.Text("%s\t%s\t%s\t%s", entity.ID, firstNonEmpty(entity.Type, "-"), firstNonEmpty(entity.Status, "-"), firstNonEmpty(entity.Title, "-"))
	if entity.SpaceID != "" {
		w.Text("\tspace=%s", entity.SpaceID)
	}
	if entity.ParentID != "" {
		w.Text("\tparent=%s", entity.ParentID)
	}
	if entity.EmbedURL != "" {
		w.Text("\turl=%s", entity.EmbedURL)
	}
	w.Text("\n")
}

func printContentTreeChildren(w contentTreeTextWriter, children []conf.Content) {
	for _, child := range children {
		position := "-"
		if child.ChildPosition > 0 {
			position = fmt.Sprintf("%d", child.ChildPosition)
		}
		w.Text("%s\t%s\tpos=%s\t%s\n", child.ID, firstNonEmpty(child.Type, "-"), position, firstNonEmpty(child.Title, "-"))
	}
}
