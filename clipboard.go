package main

import "fmt"

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
