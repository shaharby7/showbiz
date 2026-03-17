package commands

import (
	"strings"

	"github.com/shaharby7/showbiz/cli/internal/output"
	"github.com/spf13/cobra"
)

var providerCmd = &cobra.Command{
	Use:   "provider",
	Short: "Browse available providers",
}

func init() {
	providerListCmd := &cobra.Command{
		Use:   "list",
		Short: "List available providers",
		RunE:  runProviderList,
	}

	providerGetCmd := &cobra.Command{
		Use:   "get <name>",
		Short: "Get provider details",
		Args:  cobra.ExactArgs(1),
		RunE:  runProviderGet,
	}

	providerCmd.AddCommand(providerListCmd, providerGetCmd)
}

func runProviderList(cmd *cobra.Command, args []string) error {
	p := getPrinter(cmd)
	client, err := newClient()
	if err != nil {
		return err
	}

	providers, err := client.ListProviders(cmd.Context())
	if err != nil {
		return err
	}

	t := p.Table("NAME", "RESOURCE TYPES")
	for _, prov := range providers {
		t.AddRow(prov.Name, strings.Join(prov.ResourceTypes, ", "))
	}
	return t.Print()
}

func runProviderGet(cmd *cobra.Command, args []string) error {
	p := getPrinter(cmd)
	client, err := newClient()
	if err != nil {
		return err
	}

	prov, err := client.GetProvider(cmd.Context(), args[0])
	if err != nil {
		return err
	}

	return p.PrintResource([]output.KeyValue{
		{Key: "Name", Value: prov.Name},
		{Key: "Resource Types", Value: strings.Join(prov.ResourceTypes, ", ")},
	})
}
