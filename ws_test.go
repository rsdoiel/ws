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
	if cfg.DocRoot != "" {
		t.Errorf("cfg.DocRoot should be empty, %s", cfg.DocRoot)
	}
	if cfg.SSLKey != "" {
		t.Errorf("cfg.SSLKey should be empty, %s", cfg.SSLKey)
	}
	if cfg.SSLCert != "" {
		t.Errorf("cfg.SSLCert should be empty, %s", cfg.SSLCert)
	}

	os.Setenv("WS_URL", "https://example.org:8001")
	os.Setenv("WS_DOCROOT", "htdocs")
	os.Setenv("WS_SSL_KEY", "etc/ssl/site.key")
	os.Setenv("WS_SSL_CERT", "etc/ssl/site.crt")
	err := cfg.Getenv()
	if err != nil {
		t.Errorf("cfg.Getenv() error, %s", err)
	}
	if cfg.URL.Host != "example.org:8001" {
		t.Errorf("cfg.URL.Host != example.org:8001, %s", cfg.URL.Host)
	}
	if cfg.DocRoot != "htdocs" {
		t.Errorf("cfg.DocRoot != htdocs, %s", cfg.DocRoot)
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
	docRoot := "/www/htdocs"
	sslkey := "/etc/ssl/site.key"
	sslcert := "/etc/ssl/site.crt"

	// Set some example configuration
	os.Setenv("WS_URL", u)
	os.Setenv("WS_DOCROOT", docRoot)
	os.Setenv("WS_SSL_KEY", sslkey)
	os.Setenv("WS_SSL_CERT", sslcert)
	cfg.Getenv()

	expected := fmt.Sprintf(`#!/bin/bash
# generated %d-%02d-%02d by ws version %s
export WS_URL=%q
export WS_DOCROOT=%q
export WS_SSL_KEY=%q
export WS_SSL_CERT=%q
`, yr, mn, dy, Version, u, docRoot, sslkey, sslcert)

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
	docRoot := "testout/https/htdocs"
	sslkey := "testout/https/etc/ssl/site.key"
	sslcert := "testout/https/etc/ssl/site.crt"
	// Set some example configuration
	os.Setenv("WS_URL", u)
	os.Setenv("WS_DOCROOT", docRoot)
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
		cfg.DocRoot,
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
	docRoot = "testout/http/htdocs"
	sslkey = "testout/http/etc/ssl/site.key"
	sslcert = "testout/http/etc/ssl/site.crt"
	// Set some example configuration
	os.Setenv("WS_URL", u)
	os.Setenv("WS_DOCROOT", docRoot)
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
		cfg.DocRoot,
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
	docRoot := "testout/www-01/htdocs"
	sslkey := "testout/www-01/ssl/site.key"
	sslcert := "testout/www-01/ssl/site.crt"

	// Set some example configuration
	os.Setenv("WS_URL", u)
	os.Setenv("WS_DOCROOT", docRoot)
	os.Setenv("WS_SSL_KEY", sslkey)
	os.Setenv("WS_SSL_CERT", sslcert)
	cfg.Getenv()
	cfg.InitializeProject()

	if err := cfg.Validate(); err != nil {
		t.Errorf("Should have been valid %+v, %s", cfg, err)
	}

	u = "http://example.org"
	docRoot = "testout/www-02/htdocs"
	sslkey = ""
	sslcert = ""
	os.Setenv("WS_URL", u)
	os.Setenv("WS_DOCROOT", docRoot)
	os.Setenv("WS_SSL_KEY", sslkey)
	os.Setenv("WS_SSL_CERT", sslcert)
	cfg.Getenv()
	cfg.InitializeProject()

	if err := cfg.Validate(); err != nil {
		t.Errorf("Should have been valid %+v, %s", cfg, err)
	}

	cfg = new(Configuration)
	os.Setenv("WS_URL", "https://example.org")
	os.Unsetenv("WS_SSL_KEY")
	os.Unsetenv("WS_SSL_CERT")
	cfg.Getenv()
	if err := cfg.Validate(); err == nil {
		t.Errorf("https without ssl files, should have been invalid %+v, %s", cfg, err)
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

func TestSetDefaults(t *testing.T) {
	cfg := new(Configuration)
	cfg.SetDefaults()
	u := cfg.URL.String()
	d := cfg.DocRoot
	if u != "http://localhost:8000" {
		t.Errorf("expected http://localhost:8000 got %q", u)
	}
	if d != "." {
		t.Errorf("expected '.' got %q", d)
	}
}
