/**
 * wsjs.go - A light weight webserver for static content
 * and dynamic content calculated via JavaScript files. It is
 * intended for development and prototyping route based web API.
 *
 * Supports both http and https protocols. Dynamic route
 * processing available via Otto JavaScript virtual machines.
 *
 * @author R. S. Doiel, <rsdoiel@yahoo.com>
 * copyright (c) 2014
 * All rights reserved.
 * @license BSD 2-Clause License
 */
package main

import (
	"../../src/cfg"
	cli "../../src/cli"
	"../../src/fsengine"
	"../../src/ottoengine"
	"../../src/wslog"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

// command line parameters that override environment variables
var (
	useTLS   bool
	docroot  string
	host     string
	port     int
	cert     string
	key      string
	otto     bool
	ottoPath string
	version  bool
	doKeygen bool
	doInit   bool
	help     bool
)

type stringValue string

func webserver(config *cfg.Cfg) error {
	// If otto is enabled add routes and handle them.
	if config.Otto == true {
		ottoPath, err := filepath.Abs(config.OttoPath)
		if err != nil {
			log.Fatalf("Can't read %s: %s\n", config.OttoPath, err)
		}
		programs, err := ottoengine.Load(ottoPath)
		if err != nil {
			log.Fatalf("Load error: %s\n", err)
		}
		ottoengine.AddRoutes(programs)
	}

	// Restricted FileService excluding dot files and directories
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// hande off this request/response pair to the fsengine
		fsengine.Engine(config, w, r)
	})

	// Now start up the server and log transactions
	if config.UseTLS == true {
		if config.Cert == "" || config.Key == "" {
			log.Fatalf("TLS set true but missing key or certificate")
		}
		log.Println("Starting https://" + net.JoinHostPort(config.Hostname, strconv.Itoa(config.Port)))
		return http.ListenAndServeTLS(net.JoinHostPort(config.Hostname, strconv.Itoa(config.Port)), config.Cert, config.Key, wslog.RequestLog(http.DefaultServeMux))
	}
	log.Println("Starting http://" + net.JoinHostPort(config.Hostname, strconv.Itoa(config.Port)))
	// Now start up the server and log transactions
	return http.ListenAndServe(net.JoinHostPort(config.Hostname, strconv.Itoa(config.Port)), wslog.RequestLog(http.DefaultServeMux))
}

func init() {
	const (
		helpUsage     = "This help document."
		useTLSUsage   = "When true this turns on TLS (https) support."
		keyUsage      = "Path to your SSL key pem file."
		certUsage     = "path to your SSL cert pem file."
		docrootUsage  = "This is your document root for static files."
		hostUsage     = "Set this hostname for webserver."
		portUsage     = "Set the port number to listen on."
		ottoPathUsage = "Turns on otto engine using the path for route JavaScript route handlers"
		versionUsage  = "Display the version number of ws command."
	)

	flag.BoolVar(&help, "help", false, helpUsage)
	flag.BoolVar(&help, "h", false, helpUsage)
	flag.BoolVar(&version, "version", false, versionUsage)
	flag.BoolVar(&version, "v", false, versionUsage)

	// Settable via environment
	useTLS = cli.DefaultEnvBool("WS_TLS", false)
	key = cli.DefaultEnvString("WS_KEY", "")
	cert = cli.DefaultEnvString("WS_CERT", "")
	docroot = cli.DefaultEnvString("WS_DOCROOT", "")
	otto = cli.DefaultEnvBool("WS_OTTO", true)
	ottoPath = cli.DefaultEnvString("WS_OTTO_PATH", "")
	host = cli.DefaultEnvString("WS_HOST", "localhost")
	port = cli.DefaultEnvInt("WS_PORT", 8000)
	flag.BoolVar(&useTLS, "tls", useTLS, useTLSUsage)
	flag.StringVar(&key, "key", key, keyUsage)
	flag.StringVar(&cert, "cert", cert, certUsage)
	flag.StringVar(&docroot, "docroot", docroot, docrootUsage)
	flag.StringVar(&docroot, "d", docroot, docrootUsage)
	flag.StringVar(&host, "host", host, hostUsage)
	flag.StringVar(&host, "H", host, hostUsage)
	flag.IntVar(&port, "port", port, portUsage)
	flag.IntVar(&port, "p", port, portUsage)
	flag.StringVar(&ottoPath, "otto-path", ottoPath, ottoPathUsage)
	flag.StringVar(&ottoPath, "o", ottoPath, ottoPathUsage)
}

func main() {
	usageDescription := fmt.Sprintf(`
 %s is a simple static content web server with support for 
 defining additional routes via a separate directory of JavaScript files.
 It is suitable for prototyping simple web APIs.

`, cli.CommandName(os.Args[0]))

	flag.Parse()
	if version == true {
		cli.Version()
	}
	if help == true {
		cli.Usage(0, usageDescription, "")
	}

	config, err := cfg.Configure(docroot, host, port, useTLS, cert, key, otto, ottoPath)
	if err != nil {
		cli.Usage(1, usageDescription, fmt.Sprintf("%s", err))
	}
	if ottoPath == "" {
		otto = false
	} else {
		otto = true
	}

	fmt.Printf("\n\n"+
		"          TLS: %t\n"+
		"         Cert: %s\n"+
		"          Key: %s\n"+
		"      Docroot: %s\n"+
		"         Host: %s\n"+
		"         Port: %d\n"+
		"       Run as: %s\n\n"+
		" Otto enabled: %t\n"+
		"         Path: %s\n"+
		"\n\n",
		config.UseTLS,
		config.Cert,
		config.Key,
		config.Docroot,
		config.Hostname,
		config.Port,
		config.Username,
		config.Otto,
		config.OttoPath)
	err = webserver(config)
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(0)
}
