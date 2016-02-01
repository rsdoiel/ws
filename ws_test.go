//
// ws_test.go test routines for ws.go
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
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestConfigureFromEnv(t *testing.T) {
	cfg := new(Configuration)
	if cfg.URL != nil {
		t.Errorf("cfg.URL should be nil, %v", cfg.URL)
	}
	if cfg.HTDocs != "" {
		t.Errorf("cfg.HTDocs should be empty, %s", cfg.HTDocs)
	}
	if cfg.JSDocs != "" {
		t.Errorf("cfg.JSDocs should be empty, %s", cfg.JSDocs)
	}
	if cfg.SSLKey != "" {
		t.Errorf("cfg.SSLKey should be empty, %s", cfg.SSLKey)
	}
	if cfg.SSLCert != "" {
		t.Errorf("cfg.SSLCert should be empty, %s", cfg.SSLCert)
	}

	os.Setenv("WS_URL", "https://example.org:8001")
	os.Setenv("WS_HTDOCS", "htdocs")
	os.Setenv("WS_JSDOCS", "jsdocs")
	os.Setenv("WS_SSL_KEY", "etc/ssl/site.key")
	os.Setenv("WS_SSL_CERT", "etc/ssl/site.crt")
	err := cfg.Getenv()
	if err != nil {
		t.Errorf("cfg.Getenv() error, %s", err)
	}
	if cfg.URL.Host != "example.org:8001" {
		t.Errorf("cfg.URL.Host != example.org:8001, %s", cfg.URL.Host)
	}
	if cfg.HTDocs != "htdocs" {
		t.Errorf("cfg.HTDocs != htdocs, %s", cfg.HTDocs)
	}
	if cfg.JSDocs != "jsdocs" {
		t.Errorf("cfg.JSDocs != jsdocs, %s", cfg.JSDocs)
	}
	if cfg.SSLKey != "etc/ssl/site.key" {
		t.Errorf("cfg.SSLKey != etc/ssl/site.key, %s", cfg.SSLKey)
	}
	if cfg.SSLCert != "etc/ssl/site.crt" {
		t.Errorf("cfg.SSLCert != etc/ssl/site.crt, %s", cfg.SSLCert)
	}
}

func TestConfigureString(t *testing.T) {
	cfg := new(Configuration)

	now := time.Now()
	yr, mn, dy := now.Date()
	u := "https://example.org"
	htdocs := "/www/htdocs"
	jsdocs := "/www/jsdocs"
	sslkey := "/etc/ssl/site.key"
	sslcert := "/etc/ssl/site.crt"

	// Set some example configuration
	os.Setenv("WS_URL", u)
	os.Setenv("WS_HTDOCS", htdocs)
	os.Setenv("WS_JSDOCS", jsdocs)
	os.Setenv("WS_SSL_KEY", sslkey)
	os.Setenv("WS_SSL_CERT", sslcert)
	cfg.Getenv()

	expected := fmt.Sprintf(`#!/bin/bash
# generated %d-%02d-%02d by ws version %s
export WS_URL=%q
export WS_HTDOCS=%q
export WS_JSDOCS=%q
export WS_SSL_KEY=%q
export WS_SSL_CERT=%q
`, yr, mn, dy, Version, u, htdocs, jsdocs, sslkey, sslcert)

	s := cfg.String()
	if strings.Compare(s, expected) != 0 {
		t.Errorf("found\n%s\nexpected\n%s\n", s, expected)
	}
}

