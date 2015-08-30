/**
 * wsinit.go - Initialize a project setup for use with the ws webserver.
 *
 * @author R. S. Doiel, <rsdoiel@yahoo.com>
 * copyright (c) 2015
 * All rights reserved.
 * @license BSD 2-Clause License
 */
package main

import (
	"../../cfg"
	cli "../../cli"
	"flag"
	"log"
)

// command line parameters that override environment variables
var (
	version bool
	help    bool
)

type stringValue string

func init() {
	const (
		helpUsage    = "This help document."
		versionUsage = "Display the version number of wsinit command."
	)

	flag.BoolVar(&help, "help", false, helpUsage)
	flag.BoolVar(&help, "h", false, helpUsage)
	flag.BoolVar(&version, "version", false, versionUsage)
	flag.BoolVar(&version, "v", false, versionUsage)
}

func main() {
	usageDescription := `
 Create a basic project structure in the current working 
 directory and generates SSL/TLS certificates if needed.

`

	flag.Parse()
	if version == true {
		cli.Version()
	}
	if help == true {
		cli.Usage(0, usageDescription, "")
	}

	err := cfg.InitProject()
	if err != nil {
		log.Fatalf("%s\n", err)
	}
	log.Println("Project initialization complete.")
}
