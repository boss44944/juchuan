package main

import "sync"

const (
	DeviceStatusOnline  = "online"
	DeviceStatusOffline = "offline"
)

type Device struct {
	ID           string `json:"id"`
	DisplayName  string `json:"display_name"`
	Role         string `json:"role"`
	Platform     string `json:"platform"`
	Browser      string `json:"browser"`
	DeviceSecret string `json:"device_secret,omitempty"`
	LastSeen     int64  `json:"last_seen"`
	Status       string `json:"status"`
}

type DeviceManager struct {
	mu      sync.Mutex
	devices map[string]*Device
}

func NewDeviceManager() *DeviceManager {
	return &DeviceManager{devices: make(map[string]*Device)}
}

func (m *DeviceManager) Add(d *Device) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.devices[d.ID] = d
}

func (m *DeviceManager) Get(id string) (*Device, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	d, ok := m.devices[id]
	return d, ok
}

func (m *DeviceManager) List() []*Device {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := make([]*Device, 0, len(m.devices))
	for _, d := range m.devices {
		out = append(out, d)
	}
	return out
}

func (m *DeviceManager) NameExists(name string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, d := range m.devices {
		if d.DisplayName == name {
			return true
		}
	}
	return false
}
