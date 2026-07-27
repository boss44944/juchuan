package main

import (
	"fmt"
	"os/exec"
	"runtime"
)

type Clipboard struct{}

func (c *Clipboard) Set(text string) error {
	return SetClipboard(text)
}

func (c *Clipboard) Copy(text string) error {
	return c.Set(text)
}

func (c *Clipboard) String() string {
	return fmt.Sprint("system clipboard")
}

func SetClipboard(text string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("pbcopy")
	case "windows":
		cmd = exec.Command("clip")
	default:
		cmd = exec.Command("sh", "-c", "xclip -selection clipboard")
	}

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	if _, err := stdin.Write([]byte(text)); err != nil {
		return err
	}

	if err := stdin.Close(); err != nil {
		return err
	}

	return cmd.Wait()
}
