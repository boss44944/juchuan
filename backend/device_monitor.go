package main

import "time"

func StartDeviceMonitor(manager *DeviceManager, notify func(DeviceEvent)) {
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			now := time.Now().Unix()
			for _, device := range manager.List() {
				if now-device.LastSeen > 120 && device.Status != DeviceStatusOffline {
					device.Status = DeviceStatusOffline
					if notify != nil {
						notify(DeviceEvent{
							Type: DeviceOfflineEvent,
							Data: *device,
						})
					}
				}
			}
		}
	}()
}
