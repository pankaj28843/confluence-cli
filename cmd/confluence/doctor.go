package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/pankaj28843/confluence-cli/internal/client"
	"github.com/pankaj28843/confluence-cli/internal/conf"
)

func init() {
	runDoctor = doctorRun
}

type doctorReport struct {
	BaseURL string      `json:"baseUrl"`
	Flavor  string      `json:"flavor"`
	User    interface{} `json:"authenticatedUser"`
	OK      bool        `json:"ok"`
}

func doctorRun(cmd *cobra.Command, args []string) error {
	ctx, cancel := newContext()
	defer cancel()
	w := getWriter()
	defer w.Finish()

	cfg, err := client.FromEnv()
	if err != nil {
		return newConfigError(err)
	}
	cfg.Debug = debug
	c, err := client.New(cfg)
	if err != nil {
		return newConfigError(err)
	}

	report := doctorReport{BaseURL: cfg.BaseURL, Flavor: cfg.Flavor.String()}

	user, err := conf.GetCurrentUser(ctx, c)
	if err != nil {
		if client.IsUserFixable(err) {
			hint := "https://id.atlassian.com/manage-profile/security/api-tokens"
			if cfg.Flavor == client.FlavorServer {
				hint = cfg.BaseURL + "/plugins/personalaccesstokens/usertokens.action"
			}
			return newConfigError(fmt.Errorf("authentication failed: %w\nHint: regenerate your credentials at %s", err, hint))
		}
		return fmt.Errorf("current user: %w", err)
	}
	report.User = user
	report.OK = user.ID() != ""

	if w.IsJSON() {
		return w.JSON(report)
	}
	w.Text("flavor:               %s\n", report.Flavor)
	w.Text("baseUrl:              %s\n", report.BaseURL)
	w.Text("authenticated as:     %s (%s)\n", user.Label(), user.ID())
	if report.OK {
		w.Text("\nOK — confluence is ready to use.\n")
	} else {
		w.Text("\nReached the server but no user identifier came back — check flavor.\n")
	}
	return nil
}
