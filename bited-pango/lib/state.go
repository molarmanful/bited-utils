package bitedpango

import (
	"strings"

	"barista.run/colors"
	"barista.run/pango"
)

type State struct {
	Root    *pango.Node
	Ptr     *pango.Node
	Content strings.Builder
}

func NewState() *State {
	root := pango.New()
	return &State{Root: root, Ptr: root}
}

func (state *State) Blank() {
	state.Ptr.AppendText(state.Content.String())
	state.Content.Reset()
	state.Ptr = state.Root
}

func (state *State) Clr(clr string) {
	state.Ptr.AppendText(state.Content.String())
	state.Content.Reset()
	state.Ptr = pango.New().Color(colors.Hex(clr))
	state.Root.Append(state.Ptr)
}