func TestConfigureInitializeProject(t *testing.T) {
	//
	// Test https setup
	//
	cfg := new(Configuration)
	u := "https://testout.localhost:8001"
	htdocs := "testout/https/htdocs"
	jsdocs := "testout/https/jsdocs"
	sslkey := "testout/https/etc/ssl/site.key"
	sslcert := "testout/https/etc/ssl/site.crt"
	// Set some example configuration
	os.Setenv("WS_URL", u)
	os.Setenv("WS_HTDOCS", htdocs)
	os.Setenv("WS_JSDOCS", jsdocs)
	os.Setenv("WS_SSL_KEY", sslkey)
	os.Setenv("WS_SSL_CERT", sslcert)
	cfg.Getenv()
	expectedSetup := cfg.String()
	// Now see if we can create our test environment for project
	setup, err := cfg.InitializeProject()
	if err != nil {
		t.Errorf("Can't initialize project %s\n%s", cfg.String(), err)
	}
	// New check to see if the directory paths really exist or not...
	sslKeyDir, _ := filepath.Split(cfg.SSLKey)
	sslCertDir, _ := filepath.Split(cfg.SSLCert)
	directories := []string{
		cfg.HTDocs,
		cfg.JSDocs,
		sslKeyDir,
		sslCertDir,
	}
	for _, directory := range directories {
		if stat, err := os.Stat(directory); os.IsNotExist(err) || stat.IsDir() == false {
			t.Errorf("Missing %s, %s", directory, err)
		}
	}
	filenames := []string{
		cfg.SSLKey,
		cfg.SSLCert,
	}
	for _, filename := range filenames {
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			t.Errorf("Missing %s, %s", filename, err)
		}
	}

	if setup != expectedSetup {
		t.Errorf("     Got %s\nExpected %s\n", setup, expectedSetup)
	}

	//
	// Now test http setup
	//
	cfg = new(Configuration)
	u = "http://testout.localhost:8001"
	htdocs = "testout/http/htdocs"
	jsdocs = "testout/http/jsdocs"
	sslkey = "testout/http/etc/ssl/site.key"
	sslcert = "testout/http/etc/ssl/site.crt"
	// Set some example configuration
	os.Setenv("WS_URL", u)
	os.Setenv("WS_HTDOCS", htdocs)
	os.Setenv("WS_JSDOCS", jsdocs)
	os.Setenv("WS_SSL_KEY", sslkey)
	os.Setenv("WS_SSL_CERT", sslcert)
	cfg.Getenv()
	expectedSetup = cfg.String()
	// Now see if we can create our test environment for project
	setup, err = cfg.InitializeProject()
	if err != nil {
		t.Errorf("Can't initialize project %s\n%s", cfg.String(), err)
	}
	// New check to see if the directory paths really exist or not...
	directories = []string{
		cfg.HTDocs,
		cfg.JSDocs,
	}
	for _, directory := range directories {
		if stat, err := os.Stat(directory); os.IsNotExist(err) || stat.IsDir() == false {
			t.Errorf("Missing %s, %s", directory, err)
		}
	}
	filenames = []string{
		cfg.SSLKey,
		cfg.SSLCert,
	}
	for _, filename := range filenames {
		if _, err := os.Stat(filename); os.IsExist(err) {
			t.Errorf("We should not find %s, %s", filename, err)
		}
	}

	if setup != expectedSetup {
		t.Errorf("     Got %s\nExpected %s\n", setup, expectedSetup)
	}
}

func TestConfigureValidate(t *testing.T) {
	cfg := new(Configuration)
	u := "https://example.org"
	htdocs := "testout/www-01/htdocs"
	jsdocs := "testout/www-01/jsdocs"
	sslkey := "testout/www-01/ssl/site.key"
	sslcert := "testout/www-01/ssl/site.crt"

	// Set some example configuration
	os.Setenv("WS_URL", u)
	os.Setenv("WS_HTDOCS", htdocs)
	os.Setenv("WS_JSDOCS", jsdocs)
	os.Setenv("WS_SSL_KEY", sslkey)
	os.Setenv("WS_SSL_CERT", sslcert)
	cfg.Getenv()
	cfg.InitializeProject()

	if err := cfg.Validate(); err != nil {
		t.Errorf("Should have been valid %+v, %s", cfg, err)
	}

	u = "http://example.org"
	htdocs = "testout/www-02/htdocs"
	jsdocs = "testout/www-02/jsdocs"
	sslkey = ""
	sslcert = ""
	os.Setenv("WS_URL", u)
	os.Setenv("WS_HTDOCS", htdocs)
	os.Setenv("WS_JSDOCS", jsdocs)
	os.Setenv("WS_SSL_KEY", sslkey)
	os.Setenv("WS_SSL_CERT", sslcert)
	cfg.Getenv()
	cfg.InitializeProject()

	if err := cfg.Validate(); err != nil {
		t.Errorf("Should have been valid %+v, %s", cfg, err)
	}

	os.Setenv("WS_URL", "https://example.org")
	cfg.Getenv()
	if err := cfg.Validate(); err == nil {
		t.Errorf("https without ssl files, should have been invalid %+v, %s", cfg, err)
	}

	u = "http://example.org"
	htdocs = "testout/www-02/htdocs"
	jsdocs = "testout/www-02/htdocs/js"
	sslkey = ""
	sslcert = ""
	os.Setenv("WS_URL", u)
	os.Setenv("WS_HTDOCS", htdocs)
	os.Setenv("WS_JSDOCS", jsdocs)
	os.Setenv("WS_SSL_KEY", sslkey)
	os.Setenv("WS_SSL_CERT", sslcert)
	cfg.Getenv()
	if err := cfg.Validate(); err == nil {
		t.Errorf("jsdocs a child of htdocs, should have been invalid %+v, %s", cfg, err)
	}
}

