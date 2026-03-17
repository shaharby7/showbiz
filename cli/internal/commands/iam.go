package commands

import (
	"fmt"
	"strings"

	showbiz "github.com/showbiz-io/showbiz/sdk/go"
	"github.com/shaharby7/showbiz/cli/internal/output"
	"github.com/spf13/cobra"
)

var iamCmd = &cobra.Command{
	Use:   "iam",
	Short: "Identity and access management",
}

func init() {
	// policy subgroup
	policyCmd := &cobra.Command{
		Use:   "policy",
		Short: "Manage IAM policies",
	}

	policyListCmd := &cobra.Command{
		Use:   "list",
		Short: "List policies (global + organization)",
		RunE:  runPolicyList,
	}
	policyListCmd.Flags().String("org", "", "Organization ID")

	policyGetCmd := &cobra.Command{
		Use:   "get <policyID>",
		Short: "Get policy details",
		Args:  cobra.ExactArgs(1),
		RunE:  runPolicyGet,
	}

	policyCreateCmd := &cobra.Command{
		Use:   "create",
		Short: "Create an organization policy",
		RunE:  runPolicyCreate,
	}
	policyCreateCmd.Flags().String("org", "", "Organization ID")
	policyCreateCmd.Flags().String("name", "", "Policy name")
	policyCreateCmd.Flags().String("permissions", "", "Comma-separated permissions")
	_ = policyCreateCmd.MarkFlagRequired("name")
	_ = policyCreateCmd.MarkFlagRequired("permissions")

	policyUpdateCmd := &cobra.Command{
		Use:   "update <policyID>",
		Short: "Update an organization policy",
		Args:  cobra.ExactArgs(1),
		RunE:  runPolicyUpdate,
	}
	policyUpdateCmd.Flags().String("org", "", "Organization ID")
	policyUpdateCmd.Flags().String("permissions", "", "Comma-separated permissions")
	_ = policyUpdateCmd.MarkFlagRequired("permissions")

	policyDeleteCmd := &cobra.Command{
		Use:   "delete <policyID>",
		Short: "Delete an organization policy",
		Args:  cobra.ExactArgs(1),
		RunE:  runPolicyDelete,
	}
	policyDeleteCmd.Flags().String("org", "", "Organization ID")

	policyCmd.AddCommand(policyListCmd, policyGetCmd, policyCreateCmd, policyUpdateCmd, policyDeleteCmd)

	// attach
	attachCmd := &cobra.Command{
		Use:   "attach",
		Short: "Attach a policy to a user in a project",
		RunE:  runIAMAttach,
	}
	attachCmd.Flags().String("org", "", "Organization ID")
	attachCmd.Flags().String("project", "", "Project ID")
	attachCmd.Flags().String("user", "", "User email")
	attachCmd.Flags().String("policy", "", "Policy ID")
	_ = attachCmd.MarkFlagRequired("project")
	_ = attachCmd.MarkFlagRequired("user")
	_ = attachCmd.MarkFlagRequired("policy")

	// attachments
	attachmentsCmd := &cobra.Command{
		Use:   "attachments",
		Short: "List policy attachments for a project",
		RunE:  runIAMAttachments,
	}
	attachmentsCmd.Flags().String("org", "", "Organization ID")
	attachmentsCmd.Flags().String("project", "", "Project ID")
	_ = attachmentsCmd.MarkFlagRequired("project")

	// detach
	detachCmd := &cobra.Command{
		Use:   "detach",
		Short: "Detach a policy from a user in a project",
		RunE:  runIAMDetach,
	}
	detachCmd.Flags().String("org", "", "Organization ID")
	detachCmd.Flags().String("project", "", "Project ID")
	detachCmd.Flags().String("user", "", "User email")
	detachCmd.Flags().String("policy", "", "Policy ID")
	_ = detachCmd.MarkFlagRequired("project")
	_ = detachCmd.MarkFlagRequired("user")
	_ = detachCmd.MarkFlagRequired("policy")

	iamCmd.AddCommand(policyCmd, attachCmd, attachmentsCmd, detachCmd)
}

func runPolicyList(cmd *cobra.Command, args []string) error {
	p := getPrinter(cmd)
	orgID, err := getOrg(cmd)
	if err != nil {
		return err
	}

	client, err := newClient()
	if err != nil {
		return err
	}

	global, err := client.ListGlobalPolicies(cmd.Context())
	if err != nil {
		return err
	}

	org, err := client.ListOrgPolicies(cmd.Context(), orgID)
	if err != nil {
		return err
	}

	all := append(global, org...)

	t := p.Table("ID", "NAME", "SCOPE", "PERMISSIONS")
	for _, pol := range all {
		t.AddRow(pol.ID, pol.Name, pol.Scope, strings.Join(pol.Permissions, ", "))
	}
	return t.Print()
}

