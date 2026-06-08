package main

import (
	"github.com/spf13/cobra"

	"github.com/pankaj28843/confluence-cli/internal/conf"
)

func templateCmdReal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "template",
		Short: "Cloud content and blueprint templates",
		Long: `Cloud template read operations.

Confluence Server/Data Center template endpoints are not exposed in the
current official REST OpenAPI, so typed template commands are Cloud-only.

Examples:
  confluence template list --limit 10
  confluence template blueprint list --space ENG --json
  confluence template view 12345 --expand body.storage`,
	}
	cmd.AddCommand(templateListCmd())
	cmd.AddCommand(templateBlueprintCmd())
	cmd.AddCommand(templateViewCmd())
	return cmd
}

func templateListCmd() *cobra.Command {
	var spaceKey string
	var expand []string
	var limit int
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List Cloud content templates",
		Long: `List Confluence Cloud content templates. Use --space to list templates in
one space, or omit it for global templates.

Examples:
  confluence template list
  confluence template list --space ENG --limit 10 --json
  confluence template list --expand body.storage --expand space`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := newContext()
			defer cancel()
			w := getWriter()
			defer w.Finish()
			c, err := newClient()
			if err != nil {
				return err
			}
			templates, err := conf.ListContentTemplates(ctx, c, conf.TemplateListOptions{
				SpaceKey: spaceKey,
				Expand:   expand,
				Limit:    limit,
			})
			if err != nil {
				return err
			}
			if w.IsJSON() {
				return w.JSON(templates)
			}
			printTemplates(w, templates)
			return nil
		},
	}
	cmd.Flags().StringVar(&spaceKey, "space", "", "Space key; omit for global templates")
	cmd.Flags().StringSliceVar(&expand, "expand", nil, "Expand value; repeatable or comma-separated")
	cmd.Flags().IntVar(&limit, "limit", 25, "Max templates (hard cap 200)")
	return cmd
}

func templateBlueprintCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "blueprint",
		Short: "Cloud blueprint templates",
		Long: `Cloud blueprint template operations.

Examples:
  confluence template blueprint list
  confluence template blueprint list --space ENG --limit 10 --json`,
	}
	cmd.AddCommand(templateBlueprintListCmd())
	return cmd
}

func templateBlueprintListCmd() *cobra.Command {
	var spaceKey string
	var expand []string
	var limit int
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List Cloud blueprint templates",
		Long: `List Confluence Cloud blueprint templates. Use --space to list blueprints in
one space, or omit it for global blueprints.

Examples:
  confluence template blueprint list
  confluence template blueprint list --space ENG --limit 10 --json
  confluence template blueprint list --expand body.storage`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := newContext()
			defer cancel()
			w := getWriter()
			defer w.Finish()
			c, err := newClient()
			if err != nil {
				return err
			}
			templates, err := conf.ListBlueprintTemplates(ctx, c, conf.TemplateListOptions{
				SpaceKey: spaceKey,
				Expand:   expand,
				Limit:    limit,
			})
			if err != nil {
				return err
			}
			if w.IsJSON() {
				return w.JSON(templates)
			}
			printTemplates(w, templates)
			return nil
		},
	}
	cmd.Flags().StringVar(&spaceKey, "space", "", "Space key; omit for global blueprints")
	cmd.Flags().StringSliceVar(&expand, "expand", nil, "Expand value; repeatable or comma-separated")
	cmd.Flags().IntVar(&limit, "limit", 25, "Max templates (hard cap 200)")
	return cmd
}

func templateViewCmd() *cobra.Command {
	var expand []string
	cmd := &cobra.Command{
		Use:   "view <id>",
		Short: "Show one Cloud content template",
		Long: `Show one Confluence Cloud content template by id.

Examples:
  confluence template view 12345
  confluence template view 12345 --expand body.storage --json`,
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
			template, err := conf.GetContentTemplate(ctx, c, args[0], expand)
			if err != nil {
				return err
			}
			if w.IsJSON() {
				return w.JSON(template)
			}
			printTemplate(w, *template)
			return nil
		},
	}
	cmd.Flags().StringSliceVar(&expand, "expand", nil, "Expand value; repeatable or comma-separated")
	return cmd
}

type templateTextWriter interface {
	Text(format string, args ...interface{})
}

func printTemplates(w templateTextWriter, templates []conf.ContentTemplate) {
	for _, template := range templates {
		printTemplate(w, template)
	}
}

func printTemplate(w templateTextWriter, template conf.ContentTemplate) {
	templateType := firstNonEmpty(template.TemplateType, "-")
	name := firstNonEmpty(template.Name, "-")
	w.Text("%s\t%s\t%s", template.TemplateID, templateType, name)
	if template.ReferencingBlueprint != "" {
		w.Text("\t%s", template.ReferencingBlueprint)
	}
	w.Text("\n")
}
