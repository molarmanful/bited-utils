package bitedutils_test

import (
	_ "embed"
	"strings"
	"testing"

	bitedutils "github.com/molarmanful/bited-utils"
	"github.com/stretchr/testify/assert"
)

//go:embed test0.bdf
var test0 string

//go:embed test0_2x.bdf
var test0_2x string

//go:embed test0_3x.bdf
var test0_3x string

func TestBDFRoundTrip(t *testing.T) {
	bdf, err := bitedutils.R2BDF(strings.NewReader(test0))
	if err != nil {
		t.Error(err)
	}
	var w strings.Builder
	if err := bdf.BDF2W(&w); err != nil {
		t.Error(err)
	}
	res := w.String()
	assert.Equal(t, test0, res)
}

func TestBDFScale1(t *testing.T) {
	for scale, expected := range map[int]string{
		1: test0,
		2: test0_2x,
		3: test0_3x,
	} {
		bdf, err := bitedutils.R2BDF(strings.NewReader(test0))
		if err != nil {
			t.Error(err)
		}
		bdf.Scale(scale)
		var w strings.Builder
		if err := bdf.BDF2W(&w); err != nil {
			t.Error(err)
		}
		res := w.String()
		assert.Equal(t, expected, res)
	}
}
