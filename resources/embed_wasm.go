//go:build wasm

package resources

import "embed"

//go:embed fonts/UbuntuMono-Regular.ttf
var Static embed.FS
