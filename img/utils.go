package img

import (
	stdErr "errors"
	"github.com/pkg/errors"
	"image/color"
	"regexp"
	"strconv"
)

var (
	ErrorEmptySize     = stdErr.New("size is empty")
	ErrorMalformedSize = stdErr.New("malformed size")
	ErrorUnexpected    = stdErr.New("unexpected error")
)

var (
	sizePattern = regexp.MustCompile(`^([1-9][0-9]*)x([1-9][0-9]*)$`)
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

	height, err := strconv.Atoi(ss[2])
	if err != nil {
		return 0, 0, errors.Wrap(ErrorUnexpected, "bad height")
	}

	return width, height, nil
}

func parseColor(s string) (color.RGBA, error) {
	return color.RGBA{}, nil
}

func parseText(s string) (string, error) {
	return "", nil
}
