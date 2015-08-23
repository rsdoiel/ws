/**
 * slugify.go - command line utility to change phrases into
 * URL frienld and human readable text strings. E.g. "This famous poem" would become
 * "This_famous_poem".
 *
 * @author R. S. Doiel, <rsdoiel@gmail.com>
 * copyright (c) 2015 all rights reserved.
 * Released under the BSD 2-Clause license.
 * See: http://opensource.org/licenses/BSD-2-Clause
 */
package main

import (
	"../../slug"
	ver "../../version"
	"flag"
	"fmt"
	"os"
)

var usage = func(exit_code int, msg string) {
	var fh = os.Stderr
	if exit_code == 0 {
		fh = os.Stdout
	}
	cmdName := os.Args[0]

	fmt.Fprintf(fh, `%s
USAGE %s [options]

%s is a command line utility to changing phrase into URL friendly
and human readable strings. E.g. "This famous poem" would become "This_famous_poem".

EXAMPLE

    %s "The World in a Nutshell"

Would yield "The_World_in_a_Nutshell"

OPTIONS
`, msg, cmdName, cmdName, cmdName)

	flag.VisitAll(func(f *flag.Flag) {
		fmt.Fprintf(fh, "\t-%s\t\t%s\n", f.Name, f.Usage)
	})

	fmt.Fprintf(fh, `
copyright (c) 2015 all rights reserved.
Released under the BSD 2-Clause license.
See: http://opensource.org/licenses/BSD-2-Clause
`)
	os.Exit(exit_code)
}

func slugify(s string, extension string) string {
	if extension != "" {
		return slug.Slugify(s) + extension
	}
	return slug.Slugify(s)
}

func main() {
	help := false
	version := false
	extension := ""
	flag.StringVar(&extension, "e", extension, "Add an extention from the slug phrase. E.g. .html")
	flag.BoolVar(&help, "h", help, "Display this help document.")
	flag.BoolVar(&help, "help", help, "Display this help document.")
	flag.BoolVar(&version, "v", version, "Display program version.")
	flag.BoolVar(&version, "version", version, "Display program version.")
	flag.Parse()
	if help == true {
		usage(0, "")
	}
	if version == true {
		fmt.Printf("%s version %s\n", os.Args[0], ver.Revision)
		os.Exit(0)
	}

	if flag.NArg() < 1 {
		usage(1, "Missing phrase to unslugify")
	}
	for _, arg := range flag.Args() {
		fmt.Println(slugify(arg, extension))
	}
}
