package output

import (
	"encoding/json"
	"fmt"

	"github.com/meganerd/go-audio-interrogator/internal/audio"
)

func PrintJSON(info audio.SystemAudioInfo) error {
	data, err := json.MarshalIndent(info, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling JSON: %w", err)
	}
	fmt.Println(string(data))
	return nil
}
