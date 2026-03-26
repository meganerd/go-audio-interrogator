package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/meganerd/go-audio-interrogator/internal/audio"
	"github.com/meganerd/go-audio-interrogator/internal/backend"
	"github.com/meganerd/go-audio-interrogator/internal/output"
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
	be := backend.NewPlatformBackend()

	// Handle --list mode
	if listCards {
		cards, err := be.ListCards()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Failed to list cards: %v\n", err)
			return nil
		}
		output.PrintCardList(cards)
		return nil
	}

	if verbose && !jsonOutput {
		fmt.Println("🎵 Audio Interrogator - Scanning system audio devices...")
	}

	// Enumerate devices
	opts := audio.EnumerateOpts{NoProc: noProc}
	devices, err := be.EnumerateDevices(opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Failed to enumerate devices: %v\n", err)
	}

	// Apply filters
	devices = audio.FilterDevices(devices, cardFilter, devFilter, showAll)

	// Build system info
	info := audio.NewSystemAudioInfo(devices)

	// Detect defaults
	for _, d := range devices {
		if d.HasInput() && info.DefaultInput == "" {
			info.DefaultInput = d.Name
		}
		if d.HasOutput() && info.DefaultOutput == "" {
			info.DefaultOutput = d.Name
		}
	}

	if jsonOutput {
		return output.PrintJSON(info)
	}

	// Text output
	cards, err := be.ListCards()
	if err == nil {
		output.PrintCardSummary(cards)
	}
	output.PrintSystemSummary(info)
	output.PrintDeviceList(info.Devices, verbose)

	return nil
}
