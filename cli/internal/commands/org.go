package commands

import (
	"fmt"

	showbiz "github.com/showbiz-io/showbiz/sdk/go"
	"github.com/shaharby7/showbiz/cli/internal/output"
	"github.com/spf13/cobra"
)

var orgCmd = &cobra.Command{
	Use:   "org",
	Short: "Manage organizations",
}

func init() {
	// org list
	orgListCmd := &cobra.Command{
		Use:   "list",
		Short: "List organizations",
		RunE:  runOrgList,
	}

	// org create
	orgCreateCmd := &cobra.Command{
		Use:   "create",
		Short: "Create an organization",
		RunE:  runOrgCreate,
	}
	orgCreateCmd.Flags().String("name", "", "Organization name")
	orgCreateCmd.Flags().String("display-name", "", "Display name")
	_ = orgCreateCmd.MarkFlagRequired("name")

	// org get
	orgGetCmd := &cobra.Command{
		Use:   "get <id>",
		Short: "Get organization details",
		Args:  cobra.ExactArgs(1),
		RunE:  runOrgGet,
	}

	// org update
	orgUpdateCmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update an organization",
		Args:  cobra.ExactArgs(1),
		RunE:  runOrgUpdate,
	}
	orgUpdateCmd.Flags().String("display-name", "", "New display name")
	_ = orgUpdateCmd.MarkFlagRequired("display-name")

	// org deactivate
	orgDeactivateCmd := &cobra.Command{
		Use:   "deactivate <id>",
		Short: "Deactivate an organization",
		Args:  cobra.ExactArgs(1),
		RunE:  runOrgDeactivate,
	}

	// org activate
	orgActivateCmd := &cobra.Command{
		Use:   "activate <id>",
		Short: "Activate an organization",
		Args:  cobra.ExactArgs(1),
		RunE:  runOrgActivate,
	}

	// members subgroup
	membersCmd := &cobra.Command{
		Use:   "members",
		Short: "Manage organization members",
	}

	membersListCmd := &cobra.Command{
		Use:   "list <orgID>",
		Short: "List organization members",
		Args:  cobra.ExactArgs(1),
		RunE:  runMembersList,
	}

	membersAddCmd := &cobra.Command{
		Use:   "add <orgID>",
		Short: "Add a member to the organization",
		Args:  cobra.ExactArgs(1),
		RunE:  runMembersAdd,
	}
	membersAddCmd.Flags().String("email", "", "Member email")
	_ = membersAddCmd.MarkFlagRequired("email")

	membersRemoveCmd := &cobra.Command{
		Use:   "remove <orgID>",
		Short: "Remove a member from the organization",
		Args:  cobra.ExactArgs(1),
		RunE:  runMembersRemove,
	}
	membersRemoveCmd.Flags().String("email", "", "Member email")
	_ = membersRemoveCmd.MarkFlagRequired("email")

	membersCmd.AddCommand(membersListCmd, membersAddCmd, membersRemoveCmd)
	orgCmd.AddCommand(orgListCmd, orgCreateCmd, orgGetCmd, orgUpdateCmd, orgDeactivateCmd, orgActivateCmd, membersCmd)
}

func runOrgList(cmd *cobra.Command, args []string) error {
	p := getPrinter(cmd)
	client, err := newClient()
	if err != nil {
		return err
	}

	result, err := client.ListOrganizations(cmd.Context(), nil)
	if err != nil {
		return err
	}

	t := p.Table("ID", "NAME", "DISPLAY NAME", "ACTIVE")
	for _, o := range result.Data {
		t.AddRow(o.ID, o.Name, o.DisplayName, fmt.Sprintf("%t", o.Active))
	}
	return t.Print()
}

func runOrgCreate(cmd *cobra.Command, args []string) error {
	p := getPrinter(cmd)
	client, err := newClient()
	if err != nil {
		return err
	}

	name, _ := cmd.Flags().GetString("name")
	displayName, _ := cmd.Flags().GetString("display-name")

	org, err := client.CreateOrganization(cmd.Context(), showbiz.CreateOrganizationInput{
		Name:        name,
		DisplayName: displayName,
	})
	if err != nil {
		return err
	}

	p.Success("Created organization %s (ID: %s)", org.Name, org.ID)
	return nil
}

func runOrgGet(cmd *cobra.Command, args []string) error {
	p := getPrinter(cmd)
	client, err := newClient()
	if err != nil {
		return err
	}

	org, err := client.GetOrganization(cmd.Context(), args[0])
	if err != nil {
		return err
	}

	return p.PrintResource([]output.KeyValue{
		{Key: "ID", Value: org.ID},
		{Key: "Name", Value: org.Name},
		{Key: "Display Name", Value: org.DisplayName},
		{Key: "Active", Value: org.Active},
		{Key: "Created At", Value: org.CreatedAt},
		{Key: "Updated At", Value: org.UpdatedAt},
	})
}

func runOrgUpdate(cmd *cobra.Command, args []string) error {
	p := getPrinter(cmd)
	client, err := newClient()
	if err != nil {
		return err
	}

	displayName, _ := cmd.Flags().GetString("display-name")
	org, err := client.UpdateOrganization(cmd.Context(), args[0], showbiz.UpdateOrganizationInput{
		DisplayName: displayName,
	})
	if err != nil {
		return err
	}

	p.Success("Updated organization %s", org.ID)
	return nil
}

func runOrgDeactivate(cmd *cobra.Command, args []string) error {
	p := getPrinter(cmd)
	if !confirmAction(cmd, fmt.Sprintf("Deactivate organization %s?", args[0])) {
		return nil
	}

	client, err := newClient()
	if err != nil {
		return err
	}

	if err := client.DeactivateOrganization(cmd.Context(), args[0]); err != nil {
		return err
	}

	p.Success("Deactivated organization %s", args[0])
	return nil
}

func runOrgActivate(cmd *cobra.Command, args []string) error {
	p := getPrinter(cmd)
	client, err := newClient()
	if err != nil {
		return err
	}

	if err := client.ActivateOrganization(cmd.Context(), args[0]); err != nil {
		return err
	}

	p.Success("Activated organization %s", args[0])
	return nil
}

func runMembersList(cmd *cobra.Command, args []string) error {
	p := getPrinter(cmd)
	client, err := newClient()
	if err != nil {
		return err
	}

	members, err := client.ListMembers(cmd.Context(), args[0])
	if err != nil {
		return err
	}

	t := p.Table("EMAIL", "DISPLAY NAME", "ACTIVE")
	for _, m := range members {
		t.AddRow(m.Email, m.DisplayName, fmt.Sprintf("%t", m.Active))
	}
	return t.Print()
}

func runMembersAdd(cmd *cobra.Command, args []string) error {
	p := getPrinter(cmd)
	client, err := newClient()
	if err != nil {
		return err
	}

	email, _ := cmd.Flags().GetString("email")
	if err := client.AddMember(cmd.Context(), args[0], email); err != nil {
		return err
	}

	p.Success("Added %s to organization %s", email, args[0])
	return nil
}

func runMembersRemove(cmd *cobra.Command, args []string) error {
	p := getPrinter(cmd)
	client, err := newClient()
	if err != nil {
		return err
	}

	email, _ := cmd.Flags().GetString("email")
	if err := client.RemoveMember(cmd.Context(), args[0], email); err != nil {
		return err
	}

	p.Success("Removed %s from organization %s", email, args[0])
	return nil
}

