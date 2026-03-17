package commands

import (
	"fmt"

	showbiz "github.com/showbiz-io/showbiz/sdk/go"
	"github.com/shaharby7/showbiz/cli/internal/output"
	"github.com/spf13/cobra"
)

var resourceCmd = &cobra.Command{
	Use:   "resource",
	Short: "Manage resources",
}

func init() {
	// resource list
	resListCmd := &cobra.Command{
		Use:   "list",
		Short: "List resources",
		RunE:  runResourceList,
	}
	resListCmd.Flags().String("project", "", "Project ID")
	_ = resListCmd.MarkFlagRequired("project")

	// resource create
	resCreateCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a resource",
		RunE:  runResourceCreate,
	}
	resCreateCmd.Flags().String("project", "", "Project ID")
	resCreateCmd.Flags().String("connection", "", "Connection ID")
	resCreateCmd.Flags().String("type", "", "Resource type")
	resCreateCmd.Flags().String("name", "", "Resource name")
	resCreateCmd.Flags().String("values", "{}", "Values JSON")
	_ = resCreateCmd.MarkFlagRequired("project")
	_ = resCreateCmd.MarkFlagRequired("connection")
	_ = resCreateCmd.MarkFlagRequired("type")
	_ = resCreateCmd.MarkFlagRequired("name")

	// resource get
	resGetCmd := &cobra.Command{
		Use:   "get <resourceID>",
		Short: "Get resource details",
		Args:  cobra.ExactArgs(1),
		RunE:  runResourceGet,
	}
	resGetCmd.Flags().String("project", "", "Project ID")
	_ = resGetCmd.MarkFlagRequired("project")

	// resource update
	resUpdateCmd := &cobra.Command{
		Use:   "update <resourceID>",
		Short: "Update a resource",
		Args:  cobra.ExactArgs(1),
		RunE:  runResourceUpdate,
	}
	resUpdateCmd.Flags().String("project", "", "Project ID")
	resUpdateCmd.Flags().String("values", "", "Values JSON")
	_ = resUpdateCmd.MarkFlagRequired("project")
	_ = resUpdateCmd.MarkFlagRequired("values")

	// resource delete
	resDeleteCmd := &cobra.Command{
		Use:   "delete <resourceID>",
		Short: "Delete a resource",
		Args:  cobra.ExactArgs(1),
		RunE:  runResourceDelete,
	}
	resDeleteCmd.Flags().String("project", "", "Project ID")
	_ = resDeleteCmd.MarkFlagRequired("project")

	resourceCmd.AddCommand(resListCmd, resCreateCmd, resGetCmd, resUpdateCmd, resDeleteCmd)
}

func runResourceList(cmd *cobra.Command, args []string) error {
	p := getPrinter(cmd)
	projectID, _ := cmd.Flags().GetString("project")

	client, err := newClient()
	if err != nil {
		return err
	}

	result, err := client.ListResources(cmd.Context(), projectID, nil)
	if err != nil {
		return err
	}

	t := p.Table("ID", "NAME", "TYPE", "STATUS")
	for _, r := range result.Data {
		t.AddRow(r.ID, r.Name, r.ResourceType, r.Status)
	}
	return t.Print()
}

func runResourceCreate(cmd *cobra.Command, args []string) error {
	p := getPrinter(cmd)
	projectID, _ := cmd.Flags().GetString("project")
	connID, _ := cmd.Flags().GetString("connection")
	resType, _ := cmd.Flags().GetString("type")
	name, _ := cmd.Flags().GetString("name")
	valuesStr, _ := cmd.Flags().GetString("values")

	values, err := parseJSONMap(valuesStr)
	if err != nil {
		return fmt.Errorf("--values: %w", err)
	}

	client, err := newClient()
	if err != nil {
		return err
	}

	res, err := client.CreateResource(cmd.Context(), projectID, showbiz.CreateResourceInput{
		Name:         name,
		ConnectionID: connID,
		ResourceType: resType,
		Values:       values,
	})
	if err != nil {
		return err
	}

	p.Success("Created resource %s (ID: %s)", res.Name, res.ID)
	return nil
}

func runResourceGet(cmd *cobra.Command, args []string) error {
	p := getPrinter(cmd)
	projectID, _ := cmd.Flags().GetString("project")

	client, err := newClient()
	if err != nil {
		return err
	}

	res, err := client.GetResource(cmd.Context(), projectID, args[0])
	if err != nil {
		return err
	}

	return p.PrintResource([]output.KeyValue{
		{Key: "ID", Value: res.ID},
		{Key: "Name", Value: res.Name},
		{Key: "Project", Value: res.ProjectID},
		{Key: "Connection", Value: res.ConnectionID},
		{Key: "Type", Value: res.ResourceType},
		{Key: "Status", Value: res.Status},
		{Key: "Values", Value: res.Values},
		{Key: "Created At", Value: res.CreatedAt},
		{Key: "Updated At", Value: res.UpdatedAt},
	})
}

func runResourceUpdate(cmd *cobra.Command, args []string) error {
	p := getPrinter(cmd)
	projectID, _ := cmd.Flags().GetString("project")
	valuesStr, _ := cmd.Flags().GetString("values")

	values, err := parseJSONMap(valuesStr)
	if err != nil {
		return fmt.Errorf("--values: %w", err)
	}

	client, err := newClient()
	if err != nil {
		return err
	}

	res, err := client.UpdateResource(cmd.Context(), projectID, args[0], showbiz.UpdateResourceInput{
		Values: values,
	})
	if err != nil {
		return err
	}

	p.Success("Updated resource %s", res.ID)
	return nil
}

func runResourceDelete(cmd *cobra.Command, args []string) error {
	p := getPrinter(cmd)
	projectID, _ := cmd.Flags().GetString("project")

	if !confirmAction(cmd, fmt.Sprintf("Delete resource %s?", args[0])) {
		return nil
	}

	client, err := newClient()
	if err != nil {
		return err
	}

	if err := client.DeleteResource(cmd.Context(), projectID, args[0]); err != nil {
		return err
	}

	p.Success("Deleted resource %s", args[0])
	return nil
}
