//
// Package keygen is based on the keygen code in the TLS package adapted
// to fit the demands of this ws.go project.
//
package keygen

import (
	"../prompt"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"os"
	"path"
	"strings"
	"time"
)

// Keygen will generate the contents of a cert and key PEM file.
func Keygen(basedir string, certPem string, keyPem string) (string, string, error) {
	var (
		certFilename = certPem
		keyFilename  = keyPem
		hostnames    string
		sslPath      string
		rsaBits      int
		OK           = false
	)

	hostnames = os.Getenv("HOSTNAME")
	if hostnames == "" {
		hostnames = "localhost"
	}
	certFilename = "cert.pem"
	keyFilename = "key.pem"
	sslPath = path.Join(basedir)
	for OK == false {
		sslPath = prompt.Question(fmt.Sprintf("Use %s for cert and key? (enter accepts the default) ",
			sslPath), sslPath)
		certFilename = prompt.Question("Certificate filename? (default is cert.pem) ", "cert.pem")
		keyFilename = prompt.Question("Key filename? (default key.pem) ", "key.pem")
		hostnames = prompt.Question(fmt.Sprintf("SSL certs for %s? (enter accepts default, use comma to separate hostnames) ", hostnames), hostnames)
		if hostnames == "" {
			hostnames = "localhost"
		}
		fmt.Printf("\n\n"+
			"    Cert: %s\n"+
			"     Key: %s\n"+
			" Host(s): %s\n"+
			"\n\n",
			path.Join(sslPath, certFilename),
			path.Join(sslPath, keyFilename),
			hostnames)
		OK = prompt.YesNo("Is this correct?")
	}

	// FIXME see if directory exists first
	if sslPath != "" {
		fmt.Printf("Creating %s\n", sslPath)
		err := os.MkdirAll(sslPath, 0770)
		if err != nil {
			return "", "", err
		}
	}

	rsaBits = 2048
	fmt.Println("Generating 2048 bit key")
	priv, err := rsa.GenerateKey(rand.Reader, rsaBits)
	if err != nil {
		return "", "", err
	}
	notBefore := time.Now()
	notAfter := time.Date(2049, 12, 31, 2, 59, 59, 0, time.UTC)

	fmt.Println("Setting up cerificates")
	template := x509.Certificate{
		SerialNumber: new(big.Int).SetInt64(0),
		Subject: pkix.Name{
			Organization: []string{"Acme Co"},
		},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	fmt.Println("Checking IP addresses")
	hosts := strings.Split(hostnames, ",")
	for _, h := range hosts {
		if ip := net.ParseIP(h); ip != nil {
			template.IPAddresses = append(template.IPAddresses, ip)
		} else {
			template.DNSNames = append(template.DNSNames, h)
		}
	}

	template.IsCA = false
	fmt.Println("Generating x509 certs from template")
	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		return "", "", err
	}

	certFilename = path.Join(sslPath, certFilename)
	keyFilename = path.Join(sslPath, keyFilename)
	fmt.Printf("Creating %s\n", certFilename)

	certOut, err := os.Create(certFilename)
	if err != nil {
		return "", "", err
	}
	fmt.Println("Encode as pem")
	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	certOut.Close()
	fmt.Printf("Wrote %s\n", certFilename)

	fmt.Printf("Creating %s\n", keyFilename)
	keyOut, err := os.OpenFile(keyFilename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return "", "", err
	}
	pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})
	keyOut.Close()
	fmt.Printf("Wrote %s\n", keyFilename)
	// We got this far so no errors
	return certFilename, keyFilename, nil
}
