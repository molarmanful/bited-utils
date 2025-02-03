package main

import (
	"os"
	"strings"

	"github.com/gdamore/tcell/v2"
	bitedutils "github.com/molarmanful/bited-utils"
	bitedimg "github.com/molarmanful/bited-utils/bited-img/lib"
	"github.com/rivo/tview"
)

type State struct {
	View   *tview.TextView
	Txt    [][]rune
	Clr    [][]rune
	ClrF   *os.File
	Clrs   bitedimg.Clrs
	ClrMap map[rune]string
	X      int
	Y      int
}

func (state *State) MkClrs(clrs bitedimg.Clrs) {
	state.Clrs = clrs
	state.ClrMap = make(map[rune]string)
	state.ClrMap['.'] = "[" + clrs.Fg + "]"
	for i, r := range "0123456789ABCDEF" {
		state.ClrMap[r] = "[" + state.Clrs.Base[i] + "]"
	}
}

func (state *State) MkView(app *tview.Application) {
	state.View = tview.NewTextView().
		SetDynamicColors(true).SetRegions(true).
		SetChangedFunc(func() {
			app.Draw()
		})

	state.View.
		SetBackgroundColor(tcell.GetColor(state.Clrs.Bg)).
		SetBorder(true).SetBorderColor(tcell.ColorGray).
		SetTitleColor(tcell.ColorGray)

	state.View.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		ec := event.Rune()
		if _, ok := state.ClrMap[ec]; ok {
			state.ClrC(ec)
			return nil
		}
		switch event.Rune() {
		case 'h':
			state.MoveC(-1, 0)
			return nil
		case 'j':
			state.MoveC(0, 1)
			return nil
		case 'k':
			state.MoveC(0, -1)
			return nil
		case 'l':
			state.MoveC(1, 0)
			return nil
		}
		return event
	})

	state.Gen()
}

func (state *State) Gen() {
	var res strings.Builder
	var ir strings.Builder
	for y, l := range state.Txt {
		for x, c := range l {
			if ca, ok := state.ClrAt(x, y); ok {
				res.WriteString(tview.Escape(ir.String()))
				ir.Reset()
				res.WriteString(state.ClrMap[ca])
			}
			isC := x == state.X && y == state.Y
			if isC {
				res.WriteString(tview.Escape(ir.String()))
				ir.Reset()
				res.WriteString(`["c"]`)
				res.WriteRune(c)
				res.WriteString(`[""]`)
			} else {
				ir.WriteRune(c)
			}
		}
		ir.WriteRune('\n')
	}

	res.WriteString(tview.Escape(ir.String()))
	state.View.SetText(res.String()).Highlight("c").ScrollToHighlight()
}

func (state *State) MoveC(x int, y int) {
	y0 := state.Y
	state.X = max(0, min(len(state.Txt[y0])-1, state.X+x))
	state.Y = max(0, min(len(state.Txt)-1, state.Y+y))
	state.Gen()
}

func (state *State) ClrC(c rune) {
	state.SetClr(c, state.X, state.Y)
}

func (state *State) ClrAt(x int, y int) (rune, bool) {
	if len(state.Clr) > y && len(state.Clr[y]) > x && state.Clr[y][x] != ' ' {
		return state.Clr[y][x], true
	}
	return rune(0), false
}

func (state *State) SetClr(c rune, x int, y int) {
	if len(state.Clr) <= y {
		state.Clr = append(state.Clr, make([][]rune, y-len(state.Clr)+1)...)
	}
	if len(state.Clr[y]) <= x {
		xs := make([]rune, x-len(state.Clr[y])+1)
		for i := range xs {
			xs[i] = ' '
		}
		state.Clr[y] = append(state.Clr[y], xs...)
	}
	state.Clr[y][x] = c
	state.WriteClr()
	state.Gen()
}

func (state *State) WriteClr() {
	var res strings.Builder
	var ir strings.Builder
	clr := '.'
	for _, l := range state.Clr {
		for _, c := range l {
			if c == clr {
				ir.WriteRune(' ')
				continue
			}
			if _, ok := state.ClrMap[c]; !ok {
				ir.WriteRune(' ')
				continue
			}
			clr = c
			ir.WriteRune(c)
		}
		res.WriteString(strings.TrimRight(ir.String(), " "))
		ir.Reset()
		res.WriteRune('\n')
	}

	err := state.ClrF.Truncate(0)
	bitedutils.Check(err)
	_, err = state.ClrF.Seek(0, 0)
	bitedutils.Check(err)
	_, err = state.ClrF.WriteString(strings.TrimRight(res.String(), "\n"))
	bitedutils.Check(err)
}
