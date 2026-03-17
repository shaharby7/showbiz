package commands

import (
	"fmt"

	showbiz "github.com/showbiz-io/showbiz/sdk/go"
	"github.com/shaharby7/showbiz/cli/internal/output"
	"github.com/spf13/cobra"
)

var projectCmd = &cobra.Command{
	Use:   "project",
	Short: "Manage projects",
}

func init() {
	// project list
	projectListCmd := &cobra.Command{
		Use:   "list",
		Short: "List projects",
		RunE:  runProjectList,
	}
	projectListCmd.Flags().String("org", "", "Organization ID")

	// project create
	projectCreateCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a project",
		RunE:  runProjectCreate,
	}
	projectCreateCmd.Flags().String("org", "", "Organization ID")
	projectCreateCmd.Flags().String("name", "", "Project name")
	projectCreateCmd.Flags().String("description", "", "Project description")
	_ = projectCreateCmd.MarkFlagRequired("name")

	// project get
	projectGetCmd := &cobra.Command{
		Use:   "get <projectID>",
		Short: "Get project details",
		Args:  cobra.ExactArgs(1),
		RunE:  runProjectGet,
	}
	projectGetCmd.Flags().String("org", "", "Organization ID")

	// project update
	projectUpdateCmd := &cobra.Command{
		Use:   "update <projectID>",
		Short: "Update a project",
		Args:  cobra.ExactArgs(1),
		RunE:  runProjectUpdate,
	}
	projectUpdateCmd.Flags().String("org", "", "Organization ID")
	projectUpdateCmd.Flags().String("description", "", "New description")
	_ = projectUpdateCmd.MarkFlagRequired("description")

	// project delete
	projectDeleteCmd := &cobra.Command{
		Use:   "delete <projectID>",
		Short: "Delete a project",
		Args:  cobra.ExactArgs(1),
		RunE:  runProjectDelete,
	}
	projectDeleteCmd.Flags().String("org", "", "Organization ID")

	projectCmd.AddCommand(projectListCmd, projectCreateCmd, projectGetCmd, projectUpdateCmd, projectDeleteCmd)
}

func runProjectList(cmd *cobra.Command, args []string) error {
	p := getPrinter(cmd)
	orgID, err := getOrg(cmd)
	if err != nil {
		return err
	}

	client, err := newClient()
	if err != nil {
		return err
	}

	result, err := client.ListProjects(cmd.Context(), orgID, nil)
	if err != nil {
		return err
	}

	t := p.Table("ID", "NAME", "DESCRIPTION")
	for _, proj := range result.Data {
		t.AddRow(proj.ID, proj.Name, proj.Description)
	}
	return t.Print()
}

func runProjectCreate(cmd *cobra.Command, args []string) error {
	p := getPrinter(cmd)
	orgID, err := getOrg(cmd)
	if err != nil {
		return err
	}

	client, err := newClient()
	if err != nil {
		return err
	}

	name, _ := cmd.Flags().GetString("name")
	description, _ := cmd.Flags().GetString("description")

	proj, err := client.CreateProject(cmd.Context(), orgID, showbiz.CreateProjectInput{
		Name:        name,
		Description: description,
	})
	if err != nil {
		return err
	}

	p.Success("Created project %s (ID: %s)", proj.Name, proj.ID)
	return nil
}

func runProjectGet(cmd *cobra.Command, args []string) error {
	p := getPrinter(cmd)
	orgID, err := getOrg(cmd)
	if err != nil {
		return err
	}

	client, err := newClient()
	if err != nil {
		return err
	}

	proj, err := client.GetProject(cmd.Context(), orgID, args[0])
	if err != nil {
		return err
	}

	return p.PrintResource([]output.KeyValue{
		{Key: "ID", Value: proj.ID},
		{Key: "Name", Value: proj.Name},
		{Key: "Organization", Value: proj.OrganizationID},
		{Key: "Description", Value: proj.Description},
		{Key: "Created At", Value: proj.CreatedAt},
		{Key: "Updated At", Value: proj.UpdatedAt},
	})
}

func runProjectUpdate(cmd *cobra.Command, args []string) error {
	p := getPrinter(cmd)
	orgID, err := getOrg(cmd)
	if err != nil {
		return err
	}

	client, err := newClient()
	if err != nil {
		return err
	}

	description, _ := cmd.Flags().GetString("description")
	proj, err := client.UpdateProject(cmd.Context(), orgID, args[0], showbiz.UpdateProjectInput{
		Description: description,
	})
	if err != nil {
		return err
	}

	p.Success("Updated project %s", proj.ID)
	return nil
}

func runProjectDelete(cmd *cobra.Command, args []string) error {
	p := getPrinter(cmd)
	orgID, err := getOrg(cmd)
	if err != nil {
		return err
	}

	if !confirmAction(cmd, fmt.Sprintf("Delete project %s?", args[0])) {
		return nil
	}

	client, err := newClient()
	if err != nil {
		return err
	}

	if err := client.DeleteProject(cmd.Context(), orgID, args[0]); err != nil {
		return err
	}

	p.Success("Deleted project %s", args[0])
	return nil
}
