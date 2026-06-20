//go:build !wasm

package resources

import "embed"

//go:embed * fonts/*
var Static embed.FS
