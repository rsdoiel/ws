/**
 * unslugify.go - command line utility to change slug text phrases into
 * readable text strings. E.g. "This_famous_poem" would become
 * "This famous poem".
 *
 * @author R. S. Doiel, <rsdoiel@gmail.com>
 * copyright (c) 2015 all rights reserved.
 * Released under the BSD 2-Clause license.
 * See: http://opensource.org/licenses/BSD-2-Clause
 */
package main

import (
	"../../slug"
	"flag"
	"fmt"
	"os"
	"strings"
)

var usage = func(exit_code int, msg string) {
	var fh = os.Stderr
	if exit_code == 0 {
		fh = os.Stdout
	}
	cmdName := os.Args[0]

	fmt.Fprintf(fh, `%s
USAGE %s [options]

%s is a command line utility to changing slug phrases into friendly
strings. E.g. "This_famous_poem" would become "This famous poem".

EXAMPLE

    %s "The_World_in_a_Nutshell"

Would yield "The World in a Nutshell"

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

func unslugify(s string, extension string) string {
	if extension != "" {
		return slug.Unslugify(strings.TrimSuffix(s, extension))
	}
	return slug.Unslugify(s)
}

func main() {
	help := false
	extension := ""
	flag.StringVar(&extension, "e", extension, "Remove the extention from the slug phrase. E.g. .html")
	flag.BoolVar(&help, "h", help, "Display this help document.")
	flag.BoolVar(&help, "help", help, "Display this help document.")
	flag.Parse()
	if help == true {
		usage(0, "")
	}

	if flag.NArg() < 1 {
		usage(1, "Missing phrase to unslugify")
	}
	for _, arg := range flag.Args() {
		fmt.Println(unslugify(arg, extension))
	}
}
