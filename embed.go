package main

import (
	"embed"
)

// StaticFiles contains frontend assets.
//
//go:embed static/*
var StaticFiles embed.FS
