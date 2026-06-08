package main

import (
	"github.com/spf13/cobra"

	"github.com/pankaj28843/confluence-cli/internal/conf"
)

func taskCmdReal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "task",
		Short: "Tasks (list, view, complete, reopen, long)",
		Long: `Task operations.

Cloud content tasks are available through list/view/complete/reopen.
Server/Data Center long-running tasks are available under task long.

Examples:
  confluence task list --page 12345 --status incomplete
  confluence task view 42 --body-format storage
  confluence task long list --limit 10`,
	}
	cmd.AddCommand(taskListCmd())
	cmd.AddCommand(taskViewCmd())
	cmd.AddCommand(taskCompleteCmd())
	cmd.AddCommand(taskReopenCmd())
	cmd.AddCommand(taskLongCmd())
	return cmd
}

func taskListCmd() *cobra.Command {
	var status, pageID, blogPostID, bodyFormat string
	var taskIDs, spaceIDs, createdBy, assignedTo, completedBy []string
	var includeBlank bool
	var limit int
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List Cloud content tasks",
		Long: `List Confluence Cloud content tasks. Server/Data Center long-running tasks
are under task long list.

Examples:
  confluence task list --status incomplete --limit 50
  confluence task list --page 12345 --body-format storage --json
  confluence task list --assigned-to 557058:abc --include-blank`,
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
			tasks, err := conf.ListTasks(ctx, c, conf.TaskFilter{
				Status:            status,
				TaskIDs:           taskIDs,
				SpaceIDs:          spaceIDs,
				PageID:            pageID,
				BlogPostID:        blogPostID,
				CreatedBy:         createdBy,
				AssignedTo:        assignedTo,
				CompletedBy:       completedBy,
				IncludeBlankTasks: includeBlank,
				BodyFormat:        bodyFormat,
				Limit:             limit,
			})
			if err != nil {
				return err
			}
			if w.IsJSON() {
				return w.JSON(tasks)
			}
			printTasks(w, tasks)
			return nil
		},
	}
	cmd.Flags().StringVar(&status, "status", "", "Filter: complete | incomplete")
	cmd.Flags().StringSliceVar(&taskIDs, "task-id", nil, "Task id filter; repeatable or comma-separated")
	cmd.Flags().StringSliceVar(&spaceIDs, "space-id", nil, "Cloud space id filter; repeatable or comma-separated")
	cmd.Flags().StringVar(&pageID, "page", "", "Page id filter")
	cmd.Flags().StringVar(&blogPostID, "blogpost", "", "Blog post id filter")
	cmd.Flags().StringSliceVar(&createdBy, "created-by", nil, "Creator account id filter; repeatable or comma-separated")
	cmd.Flags().StringSliceVar(&assignedTo, "assigned-to", nil, "Assignee account id filter; repeatable or comma-separated")
	cmd.Flags().StringSliceVar(&completedBy, "completed-by", nil, "Completer account id filter; repeatable or comma-separated")
	cmd.Flags().BoolVar(&includeBlank, "include-blank", false, "Include blank tasks")
	cmd.Flags().StringVar(&bodyFormat, "body-format", "", "Body format: storage | atlas_doc_format")
	cmd.Flags().IntVar(&limit, "limit", 25, "Max tasks (hard cap 200)")
	return cmd
}

func taskViewCmd() *cobra.Command {
	var bodyFormat string
	cmd := &cobra.Command{
		Use:   "view <id>",
		Short: "Show one Cloud content task",
		Long: `Show one Confluence Cloud content task by id.

Examples:
  confluence task view 42
  confluence task view 42 --body-format storage --json`,
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
			task, err := conf.GetTask(ctx, c, args[0], bodyFormat)
			if err != nil {
				return err
			}
			if w.IsJSON() {
				return w.JSON(task)
			}
			printTask(w, *task)
			return nil
		},
	}
	cmd.Flags().StringVar(&bodyFormat, "body-format", "", "Body format: storage | atlas_doc_format")
	return cmd
}

func taskCompleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "complete <id>",
		Short: "Mark one Cloud content task complete",
		Long: `Mark one Confluence Cloud content task complete.

Examples:
  confluence task complete 42
  confluence task complete 42 --json`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runTaskStatusUpdate(args[0], "complete")
		},
	}
	return cmd
}

func taskReopenCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reopen <id>",
		Short: "Mark one Cloud content task incomplete",
		Long: `Mark one Confluence Cloud content task incomplete.

Examples:
  confluence task reopen 42
  confluence task reopen 42 --json`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runTaskStatusUpdate(args[0], "incomplete")
		},
	}
	return cmd
}

func runTaskStatusUpdate(id, status string) error {
	ctx, cancel := newContext()
	defer cancel()
	w := getWriter()
	defer w.Finish()
	c, err := newClient()
	if err != nil {
		return err
	}
	task, err := conf.UpdateTaskStatus(ctx, c, id, status)
	if err != nil {
		return err
	}
	if w.IsJSON() {
		return w.JSON(task)
	}
	printTask(w, *task)
	return nil
}

func taskLongCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "long",
		Short: "Server/Data Center long-running tasks",
		Long: `Server/Data Center long-running task operations.

Examples:
  confluence task long list --limit 10
  confluence task long view 123456 --expand messages`,
	}
	cmd.AddCommand(taskLongListCmd())
	cmd.AddCommand(taskLongViewCmd())
	return cmd
}

func taskLongListCmd() *cobra.Command {
	var expand string
	var limit int
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List Server/Data Center long-running tasks",
		Long: `List Server/Data Center long-running tasks.

Examples:
  confluence task long list
  confluence task long list --expand messages --limit 10 --json`,
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
			tasks, err := conf.ListLongTasks(ctx, c, conf.LongTaskFilter{Expand: expand, Limit: limit})
			if err != nil {
				return err
			}
			if w.IsJSON() {
				return w.JSON(tasks)
			}
			printLongTasks(w, tasks)
			return nil
		},
	}
	cmd.Flags().StringVar(&expand, "expand", "", "Expand parameter, for example messages")
	cmd.Flags().IntVar(&limit, "limit", 25, "Max long tasks (hard cap 200)")
	return cmd
}

func taskLongViewCmd() *cobra.Command {
	var expand string
	cmd := &cobra.Command{
		Use:   "view <id>",
		Short: "Show one Server/Data Center long-running task",
		Long: `Show one Server/Data Center long-running task by id.

Examples:
  confluence task long view 123456
  confluence task long view 123456 --expand messages --json`,
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
			task, err := conf.GetLongTask(ctx, c, args[0], expand)
			if err != nil {
				return err
			}
			if w.IsJSON() {
				return w.JSON(task)
			}
			printLongTask(w, *task)
			return nil
		},
	}
	cmd.Flags().StringVar(&expand, "expand", "", "Expand parameter, for example messages")
	return cmd
}

type taskTextWriter interface {
	Text(format string, args ...interface{})
}

func printTasks(w taskTextWriter, tasks []conf.Task) {
	for _, task := range tasks {
		printTask(w, task)
	}
}

func printTask(w taskTextWriter, task conf.Task) {
	container := firstNonEmpty(task.PageID, task.BlogPostID)
	if container == "" {
		container = "-"
	}
	w.Text("%s\t%s\t%s", task.ID, task.Status, container)
	if task.AssignedTo != "" {
		w.Text("\t%s", task.AssignedTo)
	}
	if task.DueAt != "" {
		w.Text("\t%s", task.DueAt)
	}
	w.Text("\n")
}

func printLongTasks(w taskTextWriter, tasks []conf.LongTask) {
	for _, task := range tasks {
		printLongTask(w, task)
	}
}

func printLongTask(w taskTextWriter, task conf.LongTask) {
	name := firstNonEmpty(task.Name.Translation, task.Name.Key)
	if name == "" {
		name = "-"
	}
	w.Text("%s\t%d\t%t\t%s\n", task.ID, task.PercentageComplete, task.Successful, name)
}
