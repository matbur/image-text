package image

import (
	"fmt"
	"regexp"
	"strconv"
)

var knownSizes = map[string]Size{
	// Screen Standards
	"cga":   NewSize(320, 200),
	"qvga":  NewSize(320, 240),
	"vga":   NewSize(640, 480),
	"wvga":  NewSize(800, 480),
	"svga":  NewSize(800, 480),
	"wsvga": NewSize(1024, 600),
	"xga":   NewSize(1024, 768),
	"wxga":  NewSize(1280, 800),
	"wsxga": NewSize(1440, 900),
	"wuxga": NewSize(1920, 1200),
	"wqxga": NewSize(2560, 1600),

	// Video Standards
	"ntsc":   NewSize(720, 480),
	"pal":    NewSize(768, 576),
	"hd720":  NewSize(1280, 720),
	"hd1080": NewSize(1920, 1080),
}

var sizePattern = regexp.MustCompile(`^([1-9][0-9]{0,5})x([1-9][0-9]{0,5})$`)

type Size struct {
	width, height int
}

func NewSize(width, height int) Size {
	return Size{
		width:  width,
		height: height,
	}
}

func NewSizeFromString(s string) (Size, error) {
	size, ok := knownSizes[s]
	if ok {
		return size, nil
	}

	ss := sizePattern.FindStringSubmatch(s)
	if len(ss) != 3 {
		return DefaultSize(), fmt.Errorf("size '%s' is not valid: %w", s, errorMalformed)
	}

	width, err := strconv.Atoi(ss[1])
	if err != nil {
		return DefaultSize(), fmt.Errorf("bad width '%s': %v: %w", ss[1], err, errorUnexpected)
	}

	height, err := strconv.Atoi(ss[2])
	if err != nil {
		return DefaultSize(), fmt.Errorf("bad height '%s': %v: %w", ss[2], err, errorUnexpected)
	}

	return NewSize(width, height), nil
}

func DefaultSize() Size {
	return knownSizes["vga"]
}

func (s Size) Width() int {
	return s.width
}

func (s Size) Height() int {
	return s.height
}

func (s Size) String() string {
	return fmt.Sprintf("%dx%d", s.width, s.height)
}

func KnownSizes() map[string]Size {
	return knownSizes
}
