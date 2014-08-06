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

var REVISION = "v0.0.0-alpha"

// command line parameters that override environment variables
var (
	use_tls   bool
	docroot   string
	host      string
	port      int
	cert      string
	key       string
	otto      bool
	otto_path string
	version   bool
	do_keygen bool
	do_init   bool
	help      bool
)

type stringValue string

var Usage = func(exit_code int, msg string) {
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
		fmt.Fprintf(fh, "\t-%s\t(defaults to %s) %s\n", f.Name, f.Value, f.Usage)
	})

	fmt.Fprintf(fh, `

 copyright (c) 2014 all rights reserved.
 Released under the Simplified BSD License
 See: http://opensource.org/licenses/bsd-license.php

`)
	os.Exit(exit_code)
}

func request_log(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wslog.LogRequest(r.Method, r.URL, r.RemoteAddr, r.Proto, r.Referer(), r.UserAgent())
		handler.ServeHTTP(w, r)
	})
}

func Webserver(config *cfg.Cfg) error {
	// If otto is enabled add routes and handle them.
	if config.Otto == true {
		otto_path, err := filepath.Abs(config.OttoPath)
		if err != nil {
			log.Fatalf("Can't read %s: %s\n", config.OttoPath, err)
		}
		programs, err := ottoengine.Load(otto_path)
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
		return http.ListenAndServeTLS(net.JoinHostPort(config.Hostname, strconv.Itoa(config.Port)), config.Cert, config.Key, request_log(http.DefaultServeMux))
	}
	log.Println("Starting http://" + net.JoinHostPort(config.Hostname, strconv.Itoa(config.Port)))
	// Now start up the server and log transactions
	return http.ListenAndServe(net.JoinHostPort(config.Hostname, strconv.Itoa(config.Port)), request_log(http.DefaultServeMux))
}

func defaultEnvBool(environment_var string, default_value bool) bool {
	tmp := strings.ToLower(os.Getenv(environment_var))
	if tmp == "true" {
		return true
	}
	if tmp == "false" {
		return false
	}
	return default_value
}

func defaultEnvString(environment_var string, default_value string) string {
	tmp := os.Getenv(environment_var)
	if tmp != "" {
		return tmp
	}
	return default_value
}

func defaultEnvInt(environment_var string, default_value int) int {
	tmp := os.Getenv(environment_var)
	if tmp != "" {
		i, err := strconv.Atoi(tmp)
		if err != nil {
			Usage(1, environment_var+" must be an integer.")
		}
		return i
	}
	return default_value
}

func init() {
	const (
		help_usage      = "This help document."
		keygen_usage    = "Interactive tool to generate TLS certificates and keys"
		init_usage      = "Creates a basic project structure in the current working directory"
		use_tls_usage   = "When true this turns on TLS (https) support."
		key_usage       = "Path to your SSL key pem file."
		cert_usage      = "path to your SSL cert pem file."
		docroot_usage   = "This is your document root for static files."
		host_usage      = "Set this hostname for webserver."
		port_usage      = "Set the port number to listen on."
		otto_usage      = "When true this option turns on ottoengine. Uses the path defined by WS_OTTO_PATH environment variable or one provided by -O option."
		otto_path_usage = "Turns on otto engine using the path for route JavaScript route handlers"
		version_usage   = "Display the version number of ws command."
	)

	flag.BoolVar(&help, "help", false, help_usage)
	flag.BoolVar(&help, "h", false, help_usage)
	flag.BoolVar(&do_keygen, "keygen", false, keygen_usage)
	flag.BoolVar(&do_init, "init", false, init_usage)
	flag.BoolVar(&version, "version", false, version_usage)
	flag.BoolVar(&version, "v", false, version_usage)

	// Settable via environment
	use_tls = defaultEnvBool("WS_TLS", false)
	key = defaultEnvString("WS_KEY", "")
	cert = defaultEnvString("WS_CERT", "")
	docroot = defaultEnvString("WS_DOCROOT", "")
	otto = defaultEnvBool("WS_OTTO", false)
	otto_path = defaultEnvString("WS_OTTO_PATH", "dynamic")
	host = defaultEnvString("WS_HOST", "localhost")
	port = defaultEnvInt("WS_PORT", 8000)

	flag.BoolVar(&use_tls, "tls", use_tls, use_tls_usage)
	flag.StringVar(&key, "key", key, key_usage)
	flag.StringVar(&cert, "cert", cert, cert_usage)
	flag.StringVar(&docroot, "docroot", docroot, docroot_usage)
	flag.StringVar(&docroot, "D", docroot, docroot_usage)
	flag.StringVar(&host, "host", host, host_usage)
	flag.StringVar(&host, "H", host, host_usage)
	flag.IntVar(&port, "port", port, port_usage)
	flag.IntVar(&port, "P", port, port_usage)
	flag.BoolVar(&otto, "otto", otto, otto_usage)
	flag.BoolVar(&otto, "o", otto, otto_usage)
	flag.StringVar(&otto_path, "otto-path", otto_path, otto_path_usage)
	flag.StringVar(&otto_path, "O", otto_path, otto_path_usage)
}

func main() {
	flag.Parse()
	if version == true {
		fmt.Println(REVISION)
		os.Exit(0)
	}
	if help == true {
		Usage(0, "")
	}

	config, err := cfg.Configure(docroot, host, port, use_tls, cert, key, otto, otto_path)
	if err != nil {
		Usage(1, fmt.Sprintf("%s", err))
	}

	if do_keygen == true {
		certFilename, keyFilename, err := keygen.Keygen("etc/ssl", "cert.pem", "key.pem")
		if err != nil {
			log.Fatalf("%s\n", err)
		}
		fmt.Printf("Wrote %s and %s\n", certFilename, keyFilename)
		os.Exit(0)
	}

	if do_init == true {
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
	err = Webserver(config)
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(0)
}
