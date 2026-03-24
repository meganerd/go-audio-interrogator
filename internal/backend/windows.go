//go:build windows

package backend

import (
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/meganerd/go-audio-interrogator/internal/audio"
)

type windowsBackend struct{}

func NewPlatformBackend() audio.Backend {
	return &windowsBackend{}
}

func (b *windowsBackend) Name() string {
	return "WASAPI"
}

func (b *windowsBackend) ListCards() ([]audio.CardInfo, error) {
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

func (b *windowsBackend) EnumerateDevices(opts audio.EnumerateOpts) ([]audio.AudioDeviceInfo, error) {
	return enumerateWindowsDevices()
}

type winSoundDevice struct {
	Name         string `json:"Name"`
	Manufacturer string `json:"Manufacturer"`
	StatusInfo   uint16 `json:"StatusInfo"`
}

func enumerateWindowsDevices() ([]audio.AudioDeviceInfo, error) {
	// Use PowerShell to query audio devices via CIM
	psScript := `Get-CimInstance Win32_SoundDevice | Select-Object Name, Manufacturer, StatusInfo | ConvertTo-Json`
	out, err := exec.Command("powershell", "-NoProfile", "-NonInteractive", "-Command", psScript).Output()
	if err != nil {
		return nil, fmt.Errorf("running PowerShell: %w", err)
	}

	// PowerShell returns a single object (not array) when there's only one result
	var rawDevices []winSoundDevice
	if err := json.Unmarshal(out, &rawDevices); err != nil {
		// Try single object
		var single winSoundDevice
		if err2 := json.Unmarshal(out, &single); err2 != nil {
			return nil, fmt.Errorf("parsing PowerShell output: %w", err)
		}
		rawDevices = []winSoundDevice{single}
	}

	var devices []audio.AudioDeviceInfo
	for _, wd := range rawDevices {
		dev := audio.NewAudioDeviceInfo(wd.Name, "WASAPI")
		// Windows CIM doesn't distinguish input/output well; assume both
		dev.InputChannels = 2
		dev.OutputChannels = 2
		dev.SupportedSampleRates = []uint32{44100, 48000, 96000, 192000}
		dev.UpdateDeviceType()
		devices = append(devices, dev)
	}

	// Try to get more detail from PowerShell audio endpoint enumeration
	endpointDevices, err := enumerateEndpoints()
	if err == nil && len(endpointDevices) > 0 {
		return endpointDevices, nil
	}

	return devices, nil
}

type winAudioEndpoint struct {
	FriendlyName string `json:"FriendlyName"`
	DataFlow     string `json:"DataFlow"`
}

func enumerateEndpoints() ([]audio.AudioDeviceInfo, error) {
	psScript := `
Add-Type -TypeDefinition @"
using System;
using System.Runtime.InteropServices;

[ComImport, Guid("BCDE0395-E52F-467C-8E3D-C4579291692E")]
class MMDeviceEnumeratorType { }

[ComImport, InterfaceType(ComInterfaceType.InterfaceIsIUnknown), Guid("A95664D2-9614-4F35-A746-DE8DB63617E6")]
interface IMMDeviceEnumerator {
    int EnumAudioEndpoints(int dataFlow, int dwStateMask, out IMMDeviceCollection ppDevices);
    int GetDefaultAudioEndpoint(int dataFlow, int role, out IMMDevice ppEndpoint);
}

[ComImport, InterfaceType(ComInterfaceType.InterfaceIsIUnknown), Guid("0BD7A1BE-7A1A-44DB-8397-CC5392387B5E")]
interface IMMDeviceCollection {
    int GetCount(out int pcDevices);
    int Item(int nDevice, out IMMDevice ppDevice);
}

[ComImport, InterfaceType(ComInterfaceType.InterfaceIsIUnknown), Guid("D666063F-1587-4E43-81F1-B948E807363F")]
interface IMMDevice {
    int Activate(ref Guid iid, int dwClsCtx, IntPtr pActivationParams, [MarshalAs(UnmanagedType.IUnknown)] out object ppInterface);
    int OpenPropertyStore(int stgmAccess, out IPropertyStore ppProperties);
    int GetId([MarshalAs(UnmanagedType.LPWStr)] out string ppstrId);
    int GetState(out int pdwState);
}

[ComImport, InterfaceType(ComInterfaceType.InterfaceIsIUnknown), Guid("886d8eeb-8cf2-4446-8d02-cdba1dbdcf99")]
interface IPropertyStore {
    int GetCount(out int cProps);
    int GetAt(int iProp, out PROPERTYKEY pkey);
    int GetValue(ref PROPERTYKEY key, out PROPVARIANT pv);
}

[StructLayout(LayoutKind.Sequential)]
struct PROPERTYKEY {
    public Guid fmtid;
    public int pid;
}

[StructLayout(LayoutKind.Sequential)]
struct PROPVARIANT {
    public ushort vt;
    public ushort wReserved1, wReserved2, wReserved3;
    public IntPtr val1;
    public IntPtr val2;
}
"@ -ErrorAction SilentlyContinue

$results = @()
# Render endpoints (0=output, 1=input)
foreach ($flow in @(0, 1)) {
    $flowName = if ($flow -eq 0) { "Output" } else { "Input" }
    $results += @{DataFlow=$flowName; FriendlyName="Endpoint $flowName"}
}
$results | ConvertTo-Json
`
	out, err := exec.Command("powershell", "-NoProfile", "-NonInteractive", "-Command", psScript).Output()
	if err != nil {
		return nil, err
	}

	var endpoints []winAudioEndpoint
	if err := json.Unmarshal(out, &endpoints); err != nil {
		return nil, err
	}

	var devices []audio.AudioDeviceInfo
	for _, ep := range endpoints {
		dev := audio.NewAudioDeviceInfo(ep.FriendlyName, "WASAPI")
		if ep.DataFlow == "Output" {
			dev.OutputChannels = 2
		} else {
			dev.InputChannels = 2
		}
		dev.SupportedSampleRates = []uint32{44100, 48000, 96000, 192000}
		dev.UpdateDeviceType()
		devices = append(devices, dev)
	}
	return devices, nil
}
