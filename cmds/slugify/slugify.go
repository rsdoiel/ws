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
	cli "../../src/cli"
	slug "../../src/slugify"
	"flag"
	"fmt"
	"os"
)

func slugify(s string, extension string) string {
	if extension != "" {
		return slug.Slugify(s) + extension
	}
	return slug.Slugify(s)
}

func main() {
	usageDescription := fmt.Sprintf(`
 %s is a command line utility to changing phrase into URL friendly
 and human readable strings. E.g. "This famous poem" would become "This_famous_poem".

 EXAMPLE

    %s "The World in a Nutshell"

 Would yield "The_World_in_a_Nutshell"

`, cli.CommandName(os.Args[0]), cli.CommandName(os.Args[0]))

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
		cli.Usage(0, usageDescription, "")
	}
	if version == true {
		cli.Version()
	}

	if flag.NArg() < 1 {
		cli.Usage(1, usageDescription, "Missing phrase to unslugify")
	}
	for _, arg := range flag.Args() {
		fmt.Println(slugify(arg, extension))
	}
}
