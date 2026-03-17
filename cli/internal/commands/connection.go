package commands

import (
	"fmt"

	showbiz "github.com/showbiz-io/showbiz/sdk/go"
	"github.com/shaharby7/showbiz/cli/internal/output"
	"github.com/spf13/cobra"
)

var connectionCmd = &cobra.Command{
	Use:     "connection",
	Aliases: []string{"conn"},
	Short:   "Manage connections",
}

func init() {
	// connection list
	connListCmd := &cobra.Command{
		Use:   "list",
		Short: "List connections",
		RunE:  runConnectionList,
	}
	connListCmd.Flags().String("project", "", "Project ID")
	_ = connListCmd.MarkFlagRequired("project")

	// connection create
	connCreateCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a connection",
		RunE:  runConnectionCreate,
	}
	connCreateCmd.Flags().String("project", "", "Project ID")
	connCreateCmd.Flags().String("name", "", "Connection name")
	connCreateCmd.Flags().String("provider", "", "Provider name")
	connCreateCmd.Flags().String("credentials", "{}", "Credentials JSON")
	connCreateCmd.Flags().String("config", "{}", "Config JSON")
	_ = connCreateCmd.MarkFlagRequired("project")
	_ = connCreateCmd.MarkFlagRequired("name")
	_ = connCreateCmd.MarkFlagRequired("provider")

	// connection get
	connGetCmd := &cobra.Command{
		Use:   "get <connectionID>",
		Short: "Get connection details",
		Args:  cobra.ExactArgs(1),
		RunE:  runConnectionGet,
	}
	connGetCmd.Flags().String("project", "", "Project ID")
	_ = connGetCmd.MarkFlagRequired("project")

	// connection update
	connUpdateCmd := &cobra.Command{
		Use:   "update <connectionID>",
		Short: "Update a connection",
		Args:  cobra.ExactArgs(1),
		RunE:  runConnectionUpdate,
	}
	connUpdateCmd.Flags().String("project", "", "Project ID")
	connUpdateCmd.Flags().String("config", "", "Config JSON")
	_ = connUpdateCmd.MarkFlagRequired("project")
	_ = connUpdateCmd.MarkFlagRequired("config")

	// connection delete
	connDeleteCmd := &cobra.Command{
		Use:   "delete <connectionID>",
		Short: "Delete a connection",
		Args:  cobra.ExactArgs(1),
		RunE:  runConnectionDelete,
	}
	connDeleteCmd.Flags().String("project", "", "Project ID")
	_ = connDeleteCmd.MarkFlagRequired("project")

	connectionCmd.AddCommand(connListCmd, connCreateCmd, connGetCmd, connUpdateCmd, connDeleteCmd)
}

func runConnectionList(cmd *cobra.Command, args []string) error {
	p := getPrinter(cmd)
	projectID, _ := cmd.Flags().GetString("project")

	client, err := newClient()
	if err != nil {
		return err
	}

	result, err := client.ListConnections(cmd.Context(), projectID, nil)
	if err != nil {
		return err
	}

	t := p.Table("ID", "NAME", "PROVIDER")
	for _, c := range result.Data {
		t.AddRow(c.ID, c.Name, c.Provider)
	}
	return t.Print()
}

func runConnectionCreate(cmd *cobra.Command, args []string) error {
	p := getPrinter(cmd)
	projectID, _ := cmd.Flags().GetString("project")
	name, _ := cmd.Flags().GetString("name")
	provider, _ := cmd.Flags().GetString("provider")
	credsStr, _ := cmd.Flags().GetString("credentials")
	configStr, _ := cmd.Flags().GetString("config")

	creds, err := parseJSONMap(credsStr)
	if err != nil {
		return fmt.Errorf("--credentials: %w", err)
	}
	cfg, err := parseJSONMap(configStr)
	if err != nil {
		return fmt.Errorf("--config: %w", err)
	}

	client, err := newClient()
	if err != nil {
		return err
	}

	conn, err := client.CreateConnection(cmd.Context(), projectID, showbiz.CreateConnectionInput{
		Name:        name,
		Provider:    provider,
		Credentials: creds,
		Config:      cfg,
	})
	if err != nil {
		return err
	}

	p.Success("Created connection %s (ID: %s)", conn.Name, conn.ID)
	return nil
}

func runConnectionGet(cmd *cobra.Command, args []string) error {
	p := getPrinter(cmd)
	projectID, _ := cmd.Flags().GetString("project")

	client, err := newClient()
	if err != nil {
		return err
	}

	conn, err := client.GetConnection(cmd.Context(), projectID, args[0])
	if err != nil {
		return err
	}

	return p.PrintResource([]output.KeyValue{
		{Key: "ID", Value: conn.ID},
		{Key: "Name", Value: conn.Name},
		{Key: "Project", Value: conn.ProjectID},
		{Key: "Provider", Value: conn.Provider},
		{Key: "Config", Value: conn.Config},
		{Key: "Created At", Value: conn.CreatedAt},
		{Key: "Updated At", Value: conn.UpdatedAt},
	})
}

func runConnectionUpdate(cmd *cobra.Command, args []string) error {
	p := getPrinter(cmd)
	projectID, _ := cmd.Flags().GetString("project")
	configStr, _ := cmd.Flags().GetString("config")

	cfg, err := parseJSONMap(configStr)
	if err != nil {
		return fmt.Errorf("--config: %w", err)
	}

	client, err := newClient()
	if err != nil {
		return err
	}

	conn, err := client.UpdateConnection(cmd.Context(), projectID, args[0], showbiz.UpdateConnectionInput{
		Config: cfg,
	})
	if err != nil {
		return err
	}

	p.Success("Updated connection %s", conn.ID)
	return nil
}

func runConnectionDelete(cmd *cobra.Command, args []string) error {
	p := getPrinter(cmd)
	projectID, _ := cmd.Flags().GetString("project")

	if !confirmAction(cmd, fmt.Sprintf("Delete connection %s?", args[0])) {
		return nil
	}

	client, err := newClient()
	if err != nil {
		return err
	}

	if err := client.DeleteConnection(cmd.Context(), projectID, args[0]); err != nil {
		return err
	}

	p.Success("Deleted connection %s", args[0])
	return nil
}
