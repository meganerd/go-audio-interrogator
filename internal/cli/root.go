package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	jsonOutput bool
	verbose    bool
	showAll    bool
	cardFilter string
	devFilter  string
	listCards  bool
	noProc     bool
)

func newRootCmd(version string) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "go-audio-interrogator",
		Short:   "Interrogate audio devices for their capabilities",
		Version: version,
		RunE:    run,
	}

	cmd.Flags().BoolVarP(&jsonOutput, "json", "j", false, "Output results in JSON format")
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
	cmd.Flags().BoolVarP(&showAll, "all", "a", false, "Show all devices including duplicates and virtual devices")
	cmd.Flags().StringVarP(&cardFilter, "card", "c", "", "Filter by specific card ID")
	cmd.Flags().StringVarP(&devFilter, "device", "d", "", "Filter by device name (partial match, case-insensitive)")
	cmd.Flags().BoolVarP(&listCards, "list", "l", false, "List available card IDs and exit")
	cmd.Flags().BoolVar(&noProc, "no-proc", false, "Disable /proc/asound access (Linux)")

	return cmd
}

func Execute(version string) error {
	return newRootCmd(version).Execute()
}

func run(cmd *cobra.Command, args []string) error {
	fmt.Println("go-audio-interrogator - not yet implemented")
	return nil
}
