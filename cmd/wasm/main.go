//go:build wasm

package main

import (
	"bytes"
	_ "embed"
	"encoding/base64"
	"log/slog"
	"syscall/js"

	"github.com/matbur/image-text/image"
)

func imageText(this js.Value, p []js.Value) any {
	bgColor := p[0].Get("bgColor").String()
	fgColor := p[0].Get("fgColor").String()
	size := p[0].Get("size").String()
	text := p[0].Get("text").String()

	font := p[0].Get("font").String()
	format := p[0].Get("format").String()
	if format == "undefined" {
		format = ""
	}

	img, err := image.New(size, bgColor, fgColor, text, font, format)
	if err != nil {
		slog.Error("Failed to create image", "err", err)
		return map[string]any{"err": err.Error()}
	}

	var buff bytes.Buffer
	if err := img.Draw(&buff); err != nil {
		slog.Error("Failed to draw image", "err", err)
		return map[string]any{"err": err.Error()}
	}

	return map[string]any{"imageBase64": base64.StdEncoding.EncodeToString(buff.Bytes())}
}

func registerFont(this js.Value, p []js.Value) any {
	key := p[0].String()
	length := p[1].Get("length").Int()
	data := make([]byte, length)
	js.CopyBytesToGo(data, p[1])

	if err := image.RegisterFont(key, data); err != nil {
		return map[string]any{"err": err.Error()}
	}

	return nil
}

func registerCallbacks() {
	js.Global().Set("imageText", js.FuncOf(imageText))
	js.Global().Set("registerFont", js.FuncOf(registerFont))
}

func main() {
	registerCallbacks()

	<-make(chan struct{})
}
