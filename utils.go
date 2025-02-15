// bited-utils is a set of pipeline helpers and utilities for building fonts
// from bited BDFs.
//
// Documentation for each utility:
//
//   - [github.com/molarmanful/bited-utils/bited-build]
//   - [github.com/molarmanful/bited-utils/bited-img]
//   - [github.com/molarmanful/bited-utils/bited-clr]
//   - [github.com/molarmanful/bited-utils/bited-scale]
package bitedutils

// Check panics if error is non-nil.
func Check(err error) {
	if err != nil {
		panic(err)
	}
}
