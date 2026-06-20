# image-text

Generates placeholder images with text, colors, sizes, and fonts.

Inspired by [dummyimage.com](https://dummyimage.com/).

## Quick start

With Docker:

```bash
make start
```

Open http://localhost:8080

Local development:

```bash
make install-tools   # once
make build-local
./tmp/main
```

### Live reload

For day-to-day development, use [air](https://github.com/air-verse/air). It watches source files and rebuilds automatically (via `make build-local` in `.air.toml`):

```bash
make install-tools   # once
air
```

Open http://localhost:8080 — edits to Go, templ, and WASM sources trigger a rebuild and server restart.

Environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT`   | `8080`  | HTTP port   |

## Web UI

| Path       | Description                                      |
|------------|--------------------------------------------------|
| `/`        | Home — choose online or offline mode             |
| `/online`  | Generate images on the server (HTMX)             |
| `/offline` | Generate images in the browser (WebAssembly)     |
| `/docs`    | API documentation                                |
| `/docs.json` | Machine-readable API reference                 |

## API

Generate a PNG image with a GET request:

```http
GET /{size}/{background}/{foreground}?text=hello&font=ubuntu_mono
```

### Path segments

| Segment      | Description                                      | Examples                          |
|--------------|--------------------------------------------------|-----------------------------------|
| `size`       | Named preset or `{width}x{height}`               | `vga`, `hd720`, `640x480`         |
| `background` | Named color or hex code                          | `steel_blue`, `000`, `1c6`        |
| `foreground` | Named color or hex code                          | `yellow`, `FFFF00`, `53f`         |

### Query parameters

| Parameter | Required | Default       | Description                |
|-----------|----------|---------------|----------------------------|
| `text`    | no       | size string   | Text rendered on the image |
| `font`    | no       | `ubuntu_mono` | Font name                |

### Examples

Named colors and size:

```http
GET /hd720/steel_blue/yellow?text=rendered+text&font=ubuntu_mono
```

Hex colors and custom size:

```http
GET /320x200/000/FFFF00?font=open_sans
```

### Reference data

Full lists of colors, sizes, and fonts:

- HTML: http://localhost:8080/docs
- JSON: http://localhost:8080/docs.json

Add fonts by placing `.ttf` files in `resources/fonts/`.

## Development

```bash
make test
air             # live reload (recommended for local dev)
make templ      # regenerate templ templates
make build-wasm # rebuild offline WebAssembly bundle
```
