/**
 * keygen.go - based on the keygen code in the TLS package adapted
 * to fit the demands of this ws.go project.
 */
package main

import (
    "../app"
    "os"
    "crypto/rand"
    "crypto/rsa"
    "crypto/x509"
    "crypto/x509/pkix"
    "encoding/pem"
    "math/big"
    "path"
    "strings"
    "bufio"
    "fmt"
)

func Keygen(profile *app.Profile) error {
    reader := bufio.NewReader(os.Stdin)

	home := path.Join(os.Getenv("HOME"), "etc/ws")
    fmt.Printf("Write cert/key to %s? (enter accepts the default) ", home)
    line, _ := reader.ReadString('\n')
    if line != "" {
        home = line
        //FIXME: need to do a mkdir -p to ensure the path exists
    }
    certFilename := path.Join(home, "cert.pem")
    fmt.Printf("Certificate filename is %s? (enter accepts default)", certFilename)
    line, _ := reader.ReadString('\n')
    if line != "" {
        certFilename = path.Join(home, line)
    }

	keyFilename := profile.Key
    fmt.Printf("Key filename is %s? (enter accepts default)", keyFilename)
    line, _ := reader.ReadString('\n')
	if line != "" {
		keyFilename = path.Join(home, line)
	}

	hostnames := "localhost;" + os.Getenv("Hostname")
    fmt.Printf("Hostnames for SSL certs %s? (enter accepts default)", hostnames)
    line, _ := reader.ReadString('\n')
	if line != "" {
        hostnames = line
	}

    organization := "Acme Co."
    fmt.Printf("Organization %s? (enter accepts default)?", organization)
    line, _ := reader.ReadString('\n')
    if line != "" {
        organiation = line
    }

    rsaBits := "2048"

	log.Printf("\n\n"+
		" Cert: %s\n"+
		"  Key: %s\n"+
		" Host: %s\n"+
		" Organization: %s\n"+
        // default to 2048, otherwise numeric
        " RSA Bits: %s\n"+
        // Parsable dates
        " Valid from: %s\n"+
        // isCA true/false
        " isCA: %s\n"+
		"\n\n",
		certFilename,
		keyFilename,
		hostnames,
		organization,
        rsaBits,
        validFrom,
        isCA)

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

