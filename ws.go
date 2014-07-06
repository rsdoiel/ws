/**
 * ws.go - A light weight webserver for static content
 * development and prototyping route based web API.
 *
 * Supports both http and https protocols. Dynamic route
 * processing available via Otto JavaScript virtual machines.
 *
 * @author R. S. Doiel, <rsdoiel@yahoo.com>
 * copyright (c) 2014
 * All rights reserved.
 * @license BSD 2-Clause License
 */
package main

import (
	"./ottoengine"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"flag"
	"fmt"
	"log"
	"math/big"
	"net"
	"net/url"
	"net/http"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// variables for keygen
var (
	cli_keygen   = flag.Bool("keygen", false, "Generate TLS ceriticates and keys")
	cli_ssl_host = flag.String("ssl-host", "", "Comma-separated hostnames and IPs to generate a certificate for")
	validFrom    = flag.String("start-date", "", "Creation date formatted as Jan 1 15:04:05 2011")
	validFor     = flag.Duration("duration", 365*24*time.Hour, "Duration that certificate is valid for")
	organization = flag.String("organization", "Acme Co.", "Organization used to sign certificate")
	isCA         = flag.Bool("ca", false, "whether this cert should be its own Certificate Authority")
	rsaBits      = flag.Int("rsa-bits", 2048, "Size of RSA key to generate")
)

// command line parameters that override environment variables
var (
	cli_use_tls   = flag.Bool("tls", false, "Turn on TLS (https) support with true, off with false (default is false)")
	cli_docroot   = flag.String("docroot", "", "Path to the document root")
	cli_host      = flag.String("host", "", "Hostname http(s) server to listen for")
	cli_port      = flag.String("port", "", "Port number to listen on")
	cli_cert      = flag.String("cert", "", "Filename to your SSL cert.pem")
	cli_key       = flag.String("key", "", "Filename to your SSL key.pem")
	cli_otto      = flag.Bool("otto", false, "Enable experimental Otto JS VM support")
	cli_otto_path = flag.String("otto-path", "", "The search path for runable Otto JS Programs.")
)

var ErrHelp = errors.New("flag: Help requested")

var Usage = func() {
	fmt.Fprintf(os.Stderr, "Usage of %s\n", os.Args[0])
	flag.PrintDefaults()
}

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

func LoadProfile(cli_docroot string, cli_host string, cli_port string, cli_use_tls bool, cli_cert string, cli_key string, cli_otto bool, cli_otto_path string) (*Profile, error) {
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
	if strings.HasPrefix(otto_path, "/") == false {
		clean_otto_path, err := filepath.Abs(path.Join("./", otto_path))
		if err != nil {
			log.Fatalf("Can't expand otto_path %s: %s\n", otto_path, err)
		}
		otto_path = clean_otto_path
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

func log_response(status int, err string, method string, url *url.URL, proto, referrer, user_agent string) {
	log.Printf("{\"response\": %d, \"status\": %q, %q: %q, \"protocol\": %q, \"referrer\": %q, \"user-agent\": %q}\n",
        status,
        err,
        method,
        url,
        proto,
        referrer,
        user_agent)
}

func log_request(method string, url *url.URL, proto, referrer, user_agent string) {
	log.Printf("{\"request\": true, %q: %q, \"protocol\": %q, \"referrer\": %q, \"user-agent\": %q}\n",
        method,
        url.String(),
        proto,
        referrer,
        user_agent)
}

func request_log(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log_request(r.Method, r.URL, r.Proto, r.Referer(), r.UserAgent())
        handler.ServeHTTP(w, r)
	})
}

func Webserver(profile *Profile) error {
	// If otto is enabled add routes and handle them.
	if profile.Otto == true {
		otto_path, err := filepath.Abs(profile.Otto_Path)
		if err != nil {
			log.Fatalf("Can't read %s: %s\n", profile.Otto_Path, err)
		}
		programs, err := ottoengine.Load(otto_path)
		if err != nil {
			log.Fatalf("Load error: %s\n", err)
		}
		ottoengine.AddRoutes(programs)
	}

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
        _, err := os.Stat(resolved_path)
        if hasDotPath.MatchString(clean_path) == true || 
                strings.HasPrefix(resolved_path, profile.Docroot) == false ||
                os.IsPermission(err) == true {
			log_response(401, "Not Authorized", r.Method, r.URL, r.Proto, r.Referer(), r.UserAgent())
			http.Error(w, "Not Authorized", 401)
		} else if os.IsNotExist(err) == true {
			log_response(404, "Not Found", r.Method, r.URL, r.Proto, r.Referer(), r.UserAgent())
			http.NotFound(w, r)
		} else if err == nil {
			log_response(200, "OK", r.Method, r.URL, r.Proto, r.Referer(), r.UserAgent())
			http.ServeFile(w, r, resolved_path)
        } else {
            // Easter egg
			log_response(418, "I'm a teapot", r.Method, r.URL, r.Proto, r.Referer(), r.UserAgent())
			http.Error(w, "I'm a teapot", 418)
        }
	})

	// Now start up the server and log transactions
	if profile.Use_TLS == true {
		if profile.Cert == "" || profile.Key == "" {
			log.Fatalf("TLS set true but missing key or certificate")
		}
		log.Println("Starting https://" + net.JoinHostPort(profile.Hostname, profile.Port))
		return http.ListenAndServeTLS(net.JoinHostPort(profile.Hostname, profile.Port), profile.Cert, profile.Key, request_log(http.DefaultServeMux))
	}
	log.Println("Starting http://" + net.JoinHostPort(profile.Hostname, profile.Port))
	// Now start up the server and log transactions
	return http.ListenAndServe(net.JoinHostPort(profile.Hostname, profile.Port), request_log(http.DefaultServeMux))
}

