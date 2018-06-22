package img

import (
	stdErr "errors"
	"github.com/pkg/errors"
	"image/color"
	"regexp"
	"strconv"
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
		return 0, 0, errors.Wrap(errorMissing, "color is empty")
	}

	ss := sizePattern.FindStringSubmatch(s)
	if len(ss) != 3 {
		return 0, 0, errors.Wrapf(errorMalformed, "size '%s' is not valid", s)
	}

	width, err := strconv.Atoi(ss[1])
	if err != nil { // should never happen
		return 0, 0, errors.Wrap(errorUnexpected, "bad width")
	}

	height, err := strconv.Atoi(ss[2])
	if err != nil { // should never happen
		return 0, 0, errors.Wrap(errorUnexpected, "bad height")
	}

	return width, height, nil
}

func parseColor(s string) (color.RGBA, error) {
	return color.RGBA{}, nil
}

func parseText(s string) (string, error) {
	return "", nil
}