func runPolicyGet(cmd *cobra.Command, args []string) error {
	p := getPrinter(cmd)
	client, err := newClient()
	if err != nil {
		return err
	}

	pol, err := client.GetPolicy(cmd.Context(), args[0])
	if err != nil {
		return err
	}

	return p.PrintResource([]output.KeyValue{
		{Key: "ID", Value: pol.ID},
		{Key: "Name", Value: pol.Name},
		{Key: "Scope", Value: pol.Scope},
		{Key: "Organization", Value: pol.OrganizationID},
		{Key: "Permissions", Value: strings.Join(pol.Permissions, ", ")},
		{Key: "Created At", Value: pol.CreatedAt},
		{Key: "Updated At", Value: pol.UpdatedAt},
	})
}

func runPolicyCreate(cmd *cobra.Command, args []string) error {
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
	permsStr, _ := cmd.Flags().GetString("permissions")
	perms := splitPermissions(permsStr)

	pol, err := client.CreateOrgPolicy(cmd.Context(), orgID, showbiz.CreatePolicyInput{
		Name:        name,
		Permissions: perms,
	})
	if err != nil {
		return err
	}

	p.Success("Created policy %s (ID: %s)", pol.Name, pol.ID)
	return nil
}

func runPolicyUpdate(cmd *cobra.Command, args []string) error {
	p := getPrinter(cmd)
	orgID, err := getOrg(cmd)
	if err != nil {
		return err
	}

	client, err := newClient()
	if err != nil {
		return err
	}

	permsStr, _ := cmd.Flags().GetString("permissions")
	perms := splitPermissions(permsStr)

	pol, err := client.UpdateOrgPolicy(cmd.Context(), orgID, args[0], showbiz.UpdatePolicyInput{
		Permissions: perms,
	})
	if err != nil {
		return err
	}

	p.Success("Updated policy %s", pol.ID)
	return nil
}

func runPolicyDelete(cmd *cobra.Command, args []string) error {
	p := getPrinter(cmd)
	orgID, err := getOrg(cmd)
	if err != nil {
		return err
	}

	if !confirmAction(cmd, fmt.Sprintf("Delete policy %s?", args[0])) {
		return nil
	}

	client, err := newClient()
	if err != nil {
		return err
	}

	if err := client.DeleteOrgPolicy(cmd.Context(), orgID, args[0]); err != nil {
		return err
	}

	p.Success("Deleted policy %s", args[0])
	return nil
}

func runIAMAttach(cmd *cobra.Command, args []string) error {
	p := getPrinter(cmd)
	orgID, err := getOrg(cmd)
	if err != nil {
		return err
	}

	projectID, _ := cmd.Flags().GetString("project")
	userEmail, _ := cmd.Flags().GetString("user")
	policyID, _ := cmd.Flags().GetString("policy")

	client, err := newClient()
	if err != nil {
		return err
	}

	attachment, err := client.AttachPolicy(cmd.Context(), orgID, projectID, showbiz.AttachPolicyInput{
		UserEmail: userEmail,
		PolicyID:  policyID,
	})
	if err != nil {
		return err
	}

	p.Success("Attached policy %s to user %s (ID: %s)", policyID, userEmail, attachment.ID)
	return nil
}

func runIAMAttachments(cmd *cobra.Command, args []string) error {
	p := getPrinter(cmd)
	orgID, err := getOrg(cmd)
	if err != nil {
		return err
	}

	projectID, _ := cmd.Flags().GetString("project")

	client, err := newClient()
	if err != nil {
		return err
	}

	attachments, err := client.ListPolicyAttachments(cmd.Context(), orgID, projectID)
	if err != nil {
		return err
	}

	t := p.Table("ID", "USER", "POLICY")
	for _, a := range attachments {
		t.AddRow(a.ID, a.UserEmail, a.PolicyID)
	}
	return t.Print()
}

func runIAMDetach(cmd *cobra.Command, args []string) error {
	p := getPrinter(cmd)
	orgID, err := getOrg(cmd)
	if err != nil {
		return err
	}

	projectID, _ := cmd.Flags().GetString("project")
	userEmail, _ := cmd.Flags().GetString("user")
	policyID, _ := cmd.Flags().GetString("policy")

	client, err := newClient()
	if err != nil {
		return err
	}

	if err := client.DetachPolicy(cmd.Context(), orgID, projectID, showbiz.DetachPolicyInput{
		UserEmail: userEmail,
		PolicyID:  policyID,
	}); err != nil {
		return err
	}

	p.Success("Detached policy %s from user %s", policyID, userEmail)
	return nil
}

func splitPermissions(s string) []string {
	parts := strings.Split(s, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		trimmed := strings.TrimSpace(p)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
