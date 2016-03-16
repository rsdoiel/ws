//
// Package ws provides the core library used by cmds/ws/ws.go, cmds/wsinit/wsinit.go and
// cmds/wsindexer/wsindexer.go
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
package ws

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

    // Caltech Library Packages
	"github.com/caltechlibrary/otto"
)

const (
	// Version is used as a release number number for ws, wsinit, wsindexer
	Version = "0.0.7"
)

// Configuration provides the basic settings used by _ws_ and _wsint_ commands.
type Configuration struct {
	URL     *url.URL
	HTDocs  string
	JSDocs  string
	SSLKey  string
	SSLCert string
}

// JSRequest provides a struct for server side JS request processing
type JSRequest struct {
	URL    *url.URL               `json:"url"`
	Method string                 `json:"method"`
	Body   []byte                 `json:"body"`
	Form   map[string]interface{} `json:"form"`
}

// JSResponse provides a structure for server side JS responses
type JSResponse struct {
	Code    int64               `json:"code,omitempty"`
	Headers []map[string]string `json:"headers,omitempty"`
	Content string              `json:"content,omitempty"`
}

var configFileTemplate = `#!/bin/bash
# generated %d-%02d-%02d by ws version %s
export WS_URL=%q
export WS_HTDOCS=%q
export WS_JSDOCS=%q
export WS_SSL_KEY=%q
export WS_SSL_CERT=%q
`

// Getenv scans the environment variables and updates the fields in
// the configuration. If there is a problem parsing WS_URL then an
// error is returned.
func (config *Configuration) Getenv() error {
	u, err := url.Parse(os.Getenv("WS_URL"))
	config.URL = u
	config.HTDocs = os.Getenv("WS_HTDOCS")
	config.JSDocs = os.Getenv("WS_JSDOCS")
	config.SSLKey = os.Getenv("WS_SSL_KEY")
	config.SSLCert = os.Getenv("WS_SSL_CERT")
	return err
}

// String returns a multiline block of text suitable to save in a text file.
func (config *Configuration) String() string {
	var u string
	now := time.Now()
	yr, mn, dy := now.Date()
	if config.URL != nil {
		u = config.URL.String()
	}
	return fmt.Sprintf(configFileTemplate, yr, mn, dy, Version, u, config.HTDocs, config.JSDocs, config.SSLKey, config.SSLCert)
}

// Validate performs a sanity check of the values in the configuration
// Returns nil if OK otherwise returns error
func (config *Configuration) Validate() error {
	scheme := config.URL.Scheme
	htdocs := config.HTDocs
	jsdocs := config.JSDocs
	sslkey := config.SSLKey
	sslcert := config.SSLCert
	if scheme == "https" {
		if sslkey == "" || sslcert == "" {
			return fmt.Errorf("Cannot use https without specifying SSL Cert and Key")
		}
		if _, err := os.Stat(sslkey); os.IsNotExist(err) {
			return fmt.Errorf("Cannot find %s, %s", sslkey, err)
		}
		if _, err := os.Stat(sslcert); os.IsNotExist(err) {
			return fmt.Errorf("Cannot find %s, %s", sslcert, err)
		}
	}
	if _, err := os.Stat(htdocs); os.IsNotExist(err) {
		return fmt.Errorf("Can't find %s, %s", htdocs, err)
	}
	if jsdocs != "" {
		if _, err := os.Stat(jsdocs); os.IsNotExist(err) {
			return fmt.Errorf("Can't find %s, %s", jsdocs, err)
		}
		htdir, err := filepath.Abs(htdocs)
		if err != nil {
			return fmt.Errorf("Can't determine htdocs path %s, %s", htdocs, err)
		}
		jsdir, err := filepath.Abs(jsdocs)
		if err != nil {
			return fmt.Errorf("Can't determine jsdocs path %s, %s", jsdocs, err)
		}
		if strings.HasPrefix(jsdir, htdir) == true {
			return fmt.Errorf("Unsafe jsdocs location %s, is child of %s", jsdocs, htdocs)
		}
	}
	return nil
}

