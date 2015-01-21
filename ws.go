/**
 * ws.go - A light weight webserver for static content
 * development and prototyping route based web API.
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
	"./cfg"
	"./fsengine"
	"./keygen"
	"./ottoengine"
	"./wslog"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var revision = "v0.0.0-alpha"

// command line parameters that override environment variables
var (
	useTLS   bool
	docroot   string
	host      string
	port      int
	cert      string
	key       string
	otto      bool
	ottoPath string
	version   bool
	doKeygen bool
	doInit   bool
	help      bool
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

func requestLog(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wslog.LogRequest(r.Method, r.URL, r.RemoteAddr, r.Proto, r.Referer(), r.UserAgent())
		handler.ServeHTTP(w, r)
	})
}

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
		return http.ListenAndServeTLS(net.JoinHostPort(config.Hostname, strconv.Itoa(config.Port)), config.Cert, config.Key, requestLog(http.DefaultServeMux))
	}
	log.Println("Starting http://" + net.JoinHostPort(config.Hostname, strconv.Itoa(config.Port)))
	// Now start up the server and log transactions
	return http.ListenAndServe(net.JoinHostPort(config.Hostname, strconv.Itoa(config.Port)), requestLog(http.DefaultServeMux))
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
		helpUsage      = "This help document."
		keygenUsage    = "Interactive tool to generate TLS certificates and keys"
		initUsage      = "Creates a basic project structure in the current working directory"
		useTLSUsage   = "When true this turns on TLS (https) support."
		keyUsage       = "Path to your SSL key pem file."
		certUsage      = "path to your SSL cert pem file."
		docrootUsage   = "This is your document root for static files."
		hostUsage      = "Set this hostname for webserver."
		portUsage      = "Set the port number to listen on."
		ottoUsage      = "When true this option turns on ottoengine. Uses the path defined by WS_OTTO_PATH environment variable or one provided by -O option."
		ottoPathUsage = "Turns on otto engine using the path for route JavaScript route handlers"
		versionUsage   = "Display the version number of ws command."
	)

	flag.BoolVar(&help, "help", false, helpUsage)
	flag.BoolVar(&help, "h", false, helpUsage)
	flag.BoolVar(&doKeygen, "keygen", false, keygenUsage)
	flag.BoolVar(&doInit, "init", false, initUsage)
	flag.BoolVar(&version, "version", false, versionUsage)
	flag.BoolVar(&version, "v", false, versionUsage)

	// Settable via environment
	useTLS = defaultEnvBool("WS_TLS", false)
	key = defaultEnvString("WS_KEY", "")
	cert = defaultEnvString("WS_CERT", "")
	docroot = defaultEnvString("WS_DOCROOT", "")
	otto = defaultEnvBool("WS_OTTO", false)
	ottoPath = defaultEnvString("WS_OTTO_PATH", "")
	host = defaultEnvString("WS_HOST", "localhost")
	port = defaultEnvInt("WS_PORT", 8000)

	flag.BoolVar(&useTLS, "tls", useTLS, useTLSUsage)
	flag.StringVar(&key, "key", key, keyUsage)
	flag.StringVar(&cert, "cert", cert, certUsage)
	flag.StringVar(&docroot, "docroot", docroot, docrootUsage)
	flag.StringVar(&docroot, "D", docroot, docrootUsage)
	flag.StringVar(&host, "host", host, hostUsage)
	flag.StringVar(&host, "H", host, hostUsage)
	flag.IntVar(&port, "port", port, portUsage)
	flag.IntVar(&port, "P", port, portUsage)
	flag.BoolVar(&otto, "otto", otto, ottoUsage)
	flag.BoolVar(&otto, "o", otto, ottoUsage)
	flag.StringVar(&ottoPath, "otto-path", ottoPath, ottoPathUsage)
	flag.StringVar(&ottoPath, "O", ottoPath, ottoPathUsage)
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

	config, err := cfg.Configure(docroot, host, port, useTLS, cert, key, otto, ottoPath)
	if err != nil {
		usage(1, fmt.Sprintf("%s", err))
	}

	if doKeygen == true {
		certFilename, keyFilename, err := keygen.Keygen("etc/ssl", "cert.pem", "key.pem")
		if err != nil {
			log.Fatalf("%s\n", err)
		}
		fmt.Printf("Wrote %s and %s\n", certFilename, keyFilename)
		os.Exit(0)
	}

	if doInit == true {
		err := cfg.InitProject()
		if err != nil {
			log.Fatalf("%s\n", err)
		}
		os.Exit(0)
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
