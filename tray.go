//go:build darwin || windows
// +build darwin windows

package main

import (
	_ "embed"

	"github.com/getlantern/systray"
)

//go:embed app-logo.png
var trayIconPNG []byte

// StartTray starts the desktop tray menu.
// The tray remains active after the browser window is closed.
func StartTray(url string, quit chan struct{}) {
	systray.Run(func() {
		if len(trayIconPNG) > 0 {
			systray.SetIcon(trayIconPNG)
		}

		systray.SetTitle("Juchuan 菊传")
		systray.SetTooltip("局域网文件传输工具")

		open := systray.AddMenuItem("打开 Web 页面", "打开 Juchuan")
		exit := systray.AddMenuItem("退出程序", "关闭 Juchuan")

		go func() {
			for {
				select {
				case <-open.ClickedCh:
					OpenBrowser(url)
				case <-exit.ClickedCh:
					select {
					case quit <- struct{}{}:
					default:
					}
					systray.Quit()
					return
				}
			}
		}()
	}, func() {})
}
