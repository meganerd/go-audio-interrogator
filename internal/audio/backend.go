package audio

// Backend is the interface that platform-specific audio backends must implement.
type Backend interface {
	// EnumerateDevices returns all audio devices detected by this backend.
	EnumerateDevices(opts EnumerateOpts) ([]AudioDeviceInfo, error)

	// ListCards returns a human-readable summary of available audio cards.
	ListCards() ([]CardInfo, error)

	// Name returns the backend name (e.g., "ALSA", "CoreAudio", "WASAPI").
	Name() string
}

// EnumerateOpts controls device enumeration behavior.
type EnumerateOpts struct {
	// NoProc disables /proc/asound access on Linux to avoid interfering with active streams.
	NoProc bool
}

// CardInfo represents a summary of an audio card.
type CardInfo struct {
	ID          string
	ShortName   string
	Description string
}
