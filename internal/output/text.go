package output

import (
	"fmt"
	"strings"

	"github.com/meganerd/go-audio-interrogator/internal/audio"
)

func PrintCardSummary(cards []audio.CardInfo) {
	fmt.Println("════════════════════════════════════════")
	fmt.Println("        AVAILABLE AUDIO CARDS")
	fmt.Println("════════════════════════════════════════")
	if len(cards) == 0 {
		fmt.Println("  (No audio cards found)")
		return
	}
	for _, c := range cards {
		fmt.Printf("  %s [%s]: %s\n", c.ID, c.ShortName, c.Description)
	}
}

func PrintSystemSummary(info audio.SystemAudioInfo) {
	fmt.Println("\n════════════════════════════════════════")
	fmt.Println("         SYSTEM AUDIO SUMMARY")
	fmt.Println("════════════════════════════════════════")
	fmt.Printf("Total Devices Found: %d\n", len(info.Devices))
	fmt.Printf("Input Devices: %d\n", info.TotalInputDevices)
	fmt.Printf("Output Devices: %d\n", info.TotalOutputDevices)

	if info.DefaultInput != "" {
		fmt.Printf("Default Input: %s\n", info.DefaultInput)
	}
	if info.DefaultOutput != "" {
		fmt.Printf("Default Output: %s\n", info.DefaultOutput)
	}
}

func PrintDeviceList(devices []audio.AudioDeviceInfo, verbose bool) {
	fmt.Println("\n════════════════════════════════════════")
	fmt.Println("           DEVICE DETAILS")
	fmt.Println("════════════════════════════════════════\n")

	for i, d := range devices {
		if verbose {
			PrintDeviceVerbose(i+1, d)
		} else {
			PrintDeviceCompact(i+1, d)
		}
	}

	if !verbose {
		fmt.Println("\nUse --verbose flag for detailed device information")
		fmt.Println("Use --card <id> to filter by card, --device <name> to filter by name")
		fmt.Println("Use --all to show all devices including duplicates")
	}
}

func PrintDeviceVerbose(num int, d audio.AudioDeviceInfo) {
	fmt.Printf("Device #%d\n", num)
	fmt.Printf("┌─ Device: %s\n", d.Name)
	fmt.Printf("├─ Type: %s\n", d.DeviceType)
	fmt.Printf("├─ Driver: %s\n", d.Driver)
	fmt.Printf("├─ Input Channels: %d\n", d.InputChannels)
	fmt.Printf("├─ Output Channels: %d\n", d.OutputChannels)
	fmt.Printf("├─ Default Sample Rate: %d Hz\n", d.DefaultSampleRate)
	fmt.Printf("├─ Default Buffer Size: %d samples\n", d.DefaultBufferSize)

	if len(d.SupportedSampleRates) > 0 {
		rates := make([]string, len(d.SupportedSampleRates))
		for i, r := range d.SupportedSampleRates {
			rates[i] = fmt.Sprintf("%d", r)
		}
		fmt.Printf("├─ Supported Sample Rates: [%s] Hz\n", strings.Join(rates, ", "))
	}

	bufSizes := make([]string, len(d.SupportedBufferSizes))
	for i, b := range d.SupportedBufferSizes {
		bufSizes[i] = fmt.Sprintf("%d", b)
	}
	fmt.Printf("└─ Supported Buffer Sizes: [%s] samples\n", strings.Join(bufSizes, ", "))
	fmt.Println()
}

func PrintDeviceCompact(num int, d audio.AudioDeviceInfo) {
	fmt.Printf("%d: %s (%s) - In: %d, Out: %d, SR: %dHz\n",
		num, d.Name, d.Driver, d.InputChannels, d.OutputChannels, d.DefaultSampleRate)
}

func PrintCardList(cards []audio.CardInfo) {
	fmt.Println("Available Audio Cards:")
	fmt.Println("═════════════════════")
	for _, c := range cards {
		fmt.Printf("  %s [%s]: %s\n", c.ID, c.ShortName, c.Description)
	}
	fmt.Println("\nUsage examples:")
	fmt.Println("  go-audio-interrogator --card 0        # Show devices for card0")
	fmt.Println("  go-audio-interrogator --card card1    # Show devices for card1")
	fmt.Println("  go-audio-interrogator --device Audio  # Show devices matching 'Audio'")
}
