//
// Package cfg is a general configuration methods for ws.go
//
package cfg

import (
	"../keygen"
	"../prompt"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"strings"
)

// Cfg is an application's configuration - who started the process, port assignment
// configuration settings, etc.
type Cfg struct {
	Username string
	Hostname string
	Port     int
	UseTLS   bool
	Docroot  string
	Cert     string
	Key      string
	Otto     bool
	OttoPath string
}

// Configure an application.
func Configure(docroot string, hostname string, port int, useTLS bool, cert string, key string, otto bool, ottoPath string) (*Cfg, error) {
	wsUser, err := user.Current()
	if err != nil {
		return nil, err
	}

	// Normalize docroot
	if strings.HasPrefix(docroot, "/") == false {
		cleanDocroot, err := filepath.Abs(path.Join("./", docroot))
		if err != nil {
			return nil, fmt.Errorf("Can't expand docroot %s: %s\n", docroot, err)
		}
		docroot = cleanDocroot
	}
	// Normalize ottoPath
	if otto == true && strings.HasPrefix(ottoPath, "/") == false {
		cleanOttoPath, err := filepath.Abs(path.Join("./", ottoPath))
		if err != nil {
			return nil, fmt.Errorf("Can't expand ottoPath %s: %s\n", ottoPath, err)
		}
		ottoPath = cleanOttoPath
	} else if otto == false {
		ottoPath = ""
	}
	return &Cfg{
		Username: wsUser.Username,
		Hostname: hostname,
		Port:     port,
		Docroot:  docroot,
		UseTLS:   useTLS,
		Cert:     cert,
		Key:      key,
		Otto:     otto,
		OttoPath: ottoPath}, nil
}

// InitProject initializes a basic project structure (e.g. creates static, dynamic, README.md, etc/config.sh)
func InitProject() error {
	var (
		host         string
		port         string
		projectName  string
		authorName   string
		description  string
		docroot      string
		useTLS       bool
		certFilename string
		keyFilename  string
		otto         bool
		ottoPath     string
		config       string
		OK           bool
		err          error
	)

	OK = false
	useTLS = false
	docroot = "static"
	otto = false
	ottoPath = "dynamic"
	for OK == false {
		projectName = prompt.Question("Name of Project (e.g. Big Reptiles): ", "Big Reptiles")
		authorName = prompt.Question("Name of Author(s) (e.g. Mr. Lizard): ", "Mr. Lizard")
		description = prompt.Question("Description (e.g A demo project): ", "A Demo Project")
		host = prompt.Question("Hostname (e.g. localhost): ", "localhost")
		port = prompt.Question("Post (e.g. 8000): ", "8000")
		docroot = prompt.Question("Document root for static files (e.g. ./static): ", docroot)
		config = prompt.Question("Directory to use for configuration files (e.g. ./etc): ", "etc")
		otto = prompt.YesNo("Turn on Otto Engine? ")
		if otto == true {
			ottoPath = prompt.Question("Path to Otto Engine routes (e.g. ./dyanmic): ", ottoPath)
		}
		fmt.Printf("Configuration choosen\nProject: %s\nAuthor(s): %s\nDescription: %s\nDocroot: %s\n",
			projectName,
			authorName,
			description,
			docroot)
		fmt.Printf("Turn on Otto Engine: %v %s\n", otto, ottoPath)

		// Display current settings
		OK = prompt.YesNo("Is this OK? ")
	}

	useTLS = prompt.YesNo("Configure for SSL support? ")
	if useTLS == true {
		// Defer handling of SSL questions to keygen.Keygen()
		certFilename, keyFilename, err = keygen.Keygen(path.Join(config, "ssl"), "cert.pem", "key.pem")
		if err != nil {
			return err
		}
	}

	fmt.Printf("Creating %s\n", config)
	err = os.MkdirAll(config, 0770)
	if err != nil {
		return err
	}

	fmt.Printf("Creating %s\n", docroot)
	err = os.MkdirAll(docroot, 0775)
	if err != nil {
		return err
	}

	if otto == true {
		fmt.Printf("Creating %s\n", ottoPath)
		err = os.MkdirAll(ottoPath, 0775)
		if err != nil {
			return err
		}
	}

	fmt.Println("Creating README.md")
	readme := fmt.Sprintf("\n\n# %s\n\nBy %s\n\n## Overview\n\n%s\n\n", projectName, authorName, description)
	err = ioutil.WriteFile("README.md", []byte(readme), 0664)
	if err != nil {
		return err
	}

	fmt.Printf("Creating %s\n", path.Join(config, "config.sh"))
	configEnvironment := fmt.Sprintf("#!/bin/bash\n# %s configuration\n# Source this file before running ws\n#\n\nexport WS_HOST=%q\nexport WS_PORT=%q\nexport WS_DOCROOT=%q\nexport WS_OTTO=%v\nexport WS_OTTO_PATH=%q\n", projectName, host, port, docroot, otto, ottoPath)

	if useTLS == true {
		configEnvironment += fmt.Sprintf("\nexport TLS=true\nexport WS_CERT=%s\nexport WS_KEY=%s\n\n", certFilename, keyFilename)
	}
	err = ioutil.WriteFile(path.Join(config, "config.sh"), []byte(configEnvironment), 0770)
	if err != nil {
		return err
	}

	fmt.Printf("Creating %s/index.html\n", docroot)
	index := fmt.Sprintf(`<!DOCTYPE html>
<html>
    <head>
        <title>%s</title>
    </head>
    <body>
        <h1>%s</h1>
        <p>by %s</p>
        <div>%s</div>
    </body>
</html>
`, projectName, projectName, authorName, description)
	err = ioutil.WriteFile(path.Join(docroot, "index.html"), []byte(index), 0664)
	if err != nil {
		return err
	}

	if otto == true {
		testJs := `/**
 * test.js - an example Otto Engine route handler
 */
/*jslint browser: false, indent: 4 */
/*global Request, Response */
(function (req, res) {
    res.setHeader("Content-Type", "text/plain")
    res.setContent("Hello World!")
}(Request, Response))
`
		fmt.Printf("Creating %s/test.js\n", ottoPath)
		err = ioutil.WriteFile(path.Join(ottoPath, "test.js"), []byte(testJs), 0664)
		if err != nil {
			return err
		}
	}

	fmt.Println("Setup completed.")
	return nil
}
