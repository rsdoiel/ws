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
	cli "../../cli"
	"../../keygen"
	"flag"
	"fmt"
	"log"
	"os"
)

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
	key = cli.DefaultEnvString("WS_KEY", "")
	cert = cli.DefaultEnvString("WS_CERT", "")
	host = cli.DefaultEnvString("WS_HOST", "localhost")
	port = cli.DefaultEnvInt("WS_PORT", 8000)

	flag.StringVar(&key, "key", key, keyUsage)
	flag.StringVar(&cert, "cert", cert, certUsage)
	flag.StringVar(&host, "host", host, hostUsage)
	flag.IntVar(&port, "port", port, portUsage)
}

func main() {
	usageDescription := fmt.Sprintf(`
 %s is an interactive program that generates SSL/TLS certificates for use with
 the ws and wsjs web servers.

`, cli.CommandName(os.Args[0]))

	flag.Parse()
	if help == true {
		cli.Usage(0, usageDescription, "")
	}
	if version == true {
		cli.Version()
	}

	certFilename, keyFilename, err := keygen.Keygen("etc/ssl", "cert.pem", "key.pem")
	if err != nil {
		log.Fatalf("%s\n", err)
	}
	log.Printf("Wrote %s and %s\n", certFilename, keyFilename)
}
