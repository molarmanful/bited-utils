package bitedimg

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
)

var reComb = regexp.MustCompile(`\pM`)

func (unit *Unit) GenChars() error {
	log.Println("+ GEN chars")
	charsF, err := os.Create(filepath.Join(unit.TxtDir, unit.Chars.Out+".txt"))
	if err != nil {
		return err
	}
	defer charsF.Close()

	var first = true
	for ns := range slices.Chunk(unit.Codes, unit.Chars.Width) {
		if !first {
			if _, err := fmt.Fprintln(charsF); err != nil {
				return err
			}
		}

		first = false
		for i, n := range ns {
			if i > 0 {
				if _, err := fmt.Fprint(charsF, " "); err != nil {
					return err
				}
			}

			char := string(rune(n))
			if unit.HideAccents {
				char = reComb.ReplaceAllLiteralString(char, ".")
			}
			if _, err := fmt.Fprint(charsF, char); err != nil {
				return err
			}
		}
	}

	return nil
}

func (unit *Unit) GenMap() error {
	log.Println("+ GEN map")
	mapF, err := os.Create(filepath.Join(unit.TxtDir, unit.Map.Out+".txt"))
	if err != nil {
		return err
	}
	defer mapF.Close()
	mapClrF, err := os.Create(filepath.Join(unit.TxtDir, unit.Map.Out+".clr"))
	if err != nil {
		return err
	}
	defer mapClrF.Close()

	if _, err := fmt.Fprint(mapF,
		"          0 1 2 3 4 5 6 7 8 9 A B C D E F\n",
		"        ┌────────────────────────────────",
	); err != nil {
		return err
	}
	if _, err := fmt.Fprint(mapClrF, unit.Map.LabelClrs[0], "\n", unit.Map.BorderClr); err != nil {
		return err
	}

	row := -1
	var line []string
	clrLine := fmt.Sprintf("%s     %s %s.", unit.Map.LabelClrs[1], unit.Map.LabelClrs[0], unit.Map.BorderClr)
	for _, n := range unit.Codes {
		if n/16 != row {
			row = n / 16
			if _, err := fmt.Fprintln(mapF, strings.Join(line, " ")); err != nil {
				return err
			}
			line = strings.Split(strings.Repeat(" ", 16), "")
			if _, err := fmt.Fprintf(mapF, "U+%04X_ │ ", row); err != nil {
				return err
			}
			if _, err := fmt.Fprint(mapClrF, "\n", clrLine); err != nil {
				return err
			}
		}

		char := string(rune(n))
		if unit.HideAccents {
			char = reComb.ReplaceAllLiteralString(char, ".")
		}
		line[n%16] = char
	}
	if _, err := fmt.Fprint(mapF, strings.Join(line, " ")); err != nil {
		return err
	}

	return nil
}