// GenerateKeyAndCert will generate a new SSL Key/Cert pair if none exist.
// It will return a error if they already exist.
func (config *Configuration) GenerateKeyAndCert() error {
	// Double check to see if we need a self-signed cert...
	if config.URL.Scheme != "https" {
		log.Printf("Skipping key/cert creation, not needed for %s", config.URL.String())
		return nil
	}
	// Check to see if config.SSLKey and config.SSLCert already exist, if not create them.
	if config.SSLKey == "" || config.SSLCert == "" {
		log.Println("Skipping key/cert creation, not defined by config")
		return nil
	}
	if _, err := os.Stat(config.SSLCert); os.IsExist(err) == true {
		log.Printf("%s already exists, skipping key creation", config.SSLCert)
		return nil
	}
	if _, err := os.Stat(config.SSLKey); os.IsExist(err) == true {
		log.Printf("%s already exists, skipping key creation", config.SSLKey)
		return nil
	}

	hostname := config.URL.Host
	// Trim our hostname before the port number if needed
	if pos := strings.Index(hostname, ":"); pos > 0 {
		hostname = hostname[0:pos]
	}
	if hostname == "" {
		hostname = "localhost"
	}

	rsaBits := 2048
	log.Printf("Generating %d bit key", rsaBits)
	priv, err := rsa.GenerateKey(rand.Reader, rsaBits)
	if err != nil {
		return err
	}
	notBefore := time.Now()
	yr, _, _ := notBefore.Date()
	yr++
	notAfter := time.Date(yr, 12, 31, 2, 59, 59, 0, time.UTC)

	log.Println("Setting up cerificates")
	template := x509.Certificate{
		SerialNumber: new(big.Int).SetInt64(0),
		Subject: pkix.Name{
			Organization: []string{"Acme Co"},
		},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	log.Println("Checking IP address")
	if ip := net.ParseIP(hostname); ip != nil {
		template.IPAddresses = append(template.IPAddresses, ip)
	} else {
		template.DNSNames = append(template.DNSNames, hostname)
	}

	template.IsCA = false
	log.Println("Generating x509 certs from template")
	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		return err
	}

	log.Printf("Creating %s\n", config.SSLCert)
	// make the directory if needed.
	dname, _ := filepath.Split(config.SSLCert)
	if _, err := os.Stat(dname); os.IsNotExist(err) {
		err := os.MkdirAll(dname, 0770)
		if err != nil {
			return fmt.Errorf("Can't create directory %s, %s", dname, err)
		}
	}

	certOut, err := os.Create(config.SSLCert)
	if err != nil {
		return err
	}
	log.Printf("Encode %s as pem", config.SSLCert)
	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	certOut.Close()
	log.Printf("Wrote %s\n", config.SSLCert)

	// make the directory if needed.
	dname, _ = filepath.Split(config.SSLKey)
	if _, err := os.Stat(dname); os.IsNotExist(err) {
		err := os.MkdirAll(dname, 0770)
		if err != nil {
			return fmt.Errorf("Can't create directory %s, %s", dname, err)
		}
	}
	log.Printf("Creating %s\n", config.SSLKey)
	keyOut, err := os.OpenFile(config.SSLKey, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	log.Printf("Encode %s as pem", config.SSLKey)
	pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})
	keyOut.Close()
	log.Printf("Wrote %s\n", config.SSLKey)
	// We got this far so no errors
	return nil
}

// InitializeProject scans the current working path and identifies what
// directories need to be created, creates them, adds SSL keys/certs if
// needed and returns a suggested configuration file.
func (config *Configuration) InitializeProject() (string, error) {
	// Only append SSL dir and files if config.URL.Scheme is https
	directories := []string{
		config.HTDocs,
		config.JSDocs,
	}
	if config.URL.Scheme == "https" {
		sslKeyDir, _ := filepath.Split(config.SSLKey)
		sslCertDir, _ := filepath.Split(config.SSLCert)
		directories = append(directories, sslKeyDir)
		directories = append(directories, sslCertDir)
	}

	// Check if directory exists, if not created it
	for _, directory := range directories {
		if directory != "" {
			_, err := os.Stat(directory)
			if os.IsNotExist(err) == true {
				if os.MkdirAll(directory, 0775) != nil {
					return "", fmt.Errorf("Can't create %s, %s", directory, err)
				}
			}
		}
	}
	if config.URL.Scheme == "https" {
		if err := config.GenerateKeyAndCert(); err != nil {
			return config.String(), err
		}
	}
	// return a suggested config file or error
	return config.String(), nil
}

// ReadJSFiles walks a directory tree and then return the results
// as a map of paths and JS source code in []byte.
func ReadJSFiles(jsDocs string) (map[string][]byte, error) {
	jsSources := make(map[string][]byte)
	err := filepath.Walk(jsDocs, func(path string, _ os.FileInfo, _ error) error {
		if strings.HasSuffix(path, ".js") {
			src, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			jsSources[path] = src
		}
		return nil
	})
	return jsSources, err
}

