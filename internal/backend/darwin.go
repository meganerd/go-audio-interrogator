//go:build darwin

package backend

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/meganerd/go-audio-interrogator/internal/audio"
)

type darwinBackend struct{}

func NewPlatformBackend() audio.Backend {
	return &darwinBackend{}
}

func (b *darwinBackend) Name() string {
	return "CoreAudio"
}

func (b *darwinBackend) ListCards() ([]audio.CardInfo, error) {
	devices, err := b.EnumerateDevices(audio.EnumerateOpts{})
	if err != nil {
		return nil, err
	}

	var cards []audio.CardInfo
	for i, d := range devices {
		cards = append(cards, audio.CardInfo{
			ID:          fmt.Sprintf("%d", i),
			ShortName:   d.Name,
			Description: d.Driver,
		})
	}
	return cards, nil
}

func (b *darwinBackend) EnumerateDevices(opts audio.EnumerateOpts) ([]audio.AudioDeviceInfo, error) {
	return enumerateDarwinDevices()
}

type spAudioData struct {
	SPAudioDataType []spAudioDevice `json:"SPAudioDataType"`
}

type spAudioDevice struct {
	Name             string            `json:"_name"`
	CoreAudioDevices []coreAudioDevice `json:"coreaudio_device_list,omitempty"`
}

type coreAudioDevice struct {
	Name           string `json:"_name"`
	InputChannels  int    `json:"coreaudio_device_input,omitempty"`
	OutputChannels int    `json:"coreaudio_device_output,omitempty"`
	SampleRate     string `json:"coreaudio_device_srate,omitempty"`
}

func enumerateDarwinDevices() ([]audio.AudioDeviceInfo, error) {
	out, err := exec.Command("system_profiler", "SPAudioDataType", "-json").Output()
	if err != nil {
		return nil, fmt.Errorf("running system_profiler: %w", err)
	}

	var data spAudioData
	if err := json.Unmarshal(out, &data); err != nil {
		return nil, fmt.Errorf("parsing system_profiler output: %w", err)
	}

	var devices []audio.AudioDeviceInfo
	for _, section := range data.SPAudioDataType {
		for _, caDev := range section.CoreAudioDevices {
			dev := audio.NewAudioDeviceInfo(caDev.Name, "CoreAudio")
			dev.InputChannels = uint32(caDev.InputChannels)
			dev.OutputChannels = uint32(caDev.OutputChannels)

			if caDev.SampleRate != "" {
				rateStr := strings.TrimSpace(strings.Replace(caDev.SampleRate, "_Hz", "", 1))
				rateStr = strings.Replace(rateStr, " Hz", "", 1)
				var rate uint32
				if _, err := fmt.Sscanf(rateStr, "%d", &rate); err == nil && rate > 0 {
					dev.DefaultSampleRate = rate
					dev.SupportedSampleRates = []uint32{rate}
				}
			}

			dev.UpdateDeviceType()
			if dev.DeviceType != "Unknown" {
				devices = append(devices, dev)
			}
		}
	}

	return devices, nil
}
