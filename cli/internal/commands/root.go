package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	showbiz "github.com/showbiz-io/showbiz/sdk/go"
	"github.com/shaharby7/showbiz/cli/internal/config"
	"github.com/shaharby7/showbiz/cli/internal/output"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "showbiz",
	Short: "Showbiz CLI — manage your Showbiz resources",
	Long:  "The Showbiz command-line interface lets you manage organizations, projects, connections, resources, and IAM policies.",
	SilenceUsage:  true,
	SilenceErrors: true,
}

func init() {
	rootCmd.PersistentFlags().String("output", "table", "Output format: table or json")
	rootCmd.PersistentFlags().Bool("no-color", false, "Disable colored output")
	rootCmd.PersistentFlags().Bool("yes", false, "Skip confirmation prompts")

	rootCmd.AddCommand(authCmd)
	rootCmd.AddCommand(orgCmd)
	rootCmd.AddCommand(projectCmd)
	rootCmd.AddCommand(connectionCmd)
	rootCmd.AddCommand(resourceCmd)
	rootCmd.AddCommand(iamCmd)
	rootCmd.AddCommand(providerCmd)
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(completionCmd)
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// newClient creates an SDK client from config and stored credentials.
func newClient() (*showbiz.Client, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}
	apiURL := config.ResolveAPIURL(cfg)
	opts := []showbiz.Option{}
	creds, _ := config.LoadCredentials()
	if creds != nil && creds.AccessToken != "" {
		opts = append(opts, showbiz.WithToken(creds.AccessToken))
	}
	return showbiz.NewClient(apiURL, opts...), nil
}

// getPrinter creates an output.Printer from the command's flags.
func getPrinter(cmd *cobra.Command) *output.Printer {
	format, _ := cmd.Flags().GetString("output")
	noColor, _ := cmd.Flags().GetBool("no-color")
	f := output.FormatTable
	if format == "json" {
		f = output.FormatJSON
	}
	return output.NewPrinter(f, noColor)
}

// getOrg resolves the organization from the --org flag or config.
func getOrg(cmd *cobra.Command) (string, error) {
	org, _ := cmd.Flags().GetString("org")
	if org != "" {
		return org, nil
	}
	cfg, err := config.Load()
	if err != nil {
		return "", err
	}
	if cfg.Org != "" {
		return cfg.Org, nil
	}
	return "", fmt.Errorf("organization required: use --org flag or set with 'showbiz config set org <id>'")
}

// confirmAction prompts for confirmation unless --yes is set.
func confirmAction(cmd *cobra.Command, message string) bool {
	yes, _ := cmd.Flags().GetBool("yes")
	if yes {
		return true
	}
	fmt.Fprintf(os.Stderr, "%s [y/N]: ", message)
	var response string
	fmt.Scanln(&response)
	return strings.ToLower(response) == "y"
}

// parseJSONMap parses a JSON string into a map.
func parseJSONMap(s string) (map[string]interface{}, error) {
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(s), &m); err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}
	return m, nil
}

var completionCmd = &cobra.Command{
	Use:       "completion [bash|zsh|fish|powershell]",
	Short:     "Generate shell completion scripts",
	Long:      "Generate shell completion scripts for the specified shell.",
	Args:      cobra.ExactArgs(1),
	ValidArgs: []string{"bash", "zsh", "fish", "powershell"},
	RunE: func(cmd *cobra.Command, args []string) error {
		switch args[0] {
		case "bash":
			return rootCmd.GenBashCompletion(os.Stdout)
		case "zsh":
			return rootCmd.GenZshCompletion(os.Stdout)
		case "fish":
			return rootCmd.GenFishCompletion(os.Stdout, true)
		case "powershell":
			return rootCmd.GenPowerShellCompletionWithDesc(os.Stdout)
		}
		return fmt.Errorf("unsupported shell: %s", args[0])
	},
}
