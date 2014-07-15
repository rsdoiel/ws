/**
 * app.go - general configuration methods for ws.go
 */
package app

import (
    "fmt"
    "io"
    "bufio"
    "os"
    "os/user"
    "strings"
    "path"
    "path/filepath"
    "log"
    "errors"
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
func Init() {
    var line string

    // Get the name of the project
    project_name := "Big Reptiles"
    author_name := "Mr. Lizard"
    description := "A demo project"
    config := "etc"
    use_tls := false
    docroot := "static"
    otto := false
    otto_path := "dynamic"

    fmt.Printf("Name of project: %s", project_name)
    line, _ = bufio.ReadString('\n')
    line = strings.TrimSpace(line)
    if line != "" {
        project_name = line
    }

    fmt.Printf("Author(s): %s", author_name)
    line, _ = bufio.ReadString('\n')
    line = strings.TrimSpace(line)
    if line != "" {
        author_name = line
    }

    fmt.Printf("Description: %s", description)
    line, _ = bufio.ReadString('\n')
    line = strings.TrimSpace(line)
    if line != "" {
        description = line
    }

    fmt.Printf("Config directory: %s", config)
    line, _ = bufio.ReadString('\n')
    line = strings.TrimSpace(line)
    if line != "" {
        config = line
    }

    fmt.Print("Run with SSL support? (y/n)")
    line, _ =  bufio.ReadString('\n')
    line = strings.ToLowerCase(strings.Trim(line))
    if line == "y" {
        use_tls = true
    }

    fmt.Printf("Docroot: %s", docroot)
    line, _ := bufio.ReadString('\n')
    line = strings.TrimSpace(line)
    if line != "" {
        docroot = line
    }

    fmt.Print("Turn Otto Engine on? (y/n)")
    line, _ := bufio.ReadString('\n')
    line = strins.ToLowerCase(strings.TrimSpace(line))
    if line == "y" {
        otto = true
        fmt.Printf("Otto path: %s", otto_path)
        line, _ := bufio.ReadString('\n')
        line = strings.TrimSpace(line)
        if line != "" {
            otto_path = line
        }
    }

    os.MkdirAll(config, 0770)
    os.MkdirAll(docroot, 0775)
    os.MkdirAll(otto_path, 0775)

    readme := fmt.Sprintf("\n# %s\n\n By %s\n\n %s\n\n", project_name, author_name, description)
    ioutil.WriteFile("README.md", []byte(readme), 0664)

    config_environment := fmt.Sprintf("#!/bin/bash\n# %s configuration\n# Source this file before running ws\n\nexport WS_DOCROOT=%q\nexport WS_OTTO=%s\nexport WS_OTTO_PATH=%q\n")

    if use_tls == true {
        config_environment += fmt.Sprintf("\nexport TLS=true\nexport WS_CERT=%s/cert.pem\nexport WS_KEY=%s/key.pem\n\n", config)
        fmt.Println("Run ws -keygen and set the directory to match %s with cert filename of cert.pem and key filename of key.pem")
    } 
    ioutil.WriteFile(path.Join(config, "config.sh"), config_environment, 0770)
}

