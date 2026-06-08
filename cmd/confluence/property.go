package main

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/pankaj28843/confluence-cli/internal/conf"
)

func propertyCmdReal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "property",
		Short: "Content and space properties",
		Long: `Content and space property operations.

Examples:
  confluence property content list --page 12345
  confluence property content get --page 12345 --key release
  confluence property space set --space ENG --key retention --value '{"days":30}'`,
	}
	cmd.AddCommand(propertyContentCmd())
	cmd.AddCommand(propertySpaceCmd())
	return cmd
}

func propertyContentCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "content",
		Short: "Content properties (list, get, set, delete)",
		Long: `Content property operations for a page/content id.

Examples:
  confluence property content list --page 12345
  confluence property content set --page 12345 --key release --value '{"ready":true}'`,
	}
	cmd.AddCommand(propertyContentListCmd())
	cmd.AddCommand(propertyContentGetCmd())
	cmd.AddCommand(propertyContentSetCmd())
	cmd.AddCommand(propertyContentDeleteCmd())
	return cmd
}

func propertySpaceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "space",
		Short: "Space properties (list, get, set, delete)",
		Long: `Space property operations for a space key.

Examples:
  confluence property space list --space ENG
  confluence property space set --space ENG --key retention --value '{"days":30}'`,
	}
	cmd.AddCommand(propertySpaceListCmd())
	cmd.AddCommand(propertySpaceGetCmd())
	cmd.AddCommand(propertySpaceSetCmd())
	cmd.AddCommand(propertySpaceDeleteCmd())
	return cmd
}

func propertyContentListCmd() *cobra.Command {
	var page, key string
	var limit int
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List content properties",
		Long: `List properties for a page/content id.

Examples:
  confluence property content list --page 12345
  confluence property content list --page 12345 --key release --json
  confluence property content list --page 12345 --limit 100`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := newContext()
			defer cancel()
			w := getWriter()
			defer w.Finish()
			if page == "" {
				return fmt.Errorf("--page is required")
			}
			c, err := newClient()
			if err != nil {
				return err
			}
			props, err := conf.ListContentProperties(ctx, c, page, key, limit)
			if err != nil {
				return err
			}
			if w.IsJSON() {
				return w.JSON(props)
			}
			printProperties(w, props)
			return nil
		},
	}
	cmd.Flags().StringVar(&page, "page", "", "Content id (required)")
	cmd.Flags().StringVar(&key, "key", "", "Optional property key filter")
	cmd.Flags().IntVar(&limit, "limit", 25, "Max properties (hard cap 200)")
	return cmd
}

func propertyContentGetCmd() *cobra.Command {
	var page, key string
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get one content property by key",
		Long: `Get a property for a page/content id by key.

Examples:
  confluence property content get --page 12345 --key release
  confluence property content get --page 12345 --key release --jq '.value'`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := newContext()
			defer cancel()
			w := getWriter()
			defer w.Finish()
			if page == "" || key == "" {
				return fmt.Errorf("--page and --key are required")
			}
			c, err := newClient()
			if err != nil {
				return err
			}
			prop, err := conf.GetContentProperty(ctx, c, page, key)
			if err != nil {
				return err
			}
			if w.IsJSON() {
				return w.JSON(prop)
			}
			printProperty(w, *prop)
			return nil
		},
	}
	cmd.Flags().StringVar(&page, "page", "", "Content id (required)")
	cmd.Flags().StringVar(&key, "key", "", "Property key (required)")
	return cmd
}

func propertyContentSetCmd() *cobra.Command {
	var page, key, valueFlag string
	cmd := &cobra.Command{
		Use:   "set",
		Short: "Create or update one content property",
		Long: `Create or update a property for a page/content id. --value must be JSON.
Use JSON strings for scalar text values, for example --value '"done"'.

Examples:
  confluence property content set --page 12345 --key release --value '{"ready":true}'
  confluence property content set --page 12345 --key owner --value '"platform"'
  confluence property content set --page 12345 --key release --value @property.json`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := newContext()
			defer cancel()
			w := getWriter()
			defer w.Finish()
			if page == "" || key == "" || valueFlag == "" {
				return fmt.Errorf("--page, --key, and --value are required")
			}
			value, err := parsePropertyValue(valueFlag)
			if err != nil {
				return err
			}
			c, err := newClient()
			if err != nil {
				return err
			}
			prop, err := conf.SetContentProperty(ctx, c, page, key, value)
			if err != nil {
				return err
			}
			if w.IsJSON() {
				return w.JSON(prop)
			}
			printProperty(w, *prop)
			return nil
		},
	}
	cmd.Flags().StringVar(&page, "page", "", "Content id (required)")
	cmd.Flags().StringVar(&key, "key", "", "Property key (required)")
	cmd.Flags().StringVar(&valueFlag, "value", "", "JSON value, @file, or @- (required)")
	return cmd
}

func propertyContentDeleteCmd() *cobra.Command {
	var page, key string
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete one content property by key",
		Long: `Delete a property for a page/content id by key.

Examples:
  confluence property content delete --page 12345 --key release
  confluence property content delete --page 12345 --key release --json`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := newContext()
			defer cancel()
			w := getWriter()
			defer w.Finish()
			if page == "" || key == "" {
				return fmt.Errorf("--page and --key are required")
			}
			c, err := newClient()
			if err != nil {
				return err
			}
			if err := conf.DeleteContentProperty(ctx, c, page, key); err != nil {
				return err
			}
			if w.IsJSON() {
				return w.JSON(map[string]any{"ok": true, "page": page, "key": key})
			}
			w.Text("deleted property %q from page %s\n", key, page)
			return nil
		},
	}
	cmd.Flags().StringVar(&page, "page", "", "Content id (required)")
	cmd.Flags().StringVar(&key, "key", "", "Property key (required)")
	return cmd
}

