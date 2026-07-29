package main

import (
	"database/sql"
)

func LoadDevices(db *sql.DB, manager *DeviceManager) error {
	rows, err := db.Query(`SELECT id, display_name, role, platform, browser, device_secret, last_seen FROM devices`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		d := &Device{}
		if err := rows.Scan(&d.ID, &d.DisplayName, &d.Role, &d.Platform, &d.Browser, &d.DeviceSecret, &d.LastSeen); err != nil {
			return err
		}
		d.Status = DeviceStatusOffline
		manager.Add(d)
	}

	return rows.Err()
}

func SaveDevice(db *sql.DB, d *Device) error {
	_, err := db.Exec(`
INSERT INTO devices(id, display_name, role, platform, browser, device_secret, last_seen)
VALUES(?,?,?,?,?,?,?)
ON CONFLICT(id) DO UPDATE SET
 display_name=excluded.display_name,
 role=excluded.role,
 platform=excluded.platform,
 browser=excluded.browser,
 device_secret=excluded.device_secret,
 last_seen=excluded.last_seen
`, d.ID, d.DisplayName, d.Role, d.Platform, d.Browser, d.DeviceSecret, d.LastSeen)
	return err
}

func UpdateDeviceName(db *sql.DB, id string, name string) error {
	_, err := db.Exec(`UPDATE devices SET display_name=? WHERE id=?`, name, id)
	return err
}

func DeleteDevice(db *sql.DB, id string) error {
	_, err := db.Exec(`DELETE FROM devices WHERE id=?`, id)
	return err
}
