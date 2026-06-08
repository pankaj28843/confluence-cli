package main

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/pankaj28843/confluence-cli/internal/conf"
)

func permissionCmdReal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "permission",
		Short: "Space permission reads",
		Long: `Space permission helpers.

Inspect documented space permission assignments across Confluence Cloud and
Server/Data Center. Subject-specific reads are Server/Data Center only.

Examples:
  confluence permission space list --space ENG --json
  confluence permission space available --limit 100 --json
  confluence permission space subject --space ENG --group confluence-users --json`,
	}
	cmd.AddCommand(permissionSpaceCmd())
	return cmd
}

func permissionSpaceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "space",
		Short: "Space permission reads",
		Long: `Space permission reads.

Examples:
  confluence permission space list --space ENG
  confluence permission space available --json
  confluence permission space subject --space ENG --anonymous`,
	}
	cmd.AddCommand(permissionSpaceListCmd())
	cmd.AddCommand(permissionSpaceAvailableCmd())
	cmd.AddCommand(permissionSpaceSubjectCmd())
	return cmd
}

func permissionSpaceListCmd() *cobra.Command {
	var space string
	var limit int
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List permissions assigned on a space",
		Long: `List permissions assigned on a space.

On Cloud this uses the v2 space permission assignment endpoint. A space key is
resolved to its Cloud space id before reading permissions. On Server/Data Center
this uses the space permissions endpoint.

Examples:
  confluence permission space list --space ENG
  confluence permission space list --space ENG --limit 100 --json`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := newContext()
			defer cancel()
			w := getWriter()
			defer w.Finish()
			if strings.TrimSpace(space) == "" {
				return fmt.Errorf("--space is required")
			}
			c, err := newClient()
			if err != nil {
				return err
			}
			assignments, err := conf.ListSpacePermissionAssignments(ctx, c, space, limit)
			if err != nil {
				return err
			}
			if w.IsJSON() {
				return w.JSON(assignments)
			}
			return writeSpacePermissionAssignments(w, assignments)
		},
	}
	cmd.Flags().StringVar(&space, "space", "", "Space key or Cloud space id (required)")
	cmd.Flags().IntVar(&limit, "limit", 25, "Max permission assignments (hard cap 200)")
	return cmd
}

func permissionSpaceAvailableCmd() *cobra.Command {
	var limit int
	cmd := &cobra.Command{
		Use:   "available",
		Short: "List available Cloud space permissions",
		Long: `List available Cloud space permissions.

This is the documented Cloud v2 permission catalog for tenants with RBAC.

Examples:
  confluence permission space available --json
  confluence permission space available --limit 100`,
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
			permissions, err := conf.ListAvailableSpacePermissions(ctx, c, limit)
			if err != nil {
				return err
			}
			if w.IsJSON() {
				return w.JSON(permissions)
			}
			return writeSpacePermissionDefinitions(w, permissions)
		},
	}
	cmd.Flags().IntVar(&limit, "limit", 25, "Max permissions (hard cap 200)")
	return cmd
}

func permissionSpaceSubjectCmd() *cobra.Command {
	var space string
	var anonymous bool
	var group string
	var userKey string
	var limit int
	cmd := &cobra.Command{
		Use:   "subject",
		Short: "List Server/Data Center space permissions for one subject",
		Long: `List Server/Data Center space permissions for one subject.

Select exactly one subject with --anonymous, --group, or --user-key.

Examples:
  confluence permission space subject --space ENG --anonymous
  confluence permission space subject --space ENG --group confluence-users --json
  confluence permission space subject --space ENG --user-key ada --limit 100`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := newContext()
			defer cancel()
			w := getWriter()
			defer w.Finish()
			if strings.TrimSpace(space) == "" {
				return fmt.Errorf("--space is required")
			}
			c, err := newClient()
			if err != nil {
				return err
			}
			assignments, err := conf.ListSpacePermissionsForSubject(ctx, c, space, conf.SpacePermissionSubjectSelector{
				Anonymous: anonymous,
				GroupName: group,
				UserKey:   userKey,
			}, limit)
			if err != nil {
				return err
			}
			if w.IsJSON() {
				return w.JSON(assignments)
			}
			return writeSpacePermissionAssignments(w, assignments)
		},
	}
	cmd.Flags().StringVar(&space, "space", "", "Space key (required)")
	cmd.Flags().BoolVar(&anonymous, "anonymous", false, "Read anonymous subject permissions")
	cmd.Flags().StringVar(&group, "group", "", "Read permissions for this group name")
	cmd.Flags().StringVar(&userKey, "user-key", "", "Read permissions for this Server/Data Center user key")
	cmd.Flags().IntVar(&limit, "limit", 25, "Max permission assignments (hard cap 200)")
	return cmd
}

type permissionTextWriter interface {
	Text(format string, args ...any)
}

func writeSpacePermissionAssignments(w permissionTextWriter, assignments []conf.SpacePermissionAssignment) error {
	for _, assignment := range assignments {
		w.Text("%s\t%s\n", spacePermissionSubjectLabel(assignment), spacePermissionOperationLabel(assignment.Operation))
	}
	return nil
}

func writeSpacePermissionDefinitions(w permissionTextWriter, permissions []conf.SpacePermissionDefinition) error {
	for _, permission := range permissions {
		w.Text("%s\t%s\n", permission.ID, permission.DisplayName)
	}
	return nil
}

func spacePermissionSubjectLabel(assignment conf.SpacePermissionAssignment) string {
	if assignment.Principal.Type != "" || assignment.Principal.ID != "" {
		return compactPair(assignment.Principal.Type, assignment.Principal.ID)
	}
	if assignment.Subject.Type != "" || assignment.Subject.DisplayName != "" {
		if assignment.Subject.DisplayName != "" {
			return compactPair(assignment.Subject.Type, assignment.Subject.DisplayName)
		}
		if assignment.Subject.GroupName != "" {
			return compactPair(assignment.Subject.Type, assignment.Subject.GroupName)
		}
		if assignment.Subject.UserKey != "" {
			return compactPair(assignment.Subject.Type, assignment.Subject.UserKey)
		}
		return assignment.Subject.Type
	}
	return "-"
}

func spacePermissionOperationLabel(operation conf.SpacePermissionOperation) string {
	name := operation.Name()
	if operation.TargetType == "" {
		if name == "" {
			return "-"
		}
		return name
	}
	if name == "" {
		return operation.TargetType
	}
	return operation.TargetType + "." + name
}

func compactPair(left, right string) string {
	switch {
	case left == "":
		return right
	case right == "":
		return left
	default:
		return left + ":" + right
	}
}