func keygen(profile *Profile) error {
	home := os.Getenv("HOME")
	certFilename := profile.Cert
	if certFilename == "" {
		certFilename = path.Join(home, "etc/ws/cert.pem")
	}
	keyFilename := profile.Key
	if keyFilename == "" {
		keyFilename = path.Join(home, "etc/ws/key.pem")
	}

	hostnames := profile.Hostname
	if *cli_ssl_host != "" {
		hostnames = *cli_ssl_host
	}
	if hostnames == "" {
		log.Fatalf("Missing required -ssl-host parameter")
	}

	log.Printf("\n\n"+
		" Cert: %s\n"+
		"  Key: %s\n"+
		" Host: %s\n"+
		" Organization: %v\n"+
		"\n\n",
		certFilename,
		keyFilename,
		hostnames,
		*organization)

	priv, err := rsa.GenerateKey(rand.Reader, *rsaBits)
	if err != nil {
		log.Fatalf("failed to generate private key: %s", err)
	}

	var notBefore time.Time
	if len(*validFrom) == 0 {
		notBefore = time.Now()
	} else {
		notBefore, err = time.Parse("Jan 2 15:04:05 2006", *validFrom)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to parse creation date: %s\n", err)
			return err
		}
	}

	notAfter := notBefore.Add(*validFor)

	// end of ASN.1 time
	endOfTime := time.Date(2049, 12, 31, 2, 59, 59, 0, time.UTC)
	if notAfter.After(endOfTime) {
		notAfter = endOfTime
	}

	template := x509.Certificate{
		SerialNumber: new(big.Int).SetInt64(0),
		Subject: pkix.Name{
			Organization: []string{*organization},
		},
		NotBefore: notBefore,
		NotAfter:  notAfter,

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	hosts := strings.Split(hostnames, ",")
	for _, h := range hosts {
		if ip := net.ParseIP(h); ip != nil {
			template.IPAddresses = append(template.IPAddresses, ip)
		} else {
			template.DNSNames = append(template.DNSNames, h)
		}
	}

	if *isCA {
		template.IsCA = true
		template.KeyUsage |= x509.KeyUsageCertSign
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		log.Fatalf("Failed to create certificate: %s", err)
	}

	certOut, err := os.Create(certFilename)
	if err != nil {
		log.Fatalf("failed to open cert.pem for writing: %s", err)
	}
	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	certOut.Close()
	log.Printf("written %s\n", certFilename)

	keyOut, err := os.OpenFile(keyFilename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Print("failed to open key.pem for writing:", err)
		return err
	}
	pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})
	keyOut.Close()
	log.Printf("written %s\n", keyFilename)
	// We got this for so no errors
	return nil
}

func main() {
	flag.Parse()

	profile, _ := LoadProfile(*cli_docroot, *cli_host, *cli_port, *cli_use_tls, *cli_cert, *cli_key, *cli_otto, *cli_otto_path)
	if *cli_keygen == true {
		err := keygen(profile)
		if err != nil {
			log.Fatalf("%s\n", err)
		}
		os.Exit(0)
	}

	log.Printf("\n\n"+
		"          TLS: %t\n"+
		"         Cert: %s\n"+
		"          Key: %s\n"+
		"      Docroot: %s\n"+
		"         Host: %s\n"+
		"         Port: %s\n"+
		"       Run as: %s\n\n"+
		" Otto enabled: %t\n"+
		"         Path: %s\n"+
		"\n\n",
		profile.Use_TLS,
		profile.Cert,
		profile.Key,
		profile.Docroot,
		profile.Hostname,
		profile.Port,
		profile.Username,
		profile.Otto,
		profile.Otto_Path)
	err := Webserver(profile)
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(0)
}
