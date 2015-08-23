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
	ver "../../version"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

// command line parameters that override environment variables
var (
	version bool
	help    bool
)

type stringValue string

var usage = func(exit_code int, msg string) {
	var fh = os.Stderr
	if exit_code == 0 {
		fh = os.Stdout
	}
	fmt.Fprintf(fh, `%s
 USAGE %s [options]

		Create a basic project structure in the current working 
		directory and generates SSL/TLS certificates if needed.

 OPTIONS

`, msg, os.Args[0])

	flag.VisitAll(func(f *flag.Flag) {
		fmt.Fprintf(fh, "\t-%s\t(defaults to %s) %s\n", f.Name, f.DefValue, f.Usage)
	})

	fmt.Fprintf(fh, `

 copyright (c) 2014 all rights reserved.
 Released under the Simplified BSD License
 See: http://opensource.org/licenses/bsd-license.php

`)
	os.Exit(exit_code)
}

func defaultEnvBool(environmentVar string, defaultValue bool) bool {
	tmp := strings.ToLower(os.Getenv(environmentVar))
	if tmp == "true" {
		return true
	}
	if tmp == "false" {
		return false
	}
	return defaultValue
}

func defaultEnvString(environmentVar string, defaultValue string) string {
	tmp := os.Getenv(environmentVar)
	if tmp != "" {
		return tmp
	}
	return defaultValue
}

func defaultEnvInt(environmentVar string, defaultValue int) int {
	tmp := os.Getenv(environmentVar)
	if tmp != "" {
		i, err := strconv.Atoi(tmp)
		if err != nil {
			usage(1, environmentVar+" must be an integer.")
		}
		return i
	}
	return defaultValue
}

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
	flag.Parse()
	if version == true {
		fmt.Printf("%s version %s\n", os.Args[0], ver.Revision)
		os.Exit(0)
	}
	if help == true {
		usage(0, "")
	}

	err := cfg.InitProject()
	if err != nil {
		log.Fatalf("%s\n", err)
	}
	log.Println("Project initialization complete.")
}
