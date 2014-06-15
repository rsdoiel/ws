/**
 * ws.go - A simple ad-hoc web server like httpster but providing SSL
 * support by default.
 * @author R. S. Doiel, <rsdoiel@yahoo.com>
 * copyright (c) 2014
 * Released under the BSD 2-Clause License
 */
package main

import (
	"flag"
	"math/rand"
	"os"
	"os/user"
	"time"
    "strconv"
	"net"
	"net/http"
	"log"
)

type Profile struct {
	Username string
	Hostname string
	Port     string
    Docroot  string
	Cert	 string
	Key	 string
}

func LoadProfile(cli_docroot string, cli_port int, cli_cert string, cli_key string) (*Profile, error) {
	ws_user, user_error := user.Current()
	if user_error != nil {
		return nil, user_error
	}
	hostname, hostname_error := os.Hostname()
	if hostname_error != nil {
		return nil, hostname_error
	}
	cert, err := ConfigPathTo("cert.pem")
	if err != nil {
		return nil, err
	}
	key, err := ConfigPathTo("key.pem")
	if err != nil {
		return nil, err
	}
	port := "8443"
    docroot := "./"

	// now overwrite with any environment settings found. 
	env_host := os.Getenv("WS_HOST")
	env_port := os.Getenv("WS_PORT")
	env_cert := os.Getenv("WS_CERT")
	env_key := os.Getenv("WS_KEY")
    env_docroot := os.Getenv("WS_DOCROOT")
	if env_host != "" {
		hostname = env_host
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
        Docroot:  docroot,
		Cert: 	  cert,
		Key:      key}, nil
}

func Log(handler http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
        handler.ServeHTTP(w, r)})
}


func ssl_webserver(profile *Profile) error {
    // Define a simple static file server
    http.Handle("/", http.FileServer(http.Dir(profile.Docroot)))

	// NOTE: we ONLY supports https protocol to limit leaking data 
    // over wireless and other exposed networks
    log.Printf("\n    Cert:   %s\n     Key:   %s\n" + 
               " Docroot:   %s\n    Port:   %s\n" +
	           "  Run as:   %s\n\n",
                profile.Cert, profile.Key,
                profile.Docroot, profile.Port,
                profile.Username)
	log.Println("Starting " + net.JoinHostPort(profile.Hostname, profile.Port))

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
    cli_docroot := flag.String("docroot", "", "Path to the docment root")
    cli_port := flag.Int("port", 0, "Port number to listen on")
    cli_cert := flag.String("cert", "", "Path to your SSL cert.pem")
    cli_key  := flag.String("key", "", "Path to your SSL key.pem")

    flag.Parse()

	ws_user, _ := LoadProfile(*cli_docroot, *cli_port, *cli_cert, *cli_key)
	err := ssl_webserver(ws_user)
	if err != nil {
		log.Fatal(err)
	}
}
