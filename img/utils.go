package img

import (
	stdErr "errors"
	"image/color"
	"regexp"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

var (
	errorMissing    = stdErr.New("missing value")
	errorMalformed  = stdErr.New("malformed")
	errorUnexpected = stdErr.New("unexpected error")
)

var (
	sizePattern  = regexp.MustCompile(`^([1-9][0-9]{0,5})x([1-9][0-9]{0,5})$`)
	colorPattern = regexp.MustCompile(`^(?:[0-9A-Fa-f]{3}){1,2}$`)
)

func parseSize(s string) (int, int, error) {
	if s == "" {
		return 0, 0, errors.Wrap(errorMissing, "size is empty")
	}

	ss := sizePattern.FindStringSubmatch(s)
	if len(ss) != 3 {
		return 0, 0, errors.Wrapf(errorMalformed, "size '%s' is not valid", s)
	}

	width, err := strconv.Atoi(ss[1])
	if err != nil {
		return 0, 0, errors.Wrapf(errorUnexpected, "bad width '%s'", ss[1])
	}

	height, err := strconv.Atoi(ss[2])
	if err != nil {
		return 0, 0, errors.Wrapf(errorUnexpected, "bad height '%s'", ss[2])
	}

	return width, height, nil
}

func parseColor(s string) (*color.RGBA, error) {
	if s == "" {
		return nil, errors.Wrap(errorMissing, "color is empty")
	}

	if !colorPattern.MatchString(s) {
		return nil, errors.Wrapf(errorMalformed, "color '%s' is not valid", s)
	}

	var ss []string
	if n := len(s); n == 3 {
		ss = strings.Split(s, "")
		ss[0] += ss[0]
		ss[1] += ss[1]
		ss[2] += ss[2]
	} else if n == 6 {
		ss = []string{s[0:2], s[2:4], s[4:6]}
	} else {
		return nil, errors.Wrapf(errorUnexpected, "bad color '%s'", s)
	}

	red, err := strconv.ParseUint(ss[0], 16, 8)
	if err != nil {
		return nil, errors.Wrapf(errorUnexpected, "bad red '%s'", ss[0])
	}

	green, err := strconv.ParseUint(ss[1], 16, 8)
	if err != nil {
		return nil, errors.Wrapf(errorUnexpected, "bad green '%s'", ss[1])
	}

	blue, err := strconv.ParseUint(ss[2], 16, 8)
	if err != nil {
		return nil, errors.Wrapf(errorUnexpected, "bad blue '%s'", ss[2])
	}

	return &color.RGBA{
		R: uint8(red),
		G: uint8(green),
		B: uint8(blue),
		A: 255,
	}, nil
}

func parseText(s string) (string, error) {
	return "", nil
}
