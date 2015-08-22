/**
 * wskeygen.go - Generate appropriate certs for use with
 * a webserver.
 *
 * @author R. S. Doiel, <rsdoiel@yahoo.com>
 * copyright (c) 2014
 * All rights reserved.
 * @license BSD 2-Clause License
 */
package main

import (
	"../../keygen"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

var revision = "v0.0.2"

// command line parameters that override environment variables
var (
	host    string
	port    int
	cert    string
	key     string
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

 EXAMPLES
      
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
		keygenUsage  = "Interactive tool to generate TLS certificates and keys"
		keyUsage     = "Set path to your SSL key pem file."
		certUsage    = "Set path to your SSL cert pem file."
		hostUsage    = "Set this hostname for webserver."
		portUsage    = "Set the port number to listen on."
		versionUsage = "Display the version number of ws command."
	)

	flag.BoolVar(&help, "help", false, helpUsage)
	flag.BoolVar(&help, "h", false, helpUsage)
	flag.BoolVar(&version, "version", false, versionUsage)
	flag.BoolVar(&version, "v", false, versionUsage)

	// Settable via environment
	key = defaultEnvString("WS_KEY", "")
	cert = defaultEnvString("WS_CERT", "")
	host = defaultEnvString("WS_HOST", "localhost")
	port = defaultEnvInt("WS_PORT", 8000)

	flag.StringVar(&key, "key", key, keyUsage)
	flag.StringVar(&cert, "cert", cert, certUsage)
	flag.StringVar(&host, "host", host, hostUsage)
	flag.IntVar(&port, "port", port, portUsage)
}

func main() {
	flag.Parse()
	if version == true {
		fmt.Println(revision)
		os.Exit(0)
	}
	if help == true {
		usage(0, "")
	}

	certFilename, keyFilename, err := keygen.Keygen("etc/ssl", "cert.pem", "key.pem")
	if err != nil {
		log.Fatalf("%s\n", err)
	}
	log.Printf("Wrote %s and %s\n", certFilename, keyFilename)
}
