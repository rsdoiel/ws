//
// ws.go - A simple web server for static files and limit server side JavaScript
// @author R. S. Doiel, <rsdoiel@gmail.com>
//
// Copyright (c) 2014 - 2016, R. S. Doiel
// All rights not granted herein are expressly reserved by R. S. Doiel.
//
// Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
//
// 2. Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.
//
// 3. Neither the name of the copyright holder nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"

	// Local package
	"github.com/rsdoiel/ws"
)

// Flag options
var (
	showHelp    bool
	showVersion bool
	showLicense bool
	initialize  bool
	uri         string
	docRoot     string
	sslKey      string
	sslCert     string
	cfg         *ws.Configuration
)

func usage(fp *os.File, appName string) {
	fmt.Println(`
 USAGE: %s [OPTIONS]

 OVERVIEW

 ws is a utility for prototyping web services and sites. Start
 it with the "init" options will generate a default directory structure
 in your current path along with selfsigned certs if the url to listen
 for uses the https protocol (e.g. ws -url https://localhost:8443 init).

 OPTIONS
`, appName)
	flag.VisitAll(func(f *flag.Flag) {
		if len(f.Name) > 1 {
			fmt.Fprintf(fp, "    -%s, -%s\t%s\n", f.Name[0:1], f.Name, f.Usage)
		}
	})

	fmt.Fprintln(fp, `
 EXAMPLES

 Run web server using the content in the current directory
 (assumes the environment variables WS_DOCROOT are not defined).

   ws

 Run web server using a specified directory

   ws /www/htdocs

 Setup a SSL base site saving the configuration in setup.conf.

   ws -url https://localhost:8443 -docs /www/htdocs -init
   . setup.bash
   ws

 Setup a standard project without SSL.

   ws -url http://localhost:8000 -docs $HOME/Sites/example.me -init
   . setup.bash
   ws

`)
}

func license(fp *os.File, appName string) {
	fmt.Fprintf(fp, `
%s %s

Copyright (c) 2014 - 2016, R. S. Doiel
All rights reserved.

Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.

3. Neither the name of the copyright holder nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

`, appName, ws.Version)
}

func init() {
	flag.BoolVar(&showHelp, "h", false, "Display this help message")
	flag.BoolVar(&showHelp, "help", false, "Display this help message")
	flag.BoolVar(&showVersion, "v", false, "Should version info")
	flag.BoolVar(&showVersion, "version", false, "Should version info")
	flag.BoolVar(&showLicense, "l", false, "Should license info")
	flag.BoolVar(&showLicense, "license", false, "Should license info")
	flag.BoolVar(&initialize, "i", false, "Initialize a project")
	flag.BoolVar(&initialize, "init", false, "Initialize a project")
	flag.StringVar(&docRoot, "d", "", "Set the htdocs path")
	flag.StringVar(&docRoot, "docs", "", "Set the htdocs path")
	flag.StringVar(&uri, "u", "", "The protocal and hostname listen for as a URL")
	flag.StringVar(&uri, "url", "", "The protocal and hostname listen for as a URL")
	flag.StringVar(&sslKey, "k", "", "Set the path for the SSL Key")
	flag.StringVar(&sslKey, "key", "", "Set the path for the SSL Key")
	flag.StringVar(&sslCert, "c", "", "Set the path for the SSL Cert")
	flag.StringVar(&sslCert, "cert", "", "Set the path for the SSL Cert")
}

func logRequest(r *http.Request) {
	log.Printf("Request: %s Path: %s RemoteAddr: %s UserAgent: %s\n", r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent())
}

func logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logRequest(r)
		next.ServeHTTP(w, r)
	})
}

func main() {
	appName := path.Base(os.Args[0])
	flag.Parse()

	// Process flags and update the environment as needed.
	if showHelp == true {
		usage(os.Stdout, appName)
		os.Exit(0)
	}
	if showLicense == true {
		license(os.Stdout, appName)
		os.Exit(0)
	}
	if showVersion == true {
		fmt.Printf("%s version %s\n", appName, ws.Version)
		os.Exit(0)
	}

	cfg = new(ws.Configuration)
	cfg.SetDefaults()
	cfg.Getenv()

	// Merge command line options
	if uri != "" {
		u, err := url.Parse(uri)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Can't parse url %s, %s\n", uri, err)
			os.Exit(1)
		}
		cfg.URL = u
	}
	if docRoot != "" {
		cfg.DocRoot = docRoot
	}
	if sslKey != "" {
		cfg.SSLKey = sslKey
	}
	if sslCert != "" {
		cfg.SSLCert = sslCert
	}

	// setup from command line
	args := flag.Args()
	if len(args) > 0 {
		cfg.DocRoot = args[0]
	}

	// Run through initialization process if requested.
	if initialize == true {
		if cfg.DocRoot == "" || cfg.DocRoot == "." {
			cfg.DocRoot = "htdocs"
		}
		if cfg.URL.Scheme == "https" {
			if cfg.SSLKey == "" {
				cfg.SSLKey = "etc/ssl/site.key"
			}
			if cfg.SSLCert == "" {
				cfg.SSLCert = "etc/ssl/site.crt"
			}
		}
		setup, err := cfg.InitializeProject()
		if err != nil {
			log.Fatalf("%s", err)
		}
		// Do a sanity check before generating project
		err = cfg.Validate()
		if err != nil {
			log.Fatalf("Proposed configuration not valid, %s", err)
		}
		// FIXME: this should go in etc/...
		ioutil.WriteFile("setup.bash", []byte(setup), 0660)
		log.Println("Wrote setup to setup.bash")
		// FIXME: generate a suitable systemd startup script example
		// FIXME: generate a default templates
		// FIXME: generate a css/site.css file in doc root
		os.Exit(0)
	}

	// Do a final sanity check before starting up web server
	err := cfg.Validate()
	if err != nil {
		log.Fatalf("Invalid configuration, %s", err)
	}

	log.Printf("DocRoot %s", cfg.DocRoot)
	log.Printf("Listening for %s", cfg.URL.Host)
	http.Handle("/", http.FileServer(http.Dir(cfg.DocRoot)))
	if cfg.URL.Scheme == "https" {
		http.ListenAndServeTLS(cfg.URL.Host, cfg.SSLCert, cfg.SSLKey, logger(http.DefaultServeMux))
	} else {
		http.ListenAndServe(cfg.URL.Host, logger(http.DefaultServeMux))
	}
}
