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
	"os"
	"strconv"

	// Local package
	"github.com/rsdoiel/ws"
)

// Flag options
var (
	help       bool
	initialize bool
	version    bool
	uri        string
	htdocs     string
	jsdocs     string
	sslkey     string
	sslcert    string
	cfg        *ws.Configuration
)

func usage() {
	fmt.Println(`
 USAGE: ws [OPTIONS]

 OVERVIEW

 ws is a utility for prototyping web services and sites. Start
 it with the "init" options will generate a default directory structure
 in your current path along with selfsigned certs if the url to listen
 for uses the https protocol (e.g. ws -url https://localhost:8443 init).

 OPTIONS
`)
	flag.VisitAll(func(f *flag.Flag) {
		fmt.Printf("    -%s  (defaults to %s) %s\n", f.Name, f.DefValue, f.Usage)
	})

	fmt.Println(`
 EXAMPLES

 Run a static web server using the content in the current directory
 (assumes the environment variables WS_HTDOCS and WS_JSDOCS are not defined).

   ws

 Setup a SSL base site saving the configuration in setup.conf.

   ws -url https://localhost:8443 -init
   . setup.bash
   ws

 Setup a standard project without SSL.

   ws -url http://localhost:8000 -init
   . setup.bash
   ws

 Turn on JavaScript server side support by providing a path for WS_JSDOCS
 and another location for WS_HTDOCS using local machine certs.

   ws -url https://localhost:8443 \
      -key /etc/ssl/sites/mysite.key \
      -cert /etc/ssl/sites/mysite.crt \
      -htdocs $HOME/Sites \
      -jsdocs ./jsdocs
`)
	os.Exit(0)
}

func init() {
	cfg = new(ws.Configuration)
	cfg.Getenv()
	uri = cfg.URL.String()
	if uri == "" {
		uri = "http://localhost:8000"
	}
	// If htdocs is NOT specified then turn off Server Side JavaScript support.
	htdocs = cfg.HTDocs
	jsdocs = cfg.JSDocs
	sslkey = cfg.SSLKey
	sslcert = cfg.SSLCert
	if htdocs == "" {
		htdocs = "."
		jsdocs = ""
	}

	flag.BoolVar(&help, "h", false, "Display this help message")
	flag.BoolVar(&help, "help", false, "Display this help message")
	flag.BoolVar(&version, "v", false, "Should version info")
	flag.BoolVar(&version, "version", false, "Should version info")
	flag.BoolVar(&initialize, "init", false, "Initialize a project")
	flag.StringVar(&htdocs, "htdocs", htdocs, "Set the htdocs path")
	flag.StringVar(&jsdocs, "jsdocs", jsdocs, "Set the jsdocs path, turns on server side JavaScript support")
	flag.StringVar(&uri, "url", uri, "The protocal and hostname listen for as a URL")
	flag.StringVar(&sslkey, "key", sslkey, "Set the path for the SSL Key")
	flag.StringVar(&sslcert, "cert", sslcert, "Set the path for the SSL Cert")
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

func makeJSHandler(route string, jsSource []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logRequest(r)
		vm := ws.NewJSEngine(w, r)
		_, err := vm.Eval(jsSource)
		if err != nil {
			log.Printf("JavaScript error %s, %s", route, err)
			http.Error(w, "Internal Server Error", 500)
			return
		}
		val, _ := vm.Get("Response")
		res := new(ws.JSResponse)
		err = val.ToStruct(&res)
		if err != nil {
			log.Printf("Can't unpack response %s, %s", route, err)
			http.Error(w, "Internal Server Error", 500)
			return
		}
		statusCode, _ := strconv.Atoi(fmt.Sprintf("%d", res.Code))
		w.Header().Set("Status-Code", fmt.Sprintf("%d %s", statusCode, http.StatusText(statusCode)))
		for _, header := range res.Headers {
			for k, v := range header {
				w.Header().Set(k, v)
			}
		}
		w.Write([]byte(res.Content))
	}
}

func main() {
	flag.Parse()

	// Process flags and update the environment as needed.
	if help == true {
		usage()
	}
	if version == true {
		fmt.Printf("ws version %s\n", ws.Version)
		os.Exit(0)
	}
	if uri != "" {
		os.Setenv("WS_URL", uri)
	}
	if htdocs != "" {
		os.Setenv("WS_HTDOCS", htdocs)
	}
	if jsdocs != "" {
		os.Setenv("WS_JSDOCS", jsdocs)
	}
	if sslkey != "" {
		os.Setenv("WS_SSL_KEY", sslkey)
	}
	if sslcert != "" {
		os.Setenv("WS_SSL_CERT", sslcert)
	}
	// Merge the environment changes
	cfg.Getenv()

	// Run through intialiation process if requested.
	if initialize == true {
		if cfg.HTDocs == "" || cfg.HTDocs == "." {
			cfg.HTDocs = "htdocs"
		}
		if cfg.JSDocs == "" {
			cfg.JSDocs = "jsdocs"
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
		ioutil.WriteFile("setup.bash", []byte(setup), 0660)
		log.Println("Wrote setup to setup.bash")
		os.Exit(0)
	}

	// Do a final sanity check before starting up web server
	err := cfg.Validate()
	if err != nil {
		log.Fatalf("Invalid configuration, %s", err)
	}

	if cfg.JSDocs != "" {
		log.Printf("JSDocs %s", cfg.JSDocs)
		jsSourceFiles, err := ws.ReadJSFiles(cfg.JSDocs)
		if err != nil {
			log.Fatalf("Could not read files in %s, %s", cfg.JSDocs, err)
		}
		for fname, jsSource := range jsSourceFiles {
			route, err := ws.JSPathToRoute(fname, cfg)
			if err != nil {
				log.Fatalf("%s", err)
			}
			log.Printf("Adding route %s for %s", route, fname)
			jsHandler := makeJSHandler(route, jsSource)
			http.HandleFunc(route, jsHandler)
		}
	}

	log.Printf("HTDocs %s", cfg.HTDocs)
	log.Printf("Listening for %s", cfg.URL.Host)
	http.Handle("/", http.FileServer(http.Dir(cfg.HTDocs)))
	if cfg.URL.Scheme == "https" {
		http.ListenAndServeTLS(cfg.URL.Host, cfg.SSLCert, cfg.SSLKey, logger(http.DefaultServeMux))
	} else {
		http.ListenAndServe(cfg.URL.Host, logger(http.DefaultServeMux))
	}
}
