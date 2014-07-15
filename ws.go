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
    "./fsengine"
	"./ottoengine"
	"./wslog"
    "./app"
    "./keygen"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
)

var REVISION = "v0.0.0-alpha"

// command line parameters that override environment variables
var (
	cli_use_tls   *bool
	cli_docroot   *string
	cli_host      *string
	cli_port      *string
	cli_cert      *string
	cli_key       *string
	cli_otto      *bool
	cli_otto_path *string
	cli_version   *bool
	cli_keygen   = flag.Bool("keygen", false, "Interactive tool to generate TLS certificates and keys")
    cli_init     = flag.Bool("init", false, "Creates a basic project structure in the current working directory")
)

var Usage = func() {
	flag.PrintDefaults()
}

func request_log(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wslog.LogRequest(r.Method, r.URL, r.RemoteAddr, r.Proto, r.Referer(), r.UserAgent())
		handler.ServeHTTP(w, r)
	})
}

func Webserver(profile *app.Profile) error {
	// If otto is enabled add routes and handle them.
	if profile.Otto == true {
		otto_path, err := filepath.Abs(profile.Otto_Path)
		if err != nil {
			log.Fatalf("Can't read %s: %s\n", profile.Otto_Path, err)
		}
		programs, err := ottoengine.Load(otto_path)
		if err != nil {
			log.Fatalf("Load error: %s\n", err)
		}
		ottoengine.AddRoutes(programs)
	}

	// Restricted FileService excluding dot files and directories
	http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
        // hande off this request/response pair to the fsengine
        fsengine.Engine(profile, w, r)
    })

	// Now start up the server and log transactions
	if profile.Use_TLS == true {
		if profile.Cert == "" || profile.Key == "" {
			log.Fatalf("TLS set true but missing key or certificate")
		}
		log.Println("Starting https://" + net.JoinHostPort(profile.Hostname, profile.Port))
		return http.ListenAndServeTLS(net.JoinHostPort(profile.Hostname, profile.Port), profile.Cert, profile.Key, request_log(http.DefaultServeMux))
	}
	log.Println("Starting http://" + net.JoinHostPort(profile.Hostname, profile.Port))
	// Now start up the server and log transactions
	return http.ListenAndServe(net.JoinHostPort(profile.Hostname, profile.Port), request_log(http.DefaultServeMux))
}

func init() {
	// command line parameters that override environment variables, many have short forms too.
	shortform := " (short form)"

	// No short form
	cli_use_tls = flag.Bool("tls", false, "When true this turns on TLS (https) support.")
	cli_key = flag.String("key", "", "path to your SSL key pem file.")
	cli_cert = flag.String("cert", "", "path to your SSL cert pem file.")

	// These have short forms too
	msg := "This is your document root for static files."
	cli_docroot = flag.String("docroot", "", msg)
	flag.StringVar(cli_docroot, "D", "", msg+shortform)

	msg = "Set this hostname for webserver."
	cli_host = flag.String("host", "", msg)
	flag.StringVar(cli_host, "H", "", msg+shortform)

	msg = "Set the port number to listen on."
	cli_port = flag.String("port", "", msg)
	flag.StringVar(cli_port, "P", "", msg+shortform)

	msg = "When true this option turns on ottoengine. Uses the path defined by WS_OTTO_PATH environment variable or one provided by -O option."
	cli_otto = flag.Bool("otto", false, msg)
	flag.BoolVar(cli_otto, "o", false, msg+shortform)

	msg = "Turns on otto engine using the path for route JavaScript route handlers"
	cli_otto_path = flag.String("otto-path", "", msg)
	flag.StringVar(cli_otto_path, "O", "", msg+shortform)

	msg = "Display the version number of ws command."
	cli_version = flag.Bool("version", false, msg)
	flag.BoolVar(cli_version, "v", false, msg+shortform)

}

func main() {
	flag.Parse()

	if *cli_version == true {
		fmt.Println(REVISION)
		os.Exit(0)
	}

	profile, _ := app.LoadProfile(*cli_docroot, *cli_host, *cli_port, *cli_use_tls, *cli_cert, *cli_key, *cli_otto, *cli_otto_path)
	if *cli_keygen == true {
		err := keygen.Keygen(profile)
		if err != nil {
			log.Fatalf("%s\n", err)
		}
		os.Exit(0)
	}

    if *cli_init == true {
        err := app.Init()
        if err != nil {
            log.Fatalf("%s\n", err)
        }
        os.Exit(0)
    }

	log.Printf("\n\n"+
		"          TLS: %t\n"+
		"         Cert: %s\n"+
		"          Key: %s\n"+
		"      Docroot: %s\n"+
		"         Host: %s\n"+
		"         Port: %s\n"+
		"       Run as: %s\n\n"+
		" Otto enabled: %t\n"+
		"         Path: %s\n"+
		"\n\n",
		profile.Use_TLS,
		profile.Cert,
		profile.Key,
		profile.Docroot,
		profile.Hostname,
		profile.Port,
		profile.Username,
		profile.Otto,
		profile.Otto_Path)
	err := Webserver(profile)
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(0)
}
