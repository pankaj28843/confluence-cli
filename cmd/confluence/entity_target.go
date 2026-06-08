package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

type entityTargetFlags struct {
	page          string
	blogpost      string
	attachment    string
	space         string
	footerComment string
	inlineComment string
	customContent string
	whiteboard    string
	database      string
	embed         string
	folder        string
}

type selectedEntityTarget struct {
	typ string
	id  string
}

func addOperationTargetFlags(cmd *cobra.Command, f *entityTargetFlags) {
	cmd.Flags().StringVar(&f.page, "page", "", "Page id")
	cmd.Flags().StringVar(&f.blogpost, "blogpost", "", "Blog post id")
	cmd.Flags().StringVar(&f.attachment, "attachment", "", "Attachment id")
	cmd.Flags().StringVar(&f.space, "space", "", "Cloud space key or id")
	cmd.Flags().StringVar(&f.footerComment, "footer-comment", "", "Footer comment id")
	cmd.Flags().StringVar(&f.inlineComment, "inline-comment", "", "Inline comment id")
	cmd.Flags().StringVar(&f.customContent, "custom-content", "", "Custom content id")
	cmd.Flags().StringVar(&f.whiteboard, "whiteboard", "", "Whiteboard id")
	cmd.Flags().StringVar(&f.database, "database", "", "Database id")
	cmd.Flags().StringVar(&f.embed, "embed", "", "Smart Link embed id")
	cmd.Flags().StringVar(&f.folder, "folder", "", "Folder id")
}

func addLikeTargetFlags(cmd *cobra.Command, f *entityTargetFlags) {
	cmd.Flags().StringVar(&f.page, "page", "", "Page id")
	cmd.Flags().StringVar(&f.blogpost, "blogpost", "", "Blog post id")
	cmd.Flags().StringVar(&f.footerComment, "footer-comment", "", "Footer comment id")
	cmd.Flags().StringVar(&f.inlineComment, "inline-comment", "", "Inline comment id")
}

func selectedOperationTarget(f entityTargetFlags) (selectedEntityTarget, error) {
	return selectExactlyOneTarget([]selectedEntityTarget{
		{typ: "page", id: f.page},
		{typ: "blogpost", id: f.blogpost},
		{typ: "attachment", id: f.attachment},
		{typ: "space", id: f.space},
		{typ: "footer-comment", id: f.footerComment},
		{typ: "inline-comment", id: f.inlineComment},
		{typ: "custom-content", id: f.customContent},
		{typ: "whiteboard", id: f.whiteboard},
		{typ: "database", id: f.database},
		{typ: "embed", id: f.embed},
		{typ: "folder", id: f.folder},
	})
}

func selectedLikeTarget(f entityTargetFlags) (selectedEntityTarget, error) {
	return selectExactlyOneTarget([]selectedEntityTarget{
		{typ: "page", id: f.page},
		{typ: "blogpost", id: f.blogpost},
		{typ: "footer-comment", id: f.footerComment},
		{typ: "inline-comment", id: f.inlineComment},
	})
}

func selectExactlyOneTarget(candidates []selectedEntityTarget) (selectedEntityTarget, error) {
	var selected []selectedEntityTarget
	for _, candidate := range candidates {
		if candidate.id != "" {
			selected = append(selected, candidate)
		}
	}
	if len(selected) != 1 {
		return selectedEntityTarget{}, fmt.Errorf("exactly one target flag is required")
	}
	return selected[0], nil
}
