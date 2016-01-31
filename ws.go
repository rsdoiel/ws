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
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/rsdoiel/otto"
)

const (
	// Version is used as a release number number for ws, wsinit, wsindexer
	Version = "0.0.0"
)

// Configuration provides the basic settings used by _ws_ and _wsint_ commands.
type Configuration struct {
	URL    *url.URL
	HTDocs string
	JSDocs string
	SSLKey string
	SSLPem string
}

// Getenv scans the environment variables and updates the values in
// the configuration.
func (config *Configuration) Getenv() error {
	u, err := url.Parse(os.Getenv("WS_URL"))
	config.URL = u
	config.HTDocs = os.Getenv("WS_HTDOCS")
	config.JSDocs = os.Getenv("WS_JSDOCS")
	config.SSLKey = os.Getenv("WS_SSL_KEY")
	config.SSLPem = os.Getenv("WS_SSL_PEM")
	return err
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
func NewJSEngine() *otto.Otto {
	vm := otto.New()
	vm.Set("Getenv", func(call otto.FunctionCall) otto.Value {
		if len(call.ArgumentList) != 1 {
			log.Println("Getenv() expect one environment variable name.")
			result, _ := otto.ToValue("")
			return result
		}
		env := call.Argument(0).String()
		result, err := otto.ToValue(os.Getenv(env))
		if err != nil {
			log.Printf("Getenv() error, %s", err)
		}
		return result
	})
	vm.Set("HttpGet", func(call otto.FunctionCall) otto.Value {
		var headers []map[string]string
		if len(call.ArgumentList) != 2 {
			log.Printf("HttpGet() missing parameters, got %d", len(call.ArgumentList))
		}

		uri := call.Argument(0).String()
		err := call.Argument(1).ToStruct(&headers)

		client := &http.Client{}
		req, err := http.NewRequest("GET", uri, nil)
		if err != nil {
			log.Fatalf("Can't create a GET request for %s, %s", uri, err)
		}
		for _, header := range headers {
			for k, v := range header {
				req.Header.Set(k, v)
			}
		}
		resp, err := client.Do(req)
		if err != nil {
			log.Fatalf("Can't connect to %s, %s", uri, err)
		}
		defer resp.Body.Close()
		content, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("Can't read response %s, %s", uri, err)
		}
		result, err := vm.ToValue(fmt.Sprintf("%s", content))
		if err != nil {
			log.Fatalf("HttpGet(%q) error, %s", uri, err)
		}
		return result
	})
	vm.Set("HttpPost", func(call otto.FunctionCall) otto.Value {
		var (
			headers []map[string]string
		)
		if len(call.ArgumentList) != 3 {
			log.Printf("HttpPost() missing parameters, got %d", len(call.ArgumentList))
		}

		uri := call.Argument(0).String()
		err := call.Argument(1).ToStruct(&headers)
		payload := call.Argument(2).String()
		if err != nil {
			log.Fatalf("Could not write headers to struct, %s", err)
		}
		buf := strings.NewReader(payload)

		client := &http.Client{}
		req, err := http.NewRequest("POST", uri, buf)
		if err != nil {
			log.Fatalf("Can't create a POST request %s, %s", uri, err)
		}
		for _, header := range headers {
			for k, v := range header {
				req.Header.Set(k, v)
			}
		}
		resp, err := client.Do(req)
		if err != nil {
			log.Fatalf("Can't connect to %s, %s", uri, err)
		}
		defer resp.Body.Close()
		content, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("Can't read response %s, %s", uri, err)
		}
		result, err := vm.ToValue(fmt.Sprintf("%s", content))
		if err != nil {
			log.Fatalf("HttpGet(%q) error, %s", uri, err)
		}
		return result
	})
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
