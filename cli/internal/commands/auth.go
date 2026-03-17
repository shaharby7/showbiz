package commands

import (
	"os"

	showbiz "github.com/showbiz-io/showbiz/sdk/go"
	"github.com/shaharby7/showbiz/cli/internal/config"
	"github.com/shaharby7/showbiz/cli/internal/output"
	"github.com/spf13/cobra"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authentication commands",
}

func init() {
	// login
	loginCmd := &cobra.Command{
		Use:   "login",
		Short: "Log in to Showbiz",
		RunE:  runAuthLogin,
	}
	loginCmd.Flags().String("username", os.Getenv("SHOWBIZ_USERNAME"), "Email address")
	loginCmd.Flags().String("password", os.Getenv("SHOWBIZ_PASSWORD"), "Password")
	_ = loginCmd.MarkFlagRequired("username")
	_ = loginCmd.MarkFlagRequired("password")

	// register
	registerCmd := &cobra.Command{
		Use:   "register",
		Short: "Register a new account",
		RunE:  runAuthRegister,
	}
	registerCmd.Flags().String("username", "", "Email address")
	registerCmd.Flags().String("password", "", "Password")
	registerCmd.Flags().String("display-name", "", "Display name")
	registerCmd.Flags().String("org", "", "Organization ID")
	_ = registerCmd.MarkFlagRequired("username")
	_ = registerCmd.MarkFlagRequired("password")

	// logout
	logoutCmd := &cobra.Command{
		Use:   "logout",
		Short: "Log out and clear stored credentials",
		RunE:  runAuthLogout,
	}

	// status
	statusCmd := &cobra.Command{
		Use:   "status",
		Short: "Show current authentication status",
		RunE:  runAuthStatus,
	}

	authCmd.AddCommand(loginCmd, registerCmd, logoutCmd, statusCmd)
}

func runAuthLogin(cmd *cobra.Command, args []string) error {
	p := getPrinter(cmd)
	username, _ := cmd.Flags().GetString("username")
	password, _ := cmd.Flags().GetString("password")

	cfg, err := config.Load()
	if err != nil {
		return err
	}
	apiURL := config.ResolveAPIURL(cfg)
	client := showbiz.NewClient(apiURL)

	resp, err := client.Login(cmd.Context(), showbiz.LoginInput{
		Email:    username,
		Password: password,
	})
	if err != nil {
		return err
	}

	if err := config.SaveCredentials(&config.Credentials{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
	}); err != nil {
		return err
	}

	p.Success("Logged in as %s", username)
	return nil
}

func runAuthRegister(cmd *cobra.Command, args []string) error {
	p := getPrinter(cmd)
	username, _ := cmd.Flags().GetString("username")
	password, _ := cmd.Flags().GetString("password")
	displayName, _ := cmd.Flags().GetString("display-name")
	orgID, _ := cmd.Flags().GetString("org")

	cfg, err := config.Load()
	if err != nil {
		return err
	}
	apiURL := config.ResolveAPIURL(cfg)
	client := showbiz.NewClient(apiURL)

	user, err := client.Register(cmd.Context(), showbiz.RegisterInput{
		Email:          username,
		Password:       password,
		DisplayName:    displayName,
		OrganizationID: orgID,
	})
	if err != nil {
		return err
	}

	p.Success("Registered account %s", user.Email)
	return nil
}

func runAuthLogout(cmd *cobra.Command, args []string) error {
	p := getPrinter(cmd)
	if err := config.ClearCredentials(); err != nil {
		return err
	}
	p.Success("Logged out")
	return nil
}

func runAuthStatus(cmd *cobra.Command, args []string) error {
	p := getPrinter(cmd)
	creds, err := config.LoadCredentials()
	if err != nil {
		return err
	}
	if creds == nil || creds.AccessToken == "" {
		p.PrintResource([]output.KeyValue{
			{Key: "Status", Value: "Not logged in"},
		})
		return nil
	}

	client, err := newClient()
	if err != nil {
		return err
	}

	user, err := client.Me(cmd.Context())
	if err != nil {
		p.PrintResource([]output.KeyValue{
			{Key: "Status", Value: "Logged in (unable to fetch user details)"},
		})
		return nil
	}

	return p.PrintResource([]output.KeyValue{
		{Key: "Status", Value: "Logged in"},
		{Key: "Email", Value: user.Email},
		{Key: "Display Name", Value: user.DisplayName},
		{Key: "Organization", Value: user.OrganizationID},
	})
}
