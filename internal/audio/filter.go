package audio

import "strings"

// FilterDevices applies card and device name filters, and optionally deduplicates.
func FilterDevices(devices []AudioDeviceInfo, cardFilter, deviceFilter string, showAll bool) []AudioDeviceInfo {
	filtered := devices

	if cardFilter != "" {
		cardNum := strings.TrimPrefix(cardFilter, "card")
		var result []AudioDeviceInfo
		for _, d := range filtered {
			if strings.Contains(d.Name, "hw:"+cardNum) ||
				strings.Contains(d.Name, "card"+cardNum) ||
				strings.Contains(d.Name, "CARD="+cardFilter) {
				result = append(result, d)
			}
		}
		filtered = result
	}

	if deviceFilter != "" {
		nameLower := strings.ToLower(deviceFilter)
		var result []AudioDeviceInfo
		for _, d := range filtered {
			if strings.Contains(strings.ToLower(d.Name), nameLower) {
				result = append(result, d)
			}
		}
		filtered = result
	}

	if !showAll {
		filtered = deduplicateDevices(filtered)
	}

	return filtered
}

func deduplicateDevices(devices []AudioDeviceInfo) []AudioDeviceInfo {
	seen := make(map[string]bool)
	var result []AudioDeviceInfo
	for _, d := range devices {
		// Skip virtual/duplicate devices
		if strings.HasPrefix(d.Name, "dmix:") ||
			strings.HasPrefix(d.Name, "dsnoop:") ||
			strings.HasPrefix(d.Name, "surround") ||
			strings.HasPrefix(d.Name, "iec958:") {
			continue
		}

		simplified := d.Name
		if strings.HasPrefix(simplified, "plughw:") {
			simplified = strings.Replace(simplified, "plughw:", "hw:", 1)
		}

		if seen[simplified] {
			continue
		}
		seen[simplified] = true
		result = append(result, d)
	}
	return result
}
