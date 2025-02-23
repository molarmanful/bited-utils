package bitedbuild

import (
	"os"
	"regexp"
	"strings"

	"github.com/bitfield/script"
)

var reFONT = regexp.MustCompile(`^FONT (-[^-]*-)[^-]*(-.*$)`)
var reFAM = regexp.MustCompile(`^FAMILY_NAME .*$`)

// renameBDF writes a renamed BDF copy to a tempfile.
func (unit *Unit) mkRenamedBDF(name string) (string, error) {
	tmp, err := os.CreateTemp("", "*.bdf")
	if err != nil {
		return "", err
	}
	defer tmp.Close()
	if _, err := script.File(unit.Src).
		ReplaceRegexp(reFONT, `FONT $1`+name+`$2`).
		ReplaceRegexp(reFAM, `FAMILY_NAME "`+strings.ReplaceAll(name, `"`, `""`)+`"`).
		WithStdout(tmp).Stdout(); err != nil {
		os.Remove(tmp.Name())
		return "", err
	}
	return tmp.Name(), nil
}
