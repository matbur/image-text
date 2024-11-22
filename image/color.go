package image

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var colorPattern = regexp.MustCompile(`^(?:[0-9A-Fa-f]{3}){1,2}$`)

type Color struct {
	r, g, b uint8
}

type number interface {
	int | uint8 | uint64
}

func NewColor[T number](r, g, b T) Color {
	return Color{
		r: uint8(r),
		g: uint8(g),
		b: uint8(b),
	}
}

func NewBackgroundColorFromString(s string) (Color, error) {
	return NewColorFromString(s)
}

func NewForegroundColorFromString(s string) (Color, error) {
	c, err := NewColorFromString(s)
	if err != nil {
		return DefaultForegroundColor(), err
	}
	return c, nil
}

func NewColorFromString(s string) (Color, error) {
	if s == "" {
		return DefaultBackgroundColor(), fmt.Errorf("color is empty: %w", errorMissing)
	}

	if color, ok := knownColors[s]; ok {
		return color, nil
	}

	if !colorPattern.MatchString(s) {
		return DefaultBackgroundColor(), fmt.Errorf("color '%s' is not valid: %w", s, errorMalformed)
	}

	var ss []string
	switch len(s) {
	case 3:
		ss = strings.Split(s, "")
		ss[0] += ss[0]
		ss[1] += ss[1]
		ss[2] += ss[2]
	case 6:
		ss = []string{s[0:2], s[2:4], s[4:6]}
	default:
		return DefaultBackgroundColor(), fmt.Errorf("bad color '%s': %w", s, errorUnexpected)
	}

	red, err := strconv.ParseUint(ss[0], 16, 8)
	if err != nil {
		return DefaultBackgroundColor(), fmt.Errorf("bad red '%s': %v: %w", ss[0], err, errorUnexpected)
	}

	green, err := strconv.ParseUint(ss[1], 16, 8)
	if err != nil {
		return DefaultBackgroundColor(), fmt.Errorf("bad green '%s': %v: %w", ss[1], err, errorUnexpected)
	}

	blue, err := strconv.ParseUint(ss[2], 16, 8)
	if err != nil {
		return DefaultBackgroundColor(), fmt.Errorf("bad blue '%s': %v: %w", ss[2], err, errorUnexpected)
	}

	return NewColor(red, green, blue), nil
}

func DefaultBackgroundColor() Color {
	return knownColors["steel_blue"]
}

func DefaultForegroundColor() Color {
	return knownColors["yellow"]
}

func (c Color) RGBA() (r, g, b, a uint32) {
	return uint32(c.r), uint32(c.g), uint32(c.b), 255
}

func (c Color) String() string {
	return fmt.Sprintf("#%02x%02x%02x", c.r, c.g, c.b)
}

func KnownColors() map[string]Color {
	return knownColors
}

