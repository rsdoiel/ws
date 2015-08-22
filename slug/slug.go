/**
 * slug.go - A library to convert strings to and from slug form.
 * slug forms are often used to generate strings that are URL friendly
 * with out resorting to URL encoding. This leaves them more human readable.
 * The downside is that the transformation is loosy.  Generally the approach
 * is to replace spaces, underscores and other problematic with underscores. Since more
 * than the underscore becomes overloaded reversing the map is problematic.
 * Thus Slugify and Unslugify only map to each other directly in the simplest cases.
 *
 * @author R. S. Doiel, <rsdoiel@gmail.com>
 * copyright (c) 2015 all rights reserved.
 * Released under the BSD 2-Clause license.
 * See: http://opensource.org/licenses/BSD-2-Clause
 */
package slug

import (
	"strings"
)

func Slugify(s string) string {
	return strings.Replace(s, " ", "_", -1)
}

func Unslugify(s string) string {
	return strings.Replace(s, "_", " ", -1)
}
