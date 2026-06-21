package image

import (
	"fmt"
	"path"
	"regexp"
	"sort"
	"strings"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"

	"github.com/matbur/image-text/resources"
)

var (
	knownFonts        = loadAllFonts()
	fontSuffixPattern = regexp.MustCompile(`(?i)-regular$|-variablefont[^.]*$`)
)

func loadAllFonts() map[string]*truetype.Font {
	entries, err := resources.Static.ReadDir("fonts")
	if err != nil {
		panic(fmt.Errorf("failed to read fonts directory: %w", err))
	}

	fonts := make(map[string]*truetype.Font)
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(strings.ToLower(entry.Name()), ".ttf") {
			continue
		}

		filePath := path.Join("fonts", entry.Name())
		fnt, err := loadFontFile(filePath)
		if err != nil {
			panic(fmt.Errorf("failed to load font %q: %w", entry.Name(), err))
		}

		key := fontKey(entry.Name())
		if _, exists := fonts[key]; exists {
			panic(fmt.Errorf("duplicate font key %q", key))
		}
		fonts[key] = fnt
	}

	if len(fonts) == 0 {
		panic("no fonts found in resources/fonts")
	}

	return fonts
}

func fontKey(filename string) string {
	name := strings.TrimSuffix(filename, path.Ext(filename))
	name = fontSuffixPattern.ReplaceAllString(name, "")
	return pascalToSnake(name)
}

func pascalToSnake(s string) string {
	var b strings.Builder
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			b.WriteRune('_')
		}

		switch {
		case r >= 'A' && r <= 'Z':
			b.WriteRune(r - 'A' + 'a')
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			b.WriteRune(r)
		case r == '-' || r == '_' || r == ',':
			if b.Len() > 0 && b.String()[b.Len()-1] != '_' {
				b.WriteRune('_')
			}
		}
	}

	return strings.Trim(b.String(), "_")
}

func DefaultFont() *truetype.Font {
	return knownFonts["ubuntu_mono"]
}

func NewFontFromString(s string) (*truetype.Font, error) {
	if s == "" {
		return DefaultFont(), nil
	}

	s = strings.ToLower(s)
	if fnt, ok := knownFonts[s]; ok {
		return fnt, nil
	}

	return DefaultFont(), fmt.Errorf("font '%s' is not valid: %w", s, errorMalformed)
}

func KnownFonts() map[string]*truetype.Font {
	return knownFonts
}

func KnownFontNames() []string {
	names := make([]string, 0, len(knownFonts))
	for name := range knownFonts {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func KnownFontFilenames() map[string]string {
	entries, err := resources.Static.ReadDir("fonts")
	if err != nil {
		panic(fmt.Errorf("failed to read fonts directory: %w", err))
	}

	filenames := make(map[string]string)
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(strings.ToLower(entry.Name()), ".ttf") {
			continue
		}

		key := fontKey(entry.Name())
		if _, exists := filenames[key]; exists {
			panic(fmt.Errorf("duplicate font key %q", key))
		}
		filenames[key] = entry.Name()
	}

	return filenames
}

func RegisterFont(key string, data []byte) error {
	key = strings.ToLower(strings.TrimSpace(key))
	if key == "" {
		return fmt.Errorf("font name is empty: %w", errorMalformed)
	}

	fnt, err := truetype.Parse(data)
	if err != nil {
		return fmt.Errorf("failed to parse font %q: %w", key, err)
	}

	knownFonts[key] = fnt
	return nil
}

func loadFontFile(filePath string) (*truetype.Font, error) {
	bb, err := resources.Static.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file '%s': %v: %w", filePath, err, errorUnexpected)
	}

	fnt, err := truetype.Parse(bb)
	if err != nil {
		return nil, fmt.Errorf("failed to parse font '%s': %v: %w", filePath, err, errorUnexpected)
	}

	return fnt, nil
}

func loadFace(fnt *truetype.Font, size int) font.Face {
	opts := &truetype.Options{
		Size: float64(size),
	}

	return truetype.NewFace(fnt, opts)
}
