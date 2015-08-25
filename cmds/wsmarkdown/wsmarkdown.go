/**
 * wsmarkdown.go - A command line markdown processor wrapping Blackfriday
 *
 * @author R. S. Doiel, <rsdoiel@gmail.com>
 * copyright (c) 2015
 * All rights reserved.
 * @license BSD 2-Clause License
 */
package main

import (
	cli "../../src/cli"
	"flag"
	"fmt"
	//"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
	"io/ioutil"
	"log"
	"os"
)

// command line parameters that override environment variables
var (
	version bool
	help    bool
)

type stringValue string

func init() {
	var (
		helpUsage    = "This help document."
		versionUsage = fmt.Sprintf("Display the version number for %s.", cli.CommandName(os.Args[0]))
	)

	flag.BoolVar(&help, "help", false, helpUsage)
	flag.BoolVar(&help, "h", false, helpUsage)
	flag.BoolVar(&version, "version", false, versionUsage)
	flag.BoolVar(&version, "v", false, versionUsage)
}

func main() {
	usageDescription := fmt.Sprintf(`
 Process Markdown input with Blackfriday and Bluemonday and generate HTML output.

 %s input from stdin and writes to stnout.

`, cli.CommandName(os.Args[0]))

	flag.Parse()
	if version == true {
		cli.Version()
	}
	if help == true {
		cli.Usage(0, usageDescription, "")
	}
	src, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Printf("%s", bluemonday.UGCPolicy().SanitizeBytes(blackfriday.MarkdownCommon(src)))
	fmt.Printf("%s", blackfriday.MarkdownCommon(src))
}