// NewJSEngine creates a new JavaScript version machine from otto.New() but
// adds additional functionality such as WS.Getenv(), WW.httpGet(), WS.httpPost()
func NewJSEngine(w http.ResponseWriter, r *http.Request) *otto.Otto {
	vm := otto.New()
	wsObject, _ := vm.Object(`WS = {}`)
	wsObject.Set("getEnv", func(call otto.FunctionCall) otto.Value {
		if len(call.ArgumentList) != 1 {
			log.Println("WS.getEnv() expect one environment variable name.")
			result, _ := otto.ToValue("")
			return result
		}
		env := call.Argument(0).String()
		result, err := otto.ToValue(os.Getenv(env))
		if err != nil {
			log.Printf("WS.getEnv() error, %s", err)
		}
		return result
	})
	wsObject.Set("httpGet", func(call otto.FunctionCall) otto.Value {
		var headers []map[string]string
		if len(call.ArgumentList) != 2 {
			log.Printf("WS.httpGet() missing parameters, got %d", len(call.ArgumentList))
		}

		uri := call.Argument(0).String()
		err := call.Argument(1).ToStruct(&headers)

		client := &http.Client{}
		req, err := http.NewRequest("GET", uri, nil)
		if err != nil {
			log.Printf("Can't create a GET request for %s, %s", uri, err)
		}
		for _, header := range headers {
			for k, v := range header {
				req.Header.Set(k, v)
			}
		}
		res, err := client.Do(req)
		if err != nil {
			log.Printf("Can't connect to %s, %s", uri, err)
		}
		defer res.Body.Close()
		content, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Printf("Can't read response %s, %s", uri, err)
		}
		result, err := vm.ToValue(fmt.Sprintf("%s", content))
		if err != nil {
			log.Printf("WS.httpGet(%q) error, %s", uri, err)
		}
		return result
	})
	wsObject.Set("httpPost", func(call otto.FunctionCall) otto.Value {
		var (
			headers []map[string]string
		)
		if len(call.ArgumentList) != 3 {
			log.Printf("WS.httpPost() missing parameters, got %d", len(call.ArgumentList))
		}

		uri := call.Argument(0).String()
		err := call.Argument(1).ToStruct(&headers)
		payload := call.Argument(2).String()
		if err != nil {
			log.Printf("Could not write headers to struct, %s", err)
		}
		buf := strings.NewReader(payload)

		client := &http.Client{}
		req, err := http.NewRequest("POST", uri, buf)
		if err != nil {
			log.Printf("Can't create a POST request %s, %s", uri, err)
		}
		for _, header := range headers {
			for k, v := range header {
				req.Header.Set(k, v)
			}
		}
		res, err := client.Do(req)
		if err != nil {
			log.Printf("Can't connect to %s, %s", uri, err)
		}
		defer res.Body.Close()
		content, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Printf("Can't read response %s, %s", uri, err)
		}
		result, err := vm.ToValue(fmt.Sprintf("%s", content))
		if err != nil {
			log.Printf("WS.httpGet(%q) error, %s", uri, err)
		}
		return result
	})

	// Add request and reponse of available
	req := new(JSRequest)
	res := new(JSResponse)
	route := ""
	if r != nil {
		route = r.URL.Path
		req.URL = r.URL
		req.Method = r.Method
		// FIXME: need to handle GET, POST, PUT, DELETE here...
		src, err := json.Marshal(req)
		if err != nil {
			log.Printf("Can't marshal request %s, %+v, %s", route, req, err)
			http.Error(w, "Internal Server Error", 500)
			return vm
		}
		vm.Eval(fmt.Sprintf(`var Request = %s;`, src))
	}

	if w != nil {
		src, err := json.Marshal(res)
		if err != nil {
			log.Printf("Can't marshal response %s, %+v, %s", route, res, err)
			http.Error(w, "Internal Server Error", 500)
			return vm
		}
		vm.Eval(fmt.Sprintf(`var Response = %s;`, src))
	}

	return vm
}

// JSPathToRoute converts a JSDocs path to a JavaScript file into a web server
// route.
func JSPathToRoute(p string, cfg *Configuration) (string, error) {
	// Check to see if the route is relative to JSDocs
	rel, err := filepath.Rel(cfg.JSDocs, p)
	if err != nil {
		return "", err
	}
	ext := filepath.Ext(rel)
	return "/" + strings.TrimSuffix(rel, ext), nil
}