func TestReadJSFiles(t *testing.T) {
	jsSources, err := ReadJSFiles("jsdocs")
	if err != nil {
		t.Errorf("reading jsdocs %s", err)
	}
	if len(jsSources) < 1 {
		t.Error("should find more than one JavaScript file in jsdocs")
	}
	jsFilename := "jsdocs/helloworld.js"
	jsSrc, ok := jsSources[jsFilename]
	if ok == false {
		t.Errorf("Should find %s in jsSources", jsFilename)
	}
	if jsSrc == nil {
		t.Errorf("Should find source code for %s", jsFilename)
	}
}

func TestJSEngine(t *testing.T) {
	vm := NewJSEngine(nil, nil)
	if vm == nil {
		t.Errorf("should have created a new JavaScript VM")
	}

	val, err := vm.Eval(`
    Getenv != undefined
  `)
	b, err := val.ToBoolean()
	if err != nil {
		t.Errorf("Error from Getenv != undefined, %s", err)
	}
	if b == false {
		t.Errorf("Expected Getenv != undefined to return true, %b", b)
	}
	val, err = vm.Eval(`
    HttpGet != undefined
  `)
	b, err = val.ToBoolean()
	if err != nil {
		t.Errorf("Error from HttpGet != undefined, %s", err)
	}
	if b == false {
		t.Errorf("Expected HttpGet != undefined to return true, %b", b)
	}
	val, err = vm.Eval(`
    HttpPost != undefined
  `)
	b, err = val.ToBoolean()
	if err != nil {
		t.Errorf("Error from HttpPost != undefined, %s", err)
	}
	if b == false {
		t.Errorf("Expected HttpPost != undefined to return true, %b", b)
	}

	err = os.Setenv("WS_TEST", "hello world")
	val, err = vm.Eval(`
    s = Getenv("WS_TEST");
    s;
  `)
	if err != nil {
		t.Errorf(`vm.Eval() error, %s`, err)
	}
	s, err := val.Export()
	if err != nil {
		t.Errorf("Can't export 's' from JS eval, %s", err)
	}
	if s != "hello world" {
		t.Errorf("Expected JS to return s of 'hello world', got %s", s)
	}
}

func TestJSPathToRoute(t *testing.T) {
	os.Setenv("WS_JSDOCS", "jsdocs")
	cfg := new(Configuration)
	cfg.Getenv()
	p := "jsdocs/helloworld.js"
	r, err := JSPathToRoute(p, cfg)
	if err != nil {
		t.Errorf("JSPathToRoute() error, %s", err)
	}
	if r != "/helloworld" {
		t.Errorf("Failed converting path to route /helloworld, %s", r)
	}
	p = "jsdocs/api/search.js"
	r, err = JSPathToRoute(p, cfg)
	if err != nil {
		t.Errorf("JSPathToRoute() error, %s", err)
	}
	if r != "/api/search" {
		t.Errorf("Failed converting path to route /api/search, %s", r)
	}
}

func TestMain(m *testing.M) {
	// Clean up any stale test data
	os.RemoveAll("testout")
	// Run tests
	exitCode := m.Run()
	if exitCode == 0 {
		log.Println("Cleaning up")
		os.RemoveAll("testout")
	} else {
		log.Println("Output saved in ./testout")
	}
	os.Exit(exitCode)
}
