//
// Package cli processes default environment values, log handling, version number
// and other common code blocks used to implement cli tools.
//
// @author R. S. Doiel, <rsdoiel@gmail.com>
// copyright (c) 2015 all rights reserved.
// Released under the Simplified BSD License.
// See: http://opensource.org/licenses/BSD-2-Clause
//
package cli

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const revision = "v0.0.4-alpha"

// CommandName takes a path and returns the basename of the command.
var CommandName = func(s string) string {
	return filepath.Base(s)
}

// Version info for the command
var Version = func() {
	fmt.Printf("%s version %s\n", os.Args[0], revision)
	os.Exit(0)
}

// Usage info for the command
var Usage = func(exit_code int, body string, msg string) {
	var fh = os.Stderr
	if exit_code == 0 {
		fh = os.Stdout
	}
	fmt.Fprintf(fh, `%s
 USAGE %s [options]
 %s
 OPTIONS

`, msg, CommandName(os.Args[0]), body)

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

// DefaultEnvBool returns the environment boolean value or default if not set.
func DefaultEnvBool(environmentVar string, defaultValue bool) bool {
	tmp := strings.ToLower(os.Getenv(environmentVar))
	if tmp == "true" {
		return true
	}
	if tmp == "false" {
		return false
	}
	return defaultValue
}

// DefaultEnvString returns the environment value of the string or the default if not set.
func DefaultEnvString(environmentVar string, defaultValue string) string {
	tmp := os.Getenv(environmentVar)
	if tmp != "" {
		return tmp
	}
	return defaultValue
}

// DefaultEnvInt returns the environment value of the string or the default if not set.
func DefaultEnvInt(environmentVar string, defaultValue int) int {
	tmp := os.Getenv(environmentVar)
	if tmp != "" {
		i, err := strconv.Atoi(tmp)
		if err != nil {
			Usage(1, "", environmentVar+" must be an integer.")
		}
		return i
	}
	return defaultValue
}
