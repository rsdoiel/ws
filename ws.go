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
    "./fsengine"
	"./ottoengine"
	"./wslog"
    "./app"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"flag"
	"fmt"
	"log"
	"math/big"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

var REVISION = "v0.0.0-alpha"

// variables for keygen
var (
	cli_keygen   = flag.Bool("keygen", false, "Generate TLS ceriticates and keys")
	cli_ssl_host = flag.String("keygen-ssl-host", "", "Comma-separated hostnames and IPs to generate a certificate for")
	validFrom    = flag.String("keygen-start-date", "", "Creation date formatted as Jan 1 15:04:05 2011")
	validFor     = flag.Duration("keygen-duration", 365*24*time.Hour, "Duration that certificate is valid for")
	organization = flag.String("keygen-organization", "Acme Co.", "Organization used to sign certificate")
	isCA         = flag.Bool("keygen-ca", false, "whether this cert should be its own Certificate Authority")
	rsaBits      = flag.Int("keygen-rsa-bits", 2048, "Size of RSA key to generate")
)

// command line parameters that override environment variables
var (
	cli_use_tls   *bool
	cli_docroot   *string
	cli_host      *string
	cli_port      *string
	cli_cert      *string
	cli_key       *string
	cli_otto      *bool
	cli_otto_path *string
	cli_version   *bool
)

var Usage = func() {
	flag.PrintDefaults()
}

func request_log(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wslog.LogRequest(r.Method, r.URL, r.RemoteAddr, r.Proto, r.Referer(), r.UserAgent())
		handler.ServeHTTP(w, r)
	})
}

func Webserver(profile *app.Profile) error {
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
	http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
        // hande off this request/response pair to the fsengine
        fsengine.Engine(profile, w, r)
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

func keygen(profile *app.Profile) error {
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

func init() {
	// command line parameters that override environment variables, many have short forms too.
	shortform := " (short form)"

	// No short form
	cli_use_tls = flag.Bool("tls", false, "Turn on TLS (https) support with true, off with false (default is false)")
	cli_key = flag.String("key", "", "path to your SSL key pem file.")
	cli_cert = flag.String("cert", "", "path to your SSL cert pem file.")

	// These have short forms too
	msg := "document root"
	cli_docroot = flag.String("docroot", "", msg)
	flag.StringVar(cli_docroot, "D", "", msg+shortform)

	msg = "hostname for webserver"
	cli_host = flag.String("host", "", msg)
	flag.StringVar(cli_host, "H", "", msg+shortform)

	msg = "Port number to listen on"
	cli_port = flag.String("port", "", msg)
	flag.StringVar(cli_port, "P", "", msg+shortform)

	msg = "turn on ottoengine, defaults to false"
	cli_otto = flag.Bool("otto", false, msg)
	flag.BoolVar(cli_otto, "o", false, msg+shortform)

	msg = "directory containingo your ottoengine JavaScript files"
	cli_otto_path = flag.String("otto-path", "", msg)
	flag.StringVar(cli_otto_path, "op", "", msg+shortform)

	msg = "Display the version number"
	cli_version = flag.Bool("version", false, msg)
	flag.BoolVar(cli_version, "v", false, msg+shortform)
}

func main() {
	flag.Parse()

	if *cli_version == true {
		fmt.Println(REVISION)
		os.Exit(0)
	}

	profile, _ := app.LoadProfile(*cli_docroot, *cli_host, *cli_port, *cli_use_tls, *cli_cert, *cli_key, *cli_otto, *cli_otto_path)
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
