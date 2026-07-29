package main

const (
	DeviceOnlineEvent  = "DEVICE_ONLINE"
	DeviceOfflineEvent = "DEVICE_OFFLINE"
)

type DeviceEvent struct {
	Type string  `json:"type"`
	Data Device  `json:"data"`
}
