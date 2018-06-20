package img

import (
	"image/color"
	stdErr "errors"
	"regexp"
	"strconv"
	"github.com/pkg/errors"
)

var (
	ErrorEmptySize     = stdErr.New("size cannot be empty")
	ErrorMalformedSize = stdErr.New("malformed size")
	ErrorZeroValue     = stdErr.New("value is 0")
	ErrorUnexpected    = stdErr.New("unexpected error")
)

var (
	sizePattern = regexp.MustCompile(`^(\d+)x(\d+)$`)
)

func parseSize(s string) (int, int, error) {
	if s == "" {
		return 0, 0, ErrorEmptySize
	}

	ss := sizePattern.FindStringSubmatch(s)
	if len(ss) != 3 {
		return 0, 0, ErrorMalformedSize
	}

	width, err := strconv.Atoi(ss[1])
	if err != nil {
		return 0, 0, errors.Wrap(ErrorUnexpected, "bad width")
	}
	if width == 0 {
		return 0, 0, errors.Wrap(ErrorZeroValue, "width is too small")
	}

	height, err := strconv.Atoi(ss[2])
	if err != nil {
		return 0, 0, errors.Wrap(ErrorUnexpected, "bad height")
	}
	if height == 0 {
		return 0, 0, errors.Wrap(ErrorZeroValue, "height is too small")
	}

	return width, height, nil
}

func parseColor(s string) (color.RGBA, error) {
	return color.RGBA{}, nil
}

func parseText(s string) (string, error) {
	return "", nil
}
