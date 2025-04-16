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

func TestBDFRoundTrip(t *testing.T) {
	bdf, err := bitedutils.R2BDF(strings.NewReader(test0))
	if err != nil {
		t.Error(err)
	}
	var w strings.Builder
	bdf.BDF2W(&w)
	res := w.String()
	assert.Equal(t, test0, res)
}
