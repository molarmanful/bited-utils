package bitedutils

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/bitfield/script"
)

func Check(err error) {
	if err != nil {
		panic(err)
	}
}

var reXLFD = regexp.MustCompile(`^\s*FONT\s+`)

func GetFsz(path string) (int, error) {
	xlfd, err := script.File(path).MatchRegexp(reXLFD).First(1).String()
	if err != nil {
		return 0, err
	}
	xlfields := strings.Split(xlfd, "-")
	if len(xlfields) < 8 {
		return 0, fmt.Errorf("bad XLFD: %q", xlfd)
	}
	return strconv.Atoi(xlfields[7])
}
