//go:build linux

package backend

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/meganerd/go-audio-interrogator/internal/audio"
)

type linuxBackend struct{}

func NewPlatformBackend() audio.Backend {
	return &linuxBackend{}
}

func (b *linuxBackend) Name() string {
	return "ALSA"
}

func (b *linuxBackend) ListCards() ([]audio.CardInfo, error) {
	return parseCards()
}

func (b *linuxBackend) EnumerateDevices(opts audio.EnumerateOpts) ([]audio.AudioDeviceInfo, error) {
	var devices []audio.AudioDeviceInfo

	if !opts.NoProc {
		procDevices, err := enumerateFromProc()
		if err == nil {
			devices = append(devices, procDevices...)
		}
	}

	return devices, nil
}

var cardLineRe = regexp.MustCompile(`^\s*(\d+)\s+\[(\w+)\s*\]:\s+\S+\s+-\s+(.+)$`)

func parseCards() ([]audio.CardInfo, error) {
	data, err := os.ReadFile("/proc/asound/cards")
	if err != nil {
		return nil, fmt.Errorf("reading /proc/asound/cards: %w", err)
	}

	var cards []audio.CardInfo
	for _, line := range strings.Split(string(data), "\n") {
		m := cardLineRe.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		cards = append(cards, audio.CardInfo{
			ID:          m[1],
			ShortName:   m[2],
			Description: m[3],
		})
	}
	return cards, nil
}

func enumerateFromProc() ([]audio.AudioDeviceInfo, error) {
	var devices []audio.AudioDeviceInfo

	entries, err := os.ReadDir("/proc/asound")
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		name := entry.Name()
		if !strings.HasPrefix(name, "card") {
			continue
		}
		cardNum := strings.TrimPrefix(name, "card")
		cardPath := filepath.Join("/proc/asound", name)

		cardEntries, err := os.ReadDir(cardPath)
		if err != nil {
			continue
		}

		for _, ce := range cardEntries {
			pcmName := ce.Name()
			if !strings.HasPrefix(pcmName, "pcm") {
				continue
			}

			if strings.HasSuffix(pcmName, "p") {
				if dev := readPCMInfo(cardPath, pcmName, "PLAYBACK", cardNum); dev != nil {
					devices = append(devices, *dev)
				}
			}
			if strings.HasSuffix(pcmName, "c") {
				if dev := readPCMInfo(cardPath, pcmName, "CAPTURE", cardNum); dev != nil {
					devices = append(devices, *dev)
				}
			}
		}
	}

	return devices, nil
}

func readPCMInfo(cardPath, pcmDir, streamType, cardNum string) *audio.AudioDeviceInfo {
	infoPath := filepath.Join(cardPath, pcmDir, "info")
	infoData, err := os.ReadFile(infoPath)
	if err != nil {
		return nil
	}
	infoContent := string(infoData)

	// Extract device number from pcmNp or pcmNc
	deviceNum := "0"
	if len(pcmDir) > 4 {
		deviceNum = pcmDir[3 : len(pcmDir)-1]
	}

	deviceName := fmt.Sprintf("hw:%s,%s", cardNum, deviceNum)
	dev := audio.NewAudioDeviceInfo(deviceName, "ALSA")

	// Try to read stream info for channels
	streamPath := filepath.Join(cardPath, "stream0")
	if streamData, err := os.ReadFile(streamPath); err == nil {
		if channels := parseStreamChannels(string(streamData), streamType); channels > 0 {
			switch streamType {
			case "PLAYBACK":
				dev.OutputChannels = channels
			case "CAPTURE":
				dev.InputChannels = channels
			}
		}
	} else {
		// Default to stereo
		switch streamType {
		case "PLAYBACK":
			dev.OutputChannels = 2
		case "CAPTURE":
			dev.InputChannels = 2
		}
	}

	dev.UpdateDeviceType()

	// Check if device is in use
	inUse := strings.Contains(infoContent, "subdevices_avail: 0") &&
		strings.Contains(infoContent, "subdevices_count: 1")
	if inUse {
		dev.Name = dev.Name + " (IN USE)"
	}

	// Parse supported sample rates from sub_stream info
	subPath := filepath.Join(cardPath, pcmDir, "sub0", "hw_params")
	if hwData, err := os.ReadFile(subPath); err == nil {
		dev.SupportedSampleRates = parseHWParamRates(string(hwData))
	}
	if len(dev.SupportedSampleRates) == 0 {
		dev.SupportedSampleRates = []uint32{8000, 11025, 22050, 44100, 48000, 88200, 96000, 176400, 192000}
	}

	return &dev
}

func parseStreamChannels(content, streamType string) uint32 {
	sectionStart := "Playback:"
	if streamType == "CAPTURE" {
		sectionStart = "Capture:"
	}

	idx := strings.Index(content, sectionStart)
	if idx < 0 {
		return 0
	}

	scanner := bufio.NewScanner(strings.NewReader(content[idx:]))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "Channels:") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				if ch, err := strconv.ParseUint(strings.TrimSpace(parts[1]), 10, 32); err == nil {
					return uint32(ch)
				}
			}
		}
	}
	return 0
}

func parseHWParamRates(content string) []uint32 {
	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "rate:") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				if rate, err := strconv.ParseUint(strings.TrimSpace(parts[1]), 10, 32); err == nil {
					return []uint32{uint32(rate)}
				}
			}
		}
	}
	return nil
}
