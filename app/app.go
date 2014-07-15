/**
 * app.go - general configuration methods for ws.go
 */
package app

import (
    "../prompt"
    "../keygen"
    "fmt"
    "os"
    "os/user"
    "strings"
    "path"
    "path/filepath"
    "log"
    "io/ioutil"
)

// Application's profile - who started the process, port assignment
// configuration settings, etc.
type Profile struct {
	Username  string
	Hostname  string
	Port      string
	Use_TLS   bool
	Docroot   string
	Cert      string
	Key       string
	Otto      bool
	Otto_Path string
}


// LoadProfile - load an application profile from both the environment
// and cli options.
func LoadProfile(cli_docroot string, cli_host string, cli_port string, cli_use_tls bool, cli_cert string, cli_key string, cli_otto bool, cli_otto_path string) (*Profile, error) {
	ws_user, err := user.Current()
	if err != nil {
		return nil, err
	}
    //Q: Should I really default to localhost? instead of Hostname?
    /*
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}
    */
    hostname := "localhost"
	port := "8000"
	use_tls := false
	otto := false
	otto_path := ""

	cert := ""
	key := ""
	if err != nil {
		return nil, err
	}
	docroot, _ := os.Getwd()

	// now overwrite with any environment settings found.
	env_host := os.Getenv("WS_HOST")
	env_port := os.Getenv("WS_PORT")
	env_use_tls := os.Getenv("WS_TLS")
	env_cert := os.Getenv("WS_CERT")
	env_key := os.Getenv("WS_KEY")
	env_docroot := os.Getenv("WS_DOCROOT")
	env_otto := os.Getenv("WS_OTTO")
	env_otto_path := os.Getenv("WS_OTTO_PATH")

    // merge the environment settings
	if env_host != "" {
		hostname = env_host
	}
	if env_use_tls == "true" {
		use_tls = true
		port = "8443"
	}
	if env_port != "" {
		port = env_port
	}
	if env_docroot != "" {
		docroot = env_docroot
	}
	if env_cert != "" {
		cert = env_cert
	}
	if env_key != "" {
		key = env_key
	}
	if env_otto == "true" {
		otto = true
	}
    if env_otto == "false" {
        otto = false
    }
	if env_otto_path != "" {
		otto_path = env_otto_path
	}

	// Finally resolve any command line overrides
	if cli_docroot != "" {
		docroot = cli_docroot
	}
	if cli_use_tls == true {
		use_tls = true
		if env_port == "" {
			port = "8443"
		}
	}
	if cli_host != "" {
		hostname = cli_host
	}
	if cli_port != "" {
		port = cli_port
	}
	if cli_cert != "" {
		cert = cli_cert
	}
	if cli_key != "" {
		key = cli_key
	}
	if cli_otto == true {
		otto = true
	}
	if cli_otto_path != "" {
        otto = true
		otto_path = cli_otto_path
	}

	// If TLS is false then don't expose the location of the cert/key
	if use_tls == false {
		cert = ""
		key = ""
	}

	// Normalize docroot
	if strings.HasPrefix(docroot, "/") == false {
		clean_docroot, err := filepath.Abs(path.Join("./", docroot))
		if err != nil {
			log.Fatalf("Can't expand docroot %s: %s\n", docroot, err)
		}
		docroot = clean_docroot
	}
	// Normalize otto_path
    if otto == true && strings.HasPrefix(otto_path, "/") == false {
	    clean_otto_path, err := filepath.Abs(path.Join("./", otto_path))
	    if err != nil {
		    log.Fatalf("Can't expand otto_path %s: %s\n", otto_path, err)
	    }
	    otto_path = clean_otto_path
    } else if otto == false {
        otto_path = ""
    }
	return &Profile{
		Username:  ws_user.Username,
		Hostname:  hostname,
		Port:      port,
		Docroot:   docroot,
		Use_TLS:   use_tls,
		Cert:      cert,
		Key:       key,
		Otto:      otto,
		Otto_Path: otto_path}, nil
}

// Init - initializes a basic project structure (e.g. creates static, dynamic, README.md, etc/config.sh)
func Init() error {
    var (
        project_name string
        author_name string
        description string
        docroot string
        use_tls bool
        cert_filename string
        key_filename string
        otto bool
        otto_path string
        config string
        OK  bool
        err error
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
    config_environment := fmt.Sprintf("#!/bin/bash\n# %s configuration\n# Source this file before running ws\n\nexport WS_DOCROOT=%q\nexport WS_OTTO=%v\nexport WS_OTTO_PATH=%q\n", project_name, docroot, otto, otto_path)

    if use_tls == true {
        config_environment += fmt.Sprintf("\nexport TLS=true\nexport WS_CERT=%s\nexport WS_KEY=%s\n\n", cert_filename, key_filename)
    } 
    err = ioutil.WriteFile(path.Join(config, "config.sh"), []byte(config_environment), 0770)
    if err != nil {
        return err
    }

    fmt.Printf("Creating %s/index.html", docroot)
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

    test_js := `
//
// test.js - an example Otto Engine route handler
//
(function (req, res) {
    res.setHeader("Content-Type", "text/plain")
    res.setContent("Hello World!")
}(Request, Response)
`
    fmt.Printf("Creating %s/test.js", otto_path)
    err = ioutil.WriteFile(path.Join(otto_path, "test.js"), []byte(test_js), 0664)
    if err != nil {
        return err
    }


    //FIXME: add dynamic/test.js example.
    fmt.Println("Setup completed.")
    return nil
}

