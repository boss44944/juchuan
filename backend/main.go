package main

import (
	"log"
	"strings"
)

func main() {
	s, err := NewServer()
	if err != nil {
		log.Fatal(err)
	}

	port := strings.TrimPrefix(s.Address(), ":")
	baseURL := LocalURL(port)
	serverURL := baseURL + "/server/devices"
	clientURL := baseURL + "/client/inbox"

	log.Println("==============================")
	log.Println("Juchuan 菊传")
	log.Println("")
	log.Println("电脑管理页面:")
	log.Println(serverURL)
	log.Println("手机扫码页面:")
	log.Println(clientURL)
	if cfgPath, err := configFilePath(); err == nil {
		log.Println("配置文件:")
		log.Println(cfgPath)
	}
	if pwd := s.CurrentPassword(); strings.TrimSpace(pwd) != "" {
		log.Println("当前访问密码:")
		log.Println(pwd)
	} else {
		log.Println("当前访问密码: 未启用")
	}
	log.Println("==============================")

	if s.ShouldAutoOpen() {
		OpenBrowser(serverURL)
	}

	quit := make(chan struct{}, 1)
	go StartTray(serverURL, quit)

	serverErr := make(chan error, 1)
	go func() {
		serverErr <- s.Start()
	}()

	select {
	case <-quit:
		if err := s.Shutdown(); err != nil {
			log.Println("shutdown error:", err)
		}
	case err := <-serverErr:
		if err != nil {
			log.Fatal(err)
		}
	}
}
