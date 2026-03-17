package commands

import (
	"fmt"

	"github.com/shaharby7/showbiz/cli/internal/config"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage CLI configuration",
}

func init() {
	configSetCmd := &cobra.Command{
		Use:   "set <key> <value>",
		Short: "Set a configuration value",
		Args:  cobra.ExactArgs(2),
		RunE:  runConfigSet,
	}

	configGetCmd := &cobra.Command{
		Use:   "get <key>",
		Short: "Get a configuration value",
		Args:  cobra.ExactArgs(1),
		RunE:  runConfigGet,
	}

	configCmd.AddCommand(configSetCmd, configGetCmd)
}

func runConfigSet(cmd *cobra.Command, args []string) error {
	p := getPrinter(cmd)
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	if err := cfg.Set(args[0], args[1]); err != nil {
		return err
	}

	if err := cfg.Save(); err != nil {
		return err
	}

	p.Success("Set %s = %s", args[0], args[1])
	return nil
}

func runConfigGet(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	val, err := cfg.Get(args[0])
	if err != nil {
		return err
	}

	fmt.Println(val)
	return nil
}
