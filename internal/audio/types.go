package audio

// AudioDeviceInfo holds information about a single audio device.
type AudioDeviceInfo struct {
	Name                 string   `json:"name"`
	DeviceType           string   `json:"device_type"`
	InputChannels        uint32   `json:"input_channels"`
	OutputChannels       uint32   `json:"output_channels"`
	SupportedSampleRates []uint32 `json:"supported_sample_rates"`
	SupportedBufferSizes []uint32 `json:"supported_buffer_sizes"`
	DefaultSampleRate    uint32   `json:"default_sample_rate"`
	DefaultBufferSize    uint32   `json:"default_buffer_size"`
	Driver               string   `json:"driver"`
}

// NewAudioDeviceInfo creates an AudioDeviceInfo with sensible defaults.
func NewAudioDeviceInfo(name, driver string) AudioDeviceInfo {
	return AudioDeviceInfo{
		Name:                 name,
		DeviceType:           "Unknown",
		SupportedBufferSizes: []uint32{64, 128, 256, 512, 1024, 2048, 4096},
		DefaultSampleRate:    44100,
		DefaultBufferSize:    1024,
		Driver:               driver,
	}
}

// UpdateDeviceType sets DeviceType based on channel counts.
func (d *AudioDeviceInfo) UpdateDeviceType() {
	switch {
	case d.InputChannels > 0 && d.OutputChannels > 0:
		d.DeviceType = "Input/Output"
	case d.InputChannels > 0:
		d.DeviceType = "Input"
	case d.OutputChannels > 0:
		d.DeviceType = "Output"
	default:
		d.DeviceType = "Unknown"
	}
}

// HasInput returns true if the device has input capabilities.
func (d *AudioDeviceInfo) HasInput() bool {
	return d.InputChannels > 0
}

// HasOutput returns true if the device has output capabilities.
func (d *AudioDeviceInfo) HasOutput() bool {
	return d.OutputChannels > 0
}

// SystemAudioInfo holds system-wide audio information.
type SystemAudioInfo struct {
	Devices            []AudioDeviceInfo `json:"devices"`
	DefaultInput       string            `json:"default_input,omitempty"`
	DefaultOutput      string            `json:"default_output,omitempty"`
	TotalInputDevices  int               `json:"total_input_devices"`
	TotalOutputDevices int               `json:"total_output_devices"`
}

// NewSystemAudioInfo creates a SystemAudioInfo from a device list.
func NewSystemAudioInfo(devices []AudioDeviceInfo) SystemAudioInfo {
	info := SystemAudioInfo{Devices: devices}
	for _, d := range devices {
		if d.HasInput() {
			info.TotalInputDevices++
		}
		if d.HasOutput() {
			info.TotalOutputDevices++
		}
	}
	return info
}
