/**
 * ws.go - A light weight webserver for static content development.
 * Supports both http and https protocols.
 *
 * @author R. S. Doiel, <rsdoiel@yahoo.com>
 * copyright (c) 2014
 * Released under the BSD 2-Clause License
 */
package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"os/user"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// command line parameters that override environment variables
var (
	cli_use_tls = flag.Bool("tls", false, "Turn on TLS (https) support with true, off with false (default is false)")
	cli_docroot = flag.String("docroot", "", "Path to the docment root")
	cli_port = flag.Int("port", 0, "Port number to listen on")
	cli_cert = flag.String("cert", "", "Path to your TLS cert.pem")
	cli_key = flag.String("key", "", "Path to your TLS key.pem")
)

var ErrHelp = errors.New("flag: Help requested")

var Usage = func() {
	fmt.Fprintf(os.Stderr, "Usage of %s\n", os.Args[0])
	flag.PrintDefaults()
}

// Application's profile - who started the process, port assignment
// configuration settings, etc.
type Profile struct {
	Username string
	Hostname string
	Port     string
	Use_TLS  bool
	Docroot  string
	Cert     string
	Key      string
}

func LoadProfile(cli_docroot string, cli_port int, cli_use_tls bool, cli_cert string, cli_key string) (*Profile, error) {
	ws_user, err := user.Current()
	if err != nil {
		return nil, err
	}
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}
	port := "8000"
	use_tls := false

	cert, err := ConfigPathTo("cert.pem")
	if err != nil {
		return nil, err
	}
	key, err := ConfigPathTo("key.pem")
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
	if cli_port != 0 {
		port = strconv.Itoa(cli_port)
	}
	if cli_cert != "" {
		cert = cli_cert
	}
	if cli_key != "" {
		key = cli_key
	}

	return &Profile{
		Username: ws_user.Username,
		Hostname: hostname,
		Port:     port,
		Docroot:  path.Join(docroot),
		Use_TLS:  use_tls,
		Cert:     cert,
		Key:      key}, nil
}

func Log(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func webserver(profile *Profile) error {
	// Define a simple static file server
	//http.Handle("/", http.FileServer(http.Dir(profile.Docroot)))

	// Restricted FileService excluding dot files and directories
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var hasDotPath = regexp.MustCompile(`\/\.`)
		unclean_path := r.URL.Path
		if !strings.HasPrefix(unclean_path, "/") {
			unclean_path = "/" + unclean_path
		}
		clean_path := path.Clean(unclean_path)
		r.URL.Path = clean_path
		resolved_path := path.Clean(path.Join(profile.Docroot, clean_path))
		if hasDotPath.MatchString(clean_path) {
			log.Printf("Not Authorized (401) %s\n", clean_path)
			http.Error(w, "Not Authorized", 401)
		} else if !strings.HasPrefix(resolved_path, profile.Docroot) {
			log.Printf("Not Found (404) %s\n", resolved_path)
			http.NotFound(w, r)
		} else {
			http.ServeFile(w, r, resolved_path)
		}
	})

	if profile.Use_TLS == false {
		log.Printf("\n Docroot:   %s\n    Port:   %s\n"+
			"  Run as:   %s\n\n",
			profile.Docroot, profile.Port,
			profile.Username)
		log.Println("Starting http://" + net.JoinHostPort(profile.Hostname, profile.Port))

		// Now start up the server and log transactions
		return http.ListenAndServe(net.JoinHostPort(profile.Hostname, profile.Port), Log(http.DefaultServeMux))
	}
	log.Printf("\n    Cert:   %s\n     Key:   %s\n"+
		" Docroot:   %s\n    Port:   %s\n"+
		"  Run as:   %s\n\n",
		profile.Cert, profile.Key,
		profile.Docroot, profile.Port,
		profile.Username)
	log.Println("Starting https://" + net.JoinHostPort(profile.Hostname, profile.Port))

	// Now start up the server and log transactions
	return http.ListenAndServeTLS(net.JoinHostPort(profile.Hostname, profile.Port), profile.Cert, profile.Key, Log(http.DefaultServeMux))
}

func ConfigPath() (string, error) {
	home := os.Getenv("HOME")
	config_path := home + "/etc/ws"
	err := os.MkdirAll(config_path, 0700)
	if err != nil {
		return "", err
	}
	return config_path, nil
}

func ConfigPathTo(filename string) (string, error) {
	ws_path, err := ConfigPath()
	if err != nil {
		return "", err
	}
	return ws_path + "/" + filename, nil
}

func init() {
	rand.Seed(time.Now().Unix())
}

func main() {
	flag.Parse()

	profile, _ := LoadProfile(*cli_docroot, *cli_port, *cli_use_tls, *cli_cert, *cli_key)
	err := webserver(profile)
	if err != nil {
		log.Fatal(err)
	}
}
