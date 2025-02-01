package bitedpango

import (
	"bufio"
	"io"

	"barista.run/colors"
	"barista.run/pango"
)

func Pango(txtR io.Reader, clrR io.Reader, clrs []string) *pango.Node {
	clrsMap := make(map[rune]string)
	var ks = []rune("0123456789ABCDEF")
	for i, clr := range clrs[:min(16, len(clrs))] {
		clrsMap[ks[i]] = clr
	}

	state := NewState()
	txtScan := bufio.NewScanner(txtR)
	clrScan := bufio.NewScanner(clrR)
	first := true
	for txtScan.Scan() {
		if !first {
			state.Content.WriteRune('\n')
		}
		first = false

		txtLine := []rune(txtScan.Text())
		var clrLine []rune
		if clrScan.Scan() {
			clrLine = []rune(clrScan.Text())
		}

		for i, c := range txtLine {
			if len(clrLine) > i {
				k := clrLine[i]
				switch k {
				case '.':
					state.Blank()
				default:
					if clr, ok := clrsMap[k]; ok {
						state.Clr(clr)
					}
				}
			}
			state.Content.WriteRune(c)
		}
	}

	state.Blank()

	return state.Root
}

func BgFg(node *pango.Node, bg string, fg string) {
	node.Background(colors.Hex(bg)).Color(colors.Hex(fg))
}
