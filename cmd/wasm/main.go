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

	slog.Info("imageText",
		"stringify", js.Global().Get("JSON").Call("stringify", p[0]),
	)

	img, err := image.New(size, bgColor, fgColor, text)
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

func registerCallbacks() {
	js.Global().Set("imageText", js.FuncOf(imageText))
}

func main() {
	registerCallbacks()

	<-make(chan struct{})
}
