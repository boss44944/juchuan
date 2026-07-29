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
	url := LocalURL(port)

	log.Println("==============================")
	log.Println("Juchuan 菊传")
	log.Println("")
	log.Println("访问地址:")
	log.Println(url)
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
		OpenBrowser(url)
	}

	quit := make(chan struct{}, 1)
	go StartTray(url, quit)

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
