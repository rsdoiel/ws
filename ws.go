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
	"net/http"
	"os"
	"os/user"
	"path"
	"regexp"
	"strconv"
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
	cli_use_tls = flag.Bool("tls", false, "Turn on TLS (https) support with true, off with false (default is false)")
	cli_docroot = flag.String("docroot", "", "Path to the document root")
	cli_host    = flag.String("host", "", "Hostname http(s) server to listen for")
	cli_port    = flag.Int("port", 0, "Port number to listen on")
	cli_cert    = flag.String("cert", "", "Filename to your SSL cert.pem")
	cli_key     = flag.String("key", "", "Filename to your SSL key.pem")
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

func LoadProfile(cli_docroot string, cli_host string, cli_port int, cli_use_tls bool, cli_cert string, cli_key string) (*Profile, error) {
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

	// FIXME: before we return to fail to load on *.pem, check for alternate locations
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
	if len(cli_host) != 0 {
		hostname = cli_host
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
		log.Printf("\n\n"+
			"  Docroot:   %s\n"+
			"     Host:   %s\n"+
			"     Port:   %s\n"+
			"   Run as:   %s\n\n",
			profile.Docroot, profile.Hostname, profile.Port,
			profile.Username)
		log.Println("Starting http://" + net.JoinHostPort(profile.Hostname, profile.Port))

		// Now start up the server and log transactions
		return http.ListenAndServe(net.JoinHostPort(profile.Hostname, profile.Port), Log(http.DefaultServeMux))
	}
	log.Printf("\n\n"+
		"    Cert:   %s\n"+
		"     Key:   %s\n"+
		" Docroot:   %s\n"+
		"    Host:   %s\n"+
		"    Port:   %s\n"+
		"  Run as:   %s\n\n",
		profile.Cert,
		profile.Key,
		profile.Docroot,
		profile.Hostname,
		profile.Port,
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
	log.Printf("Configuration files in %s\n", config_path)
	return config_path, nil
}

func ConfigPathTo(filename string) (string, error) {
	ws_path, err := ConfigPath()
	if err != nil {
		return "", err
	}
	return ws_path + "/" + filename, nil
}

func keygen(profile *Profile) error {
	certFilename := profile.Cert
	if certFilename == "" {
		log.Fatalln("Missing required -cert option")
		os.Exit(1)
	}
	keyFilename := profile.Key
	if keyFilename == "" {
		log.Fatalln("Missing required -key option")
		os.Exit(1)
	}

	hostnames := profile.Hostname
	if *cli_ssl_host != "" {
		hostnames = *cli_ssl_host
	}
	if hostnames == "" {
		log.Fatalf("Missing required -ssl-host parameter")
		os.Exit(1)
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
		return err
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
		return err
	}

	certOut, err := os.Create(certFilename)
	if err != nil {
		log.Fatalf("failed to open cert.pem for writing: %s", err)
		return err
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

	profile, _ := LoadProfile(*cli_docroot, *cli_host, *cli_port, *cli_use_tls, *cli_cert, *cli_key)
	if *cli_keygen == true {
		err := keygen(profile)
		if err != nil {
			log.Fatalf("%s", err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	err := webserver(profile)
	if err != nil {
		log.Fatal(err)
	}
}