var knownColors = map[string]Color{
	"alice_blue":              NewColor(240, 248, 255),
	"antique_white":           NewColor(250, 235, 215),
	"aqua":                    NewColor(0, 255, 255),
	"aquamarine":              NewColor(127, 255, 212),
	"azure":                   NewColor(240, 255, 255),
	"beige":                   NewColor(245, 245, 220),
	"bisque":                  NewColor(255, 228, 196),
	"black":                   NewColor(0, 0, 0),
	"blanched_almond":         NewColor(255, 235, 205),
	"blue":                    NewColor(0, 0, 255),
	"blue_violet":             NewColor(138, 43, 226),
	"brown":                   NewColor(165, 42, 42),
	"burly_wood":              NewColor(222, 184, 135),
	"cadet_blue":              NewColor(95, 158, 160),
	"chartreuse":              NewColor(127, 255, 0),
	"chocolate":               NewColor(210, 105, 30),
	"coral":                   NewColor(255, 127, 80),
	"cornflower_blue":         NewColor(100, 149, 237),
	"cornsilk":                NewColor(255, 248, 220),
	"crimson":                 NewColor(220, 20, 60),
	"cyan":                    NewColor(0, 255, 255),
	"dark_blue":               NewColor(0, 0, 139),
	"dark_cyan":               NewColor(0, 139, 139),
	"dark_golden_rod":         NewColor(184, 134, 11),
	"dark_gray":               NewColor(169, 169, 169),
	"dark_green":              NewColor(0, 100, 0),
	"dark_grey":               NewColor(169, 169, 169),
	"dark_khaki":              NewColor(189, 183, 107),
	"dark_magenta":            NewColor(139, 0, 139),
	"dark_olive_green":        NewColor(85, 107, 47),
	"dark_orange":             NewColor(255, 140, 0),
	"dark_orchid":             NewColor(153, 50, 204),
	"dark_red":                NewColor(139, 0, 0),
	"dark_salmon":             NewColor(233, 150, 122),
	"dark_sea_green":          NewColor(143, 188, 143),
	"dark_slate_blue":         NewColor(72, 61, 139),
	"dark_slate_gray":         NewColor(47, 79, 79),
	"dark_slate_grey":         NewColor(47, 79, 79),
	"dark_turquoise":          NewColor(0, 206, 209),
	"dark_violet":             NewColor(148, 0, 211),
	"deep_pink":               NewColor(255, 20, 147),
	"deep_sky_blue":           NewColor(0, 191, 255),
	"dim_gray":                NewColor(105, 105, 105),
	"dim_grey":                NewColor(105, 105, 105),
	"dodger_blue":             NewColor(30, 144, 255),
	"fire_brick":              NewColor(178, 34, 34),
	"floral_white":            NewColor(255, 250, 240),
	"forest_green":            NewColor(34, 139, 34),
	"fuchsia":                 NewColor(255, 0, 255),
	"gainsboro":               NewColor(220, 220, 220),
	"ghost_white":             NewColor(248, 248, 255),
	"gold":                    NewColor(255, 215, 0),
	"golden_rod":              NewColor(218, 165, 32),
	"gray":                    NewColor(128, 128, 128),
	"green":                   NewColor(0, 128, 0),
	"green_yellow":            NewColor(173, 255, 47),
	"grey":                    NewColor(128, 128, 128),
	"honey_dew":               NewColor(240, 255, 240),
	"hot_pink":                NewColor(255, 105, 180),
	"indian_red":              NewColor(205, 92, 92),
	"indigo":                  NewColor(75, 0, 130),
	"ivory":                   NewColor(255, 255, 240),
	"khaki":                   NewColor(240, 230, 140),
	"lavender":                NewColor(230, 230, 250),
	"lavender_blush":          NewColor(255, 240, 245),
	"lawn_green":              NewColor(124, 252, 0),
	"lemon_chiffon":           NewColor(255, 250, 205),
	"light_blue":              NewColor(173, 216, 230),
	"light_coral":             NewColor(240, 128, 128),
	"light_cyan":              NewColor(224, 255, 255),
	"light_golden_rod_yellow": NewColor(250, 250, 210),
	"light_gray":              NewColor(211, 211, 211),
	"light_green":             NewColor(144, 238, 144),
	"light_grey":              NewColor(211, 211, 211),
	"light_pink":              NewColor(255, 182, 193),
	"light_salmon":            NewColor(255, 160, 122),
	"light_sea_green":         NewColor(32, 178, 170),
	"light_sky_blue":          NewColor(135, 206, 250),
	"light_slate_gray":        NewColor(119, 136, 153),
	"light_slate_grey":        NewColor(119, 136, 153),
	"light_steel_blue":        NewColor(176, 196, 222),
	"light_yellow":            NewColor(255, 255, 224),
	"lime":                    NewColor(0, 255, 0),
	"lime_green":              NewColor(50, 205, 50),
	"linen":                   NewColor(250, 240, 230),
	"magenta":                 NewColor(255, 0, 255),
	"maroon":                  NewColor(128, 0, 0),
	"medium_aqua_marine":      NewColor(102, 205, 170),
	"medium_blue":             NewColor(0, 0, 205),
	"medium_orchid":           NewColor(186, 85, 211),
	"medium_purple":           NewColor(147, 112, 219),
	"medium_sea_green":        NewColor(60, 179, 113),
	"medium_slate_blue":       NewColor(123, 104, 238),
	"medium_spring_green":     NewColor(0, 250, 154),
	"medium_turquoise":        NewColor(72, 209, 204),
	"medium_violet_red":       NewColor(199, 21, 133),
	"midnight_blue":           NewColor(25, 25, 112),
	"mint_cream":              NewColor(245, 255, 250),
	"misty_rose":              NewColor(255, 228, 225),
	"moccasin":                NewColor(255, 228, 181),
	"navajo_white":            NewColor(255, 222, 173),
	"navy":                    NewColor(0, 0, 128),
	"old_lace":                NewColor(253, 245, 230),
	"olive":                   NewColor(128, 128, 0),
	"olive_drab":              NewColor(107, 142, 35),
	"orange":                  NewColor(255, 165, 0),
	"orange_red":              NewColor(255, 69, 0),
	"orchid":                  NewColor(218, 112, 214),
	"pale_golden_rod":         NewColor(238, 232, 170),
	"pale_green":              NewColor(152, 251, 152),
	"pale_turquoise":          NewColor(175, 238, 238),
	"pale_violet_red":         NewColor(219, 112, 147),
	"papaya_whip":             NewColor(255, 239, 213),
	"peach_puff":              NewColor(255, 218, 185),
	"peru":                    NewColor(205, 133, 63),
	"pink":                    NewColor(255, 192, 203),
	"plum":                    NewColor(221, 160, 221),
	"powder_blue":             NewColor(176, 224, 230),
	"purple":                  NewColor(128, 0, 128),
	"rebecca_purple":          NewColor(102, 51, 153),
	"red":                     NewColor(255, 0, 0),
	"rosy_brown":              NewColor(188, 143, 143),
	"royal_blue":              NewColor(65, 105, 225),
	"saddle_brown":            NewColor(139, 69, 19),
	"salmon":                  NewColor(250, 128, 114),
	"sandy_brown":             NewColor(244, 164, 96),
	"sea_green":               NewColor(46, 139, 87),
	"sea_shell":               NewColor(255, 245, 238),
	"sienna":                  NewColor(160, 82, 45),
	"silver":                  NewColor(192, 192, 192),
	"sky_blue":                NewColor(135, 206, 235),
	"slate_blue":              NewColor(106, 90, 205),
	"slate_gray":              NewColor(112, 128, 144),
	"slate_grey":              NewColor(112, 128, 144),
	"snow":                    NewColor(255, 250, 250),
	"spring_green":            NewColor(0, 255, 127),
	"steel_blue":              NewColor(70, 130, 180),
	"tan":                     NewColor(210, 180, 140),
	"teal":                    NewColor(0, 128, 128),
	"thistle":                 NewColor(216, 191, 216),
	"tomato":                  NewColor(255, 99, 71),
	"turquoise":               NewColor(64, 224, 208),
	"violet":                  NewColor(238, 130, 238),
	"wheat":                   NewColor(245, 222, 179),
	"white":                   NewColor(255, 255, 255),
	"white_smoke":             NewColor(245, 245, 245),
	"yellow":                  NewColor(255, 255, 0),
	"yellow_green":            NewColor(154, 205, 50),
}
