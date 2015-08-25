/**
 * slugify.go - implements tests for slug Go package.
 * @author R. S. Doiel, <rsdoiel@gmail.com>
 * copyright (c) 2015 all rights reserved.
 * Released under the BSD 2-Clause license
 * See: http://opensource.org/licenses/BSD-2-Clause
 */
package slugify

import (
	"../ok"
	"testing"
)

func TestSlugify(t *testing.T) {
	original := "Hello World"
	expected := "Hello_World"
	result := Slugify(original)
	ok.Ok(t, expected == result, result+" is invalid.")
}

func TestUnslugify(t *testing.T) {
	original := "Hello_World"
	expected := "Hello World"
	result := Unslugify(original)
	ok.Ok(t, expected == result, result+" is invalid.")
}
