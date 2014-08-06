/**
 * cfg.go - general configuration methods for ws.go
 */
package cfg

import (
	"../keygen"
	"../prompt"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"strings"
)

// Application's configuration - who started the process, port assignment
// configuration settings, etc.
type Cfg struct {
	Username  string
	Hostname  string
	Port      int
	UseTLS   bool
	Docroot   string
	Cert      string
	Key       string
	Otto      bool
	OttoPath string
}

// Configure this application.
func Configure(docroot string, hostname string, port int, use_tls bool, cert string, key string, otto bool, otto_path string) (*Cfg, error) {
	ws_user, err := user.Current()
	if err != nil {
		return nil, err
	}

	// Normalize docroot
	if strings.HasPrefix(docroot, "/") == false {
		clean_docroot, err := filepath.Abs(path.Join("./", docroot))
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Can't expand docroot %s: %s\n", docroot, err))
		}
		docroot = clean_docroot
	}
	// Normalize otto_path
	if otto == true && strings.HasPrefix(otto_path, "/") == false {
		clean_otto_path, err := filepath.Abs(path.Join("./", otto_path))
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Can't expand otto_path %s: %s\n", otto_path, err))
		}
		otto_path = clean_otto_path
	} else if otto == false {
		otto_path = ""
	}
	return &Cfg{
		Username:  ws_user.Username,
		Hostname:  hostname,
		Port:      port,
		Docroot:   docroot,
		UseTLS:   use_tls,
		Cert:      cert,
		Key:       key,
		Otto:      otto,
		OttoPath: otto_path}, nil
}

// InitProject - initializes a basic project structure (e.g. creates static, dynamic, README.md, etc/config.sh)
func InitProject() error {
	var (
		host          string
		port          string
		project_name  string
		author_name   string
		description   string
		docroot       string
		use_tls       bool
		cert_filename string
		key_filename  string
		otto          bool
		otto_path     string
		config        string
		OK            bool
		err           error
	)

	OK = false
	use_tls = false
	docroot = "static"
	otto = false
	otto_path = "dynamic"
	for OK == false {
		project_name = prompt.PromptString("Name of Project: (e.g. Big Reptiles)", "Big Reptiles")
		author_name = prompt.PromptString("Name of Author(s): (e.g. Mr. Lizard)", "Mr. Lizard")
		description = prompt.PromptString("Description (e.g A demo project)", "A Demo Project")
		host = prompt.PromptString("Hostname (e.g. localhost)", "localhost")
		port = prompt.PromptString("Post (e.g. 8000)", "8000")
		docroot = prompt.PromptString("Document root for static files (e.g. ./static)", docroot)
		config = prompt.PromptString("Directory to use for configuration files (e.g. ./etc)", "etc")
		otto = prompt.YesNo("Turn on Otto Engine?")
		if otto == true {
			otto_path = prompt.PromptString("Path to Otto Engine routes (e.g. ./dyanmic)", otto_path)
		}
		fmt.Printf("Configuration choosen\nProject: %s\nAuthor(s): %s\nDescription: %s\nDocroot: %s\n",
			project_name,
			author_name,
			description,
			docroot)
		fmt.Printf("Turn on Otto Engine: %v %s\n", otto, otto_path)

		// Display current settings
		OK = prompt.YesNo("Is this OK?")
	}

	use_tls = prompt.YesNo("Configure for SSL support?")
	if use_tls == true {
		// Defer handling of SSL questions to keygen.Keygen()
		cert_filename, key_filename, err = keygen.Keygen(path.Join(config, "ssl"), "cert.pem", "key.pem")
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
		fmt.Printf("Creating %s\n", otto_path)
		err = os.MkdirAll(otto_path, 0775)
		if err != nil {
			return err
		}
	}

	fmt.Println("Creating README.md")
	readme := fmt.Sprintf("\n# %s\n\nBy %s\n\n## Overview\n%s\n\n", project_name, author_name, description)
	err = ioutil.WriteFile("README.md", []byte(readme), 0664)
	if err != nil {
		return err
	}

	fmt.Printf("Creating %s\n", path.Join(config, "config.sh"))
	config_environment := fmt.Sprintf("#!/bin/bash\n# %s configuration\n# Source this file before running ws\n#\n\nexport WS_HOST=%q\nexport WS_PORT=%q\nexport WS_DOCROOT=%q\nexport WS_OTTO=%v\nexport WS_OTTO_PATH=%q\n", project_name, host, port, docroot, otto, otto_path)

	if use_tls == true {
		config_environment += fmt.Sprintf("\nexport TLS=true\nexport WS_CERT=%s\nexport WS_KEY=%s\n\n", cert_filename, key_filename)
	}
	err = ioutil.WriteFile(path.Join(config, "config.sh"), []byte(config_environment), 0770)
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
`, project_name, project_name, author_name, description)
	err = ioutil.WriteFile(path.Join(docroot, "index.html"), []byte(index), 0664)
	if err != nil {
		return err
	}

	if otto == true {
		test_js := `/**
 * test.js - an example Otto Engine route handler
 */
/*jslint browser: false, indent: 4 */
/*global Request, Response */
(function (req, res) {
    res.setHeader("Content-Type", "text/plain")
    res.setContent("Hello World!")
}(Request, Response))
`
		fmt.Printf("Creating %s/test.js\n", otto_path)
		err = ioutil.WriteFile(path.Join(otto_path, "test.js"), []byte(test_js), 0664)
		if err != nil {
			return err
		}
	}

	fmt.Println("Setup completed.")
	return nil
}