func propertySpaceListCmd() *cobra.Command {
	var space, key string
	var limit int
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List space properties",
		Long: `List properties for a space key.

Examples:
  confluence property space list --space ENG
  confluence property space list --space ENG --key retention --json
  confluence property space list --space ENG --limit 100`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := newContext()
			defer cancel()
			w := getWriter()
			defer w.Finish()
			if space == "" {
				return fmt.Errorf("--space is required")
			}
			c, err := newClient()
			if err != nil {
				return err
			}
			props, err := conf.ListSpaceProperties(ctx, c, space, key, limit)
			if err != nil {
				return err
			}
			if w.IsJSON() {
				return w.JSON(props)
			}
			printProperties(w, props)
			return nil
		},
	}
	cmd.Flags().StringVar(&space, "space", "", "Space key (required)")
	cmd.Flags().StringVar(&key, "key", "", "Optional property key filter")
	cmd.Flags().IntVar(&limit, "limit", 25, "Max properties (hard cap 200)")
	return cmd
}

func propertySpaceGetCmd() *cobra.Command {
	var space, key string
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get one space property by key",
		Long: `Get a space property by key.

Examples:
  confluence property space get --space ENG --key retention
  confluence property space get --space ENG --key retention --jq '.value'`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := newContext()
			defer cancel()
			w := getWriter()
			defer w.Finish()
			if space == "" || key == "" {
				return fmt.Errorf("--space and --key are required")
			}
			c, err := newClient()
			if err != nil {
				return err
			}
			prop, err := conf.GetSpaceProperty(ctx, c, space, key)
			if err != nil {
				return err
			}
			if w.IsJSON() {
				return w.JSON(prop)
			}
			printProperty(w, *prop)
			return nil
		},
	}
	cmd.Flags().StringVar(&space, "space", "", "Space key (required)")
	cmd.Flags().StringVar(&key, "key", "", "Property key (required)")
	return cmd
}

func propertySpaceSetCmd() *cobra.Command {
	var space, key, valueFlag string
	cmd := &cobra.Command{
		Use:   "set",
		Short: "Create or update one space property",
		Long: `Create or update a space property. --value must be JSON.
Use JSON strings for scalar text values, for example --value '"platform"'.

Examples:
  confluence property space set --space ENG --key retention --value '{"days":30}'
  confluence property space set --space ENG --key owner --value '"platform"'
  confluence property space set --space ENG --key retention --value @property.json`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := newContext()
			defer cancel()
			w := getWriter()
			defer w.Finish()
			if space == "" || key == "" || valueFlag == "" {
				return fmt.Errorf("--space, --key, and --value are required")
			}
			value, err := parsePropertyValue(valueFlag)
			if err != nil {
				return err
			}
			c, err := newClient()
			if err != nil {
				return err
			}
			prop, err := conf.SetSpaceProperty(ctx, c, space, key, value)
			if err != nil {
				return err
			}
			if w.IsJSON() {
				return w.JSON(prop)
			}
			printProperty(w, *prop)
			return nil
		},
	}
	cmd.Flags().StringVar(&space, "space", "", "Space key (required)")
	cmd.Flags().StringVar(&key, "key", "", "Property key (required)")
	cmd.Flags().StringVar(&valueFlag, "value", "", "JSON value, @file, or @- (required)")
	return cmd
}

func propertySpaceDeleteCmd() *cobra.Command {
	var space, key string
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete one space property by key",
		Long: `Delete a space property by key.

Examples:
  confluence property space delete --space ENG --key retention
  confluence property space delete --space ENG --key retention --json`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := newContext()
			defer cancel()
			w := getWriter()
			defer w.Finish()
			if space == "" || key == "" {
				return fmt.Errorf("--space and --key are required")
			}
			c, err := newClient()
			if err != nil {
				return err
			}
			if err := conf.DeleteSpaceProperty(ctx, c, space, key); err != nil {
				return err
			}
			if w.IsJSON() {
				return w.JSON(map[string]any{"ok": true, "space": space, "key": key})
			}
			w.Text("deleted property %q from space %s\n", key, space)
			return nil
		},
	}
	cmd.Flags().StringVar(&space, "space", "", "Space key (required)")
	cmd.Flags().StringVar(&key, "key", "", "Property key (required)")
	return cmd
}

func parsePropertyValue(valueFlag string) (any, error) {
	raw, err := readDataArg(valueFlag)
	if err != nil {
		return nil, err
	}
	if len(raw) == 0 {
		return nil, fmt.Errorf("--value is empty")
	}
	var value any
	if err := json.Unmarshal(raw, &value); err != nil {
		return nil, fmt.Errorf("--value is not valid JSON: %w", err)
	}
	return value, nil
}

type propertyTextWriter interface {
	Text(format string, args ...interface{})
}

func printProperties(w propertyTextWriter, props []conf.Property) {
	for _, prop := range props {
		w.Text("%s\t%s\t%d\n", prop.ID, prop.Key, prop.Version.Number)
	}
}

func printProperty(w propertyTextWriter, prop conf.Property) {
	w.Text("%s\t%s", prop.ID, prop.Key)
	if prop.Version.Number > 0 {
		w.Text("\tv%d", prop.Version.Number)
	}
	w.Text("\n")
}
