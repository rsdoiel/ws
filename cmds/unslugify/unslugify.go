/**
 * unslugify.go - command line utility to change slugify text phrases into
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
	cli "../../src/cli"
	slug "../../src/slugify"
	"flag"
	"fmt"
	"os"
	"strings"
)

func unslugify(s string, extension string) string {
	if extension != "" {
		return slug.Unslugify(strings.TrimSuffix(s, extension))
	}
	return slug.Unslugify(s)
}

func main() {
	usageDescription := fmt.Sprintf(`
 %s is a command line utility to changing URL friendly
 and human readable string backinto a phrase. E.g. 
 "This_famous_poem" would become "This famous poem".

 EXAMPLE

    %s "The_World_in_a_Nutshell"

 Would yield "The World in a Nutshell"

`, cli.CommandName(os.Args[0]), cli.CommandName(os.Args[0]))

	help := false
	version := false
	extension := ""
	flag.StringVar(&extension, "e", extension, "Remove the extention from the slug phrase. E.g. .html")
	flag.BoolVar(&help, "h", help, "Display this help document.")
	flag.BoolVar(&help, "help", help, "Display this help document.")
	flag.BoolVar(&version, "v", version, "Display the version number.")
	flag.BoolVar(&version, "version", version, "Display the version number.")
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
		fmt.Println(unslugify(arg, extension))
	}
}
