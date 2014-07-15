/**
 * keygen.go - based on the keygen code in the TLS package adapted
 * to fit the demands of this ws.go project.
 */
package keygen

import (
    "../app"
    "os"
    "path"
    "crypto/rand"
    "crypto/rsa"
    "crypto/x509"
    "crypto/x509/pkix"
    "encoding/pem"
    "math/big"
    "bufio"
    "fmt"
    "log"
    "time"
    "net"
    "strings"
)

func Keygen(profile *app.Profile) error {
    var (
        organization string
        rsaBits int
        line string
    )

    reader := bufio.NewReader(os.Stdin)

	basedir :=  "./etc"
    fmt.Printf("Write cert and key to directory %s? (enter accepts the default) ", basedir)
    line, _ = reader.ReadString('\n')
    line = strings.TrimSpace(line)
    if line != "" {
        basedir = line
    }
    err := os.MkdirAll(basedir, 0770)
    if err != nil {
        log.Fatalf("%s\n", err)
    }

    certFilename := path.Join(basedir, "cert.pem")
    fmt.Printf("Certificate filename is %s? (enter accepts default)", certFilename)
    line, _ = reader.ReadString('\n')
    line = strings.TrimSpace(line)
    if line != "" {
        certFilename = line
    }

	keyFilename := path.Join(basedir, "key.pem")
    fmt.Printf("Key filename is %s? (enter accepts default)", keyFilename)
    line, _ = reader.ReadString('\n')
    line = strings.TrimSpace(line)
	if line != "" {
		keyFilename = line
	}

	hostnames := os.Getenv("HOSTNAME")
    if hostnames == "" {
        hostnames = "localhost"
    }
    fmt.Printf("SSL certs for %s? (enter accepts default, use comma to separate hostnames)", hostnames)
    line, _ = reader.ReadString('\n')
    line = strings.TrimSpace(line)
	if line != "" {
        hostnames = line
	}

	fmt.Printf("\n\n"+
		" Cert: %s\n"+
		"  Key: %s\n"+
		" Host: %s\n"+
		"\n\n",
		certFilename,
		keyFilename,
		hostnames)


    organization = "Acme Co."
    rsaBits = 2048
	priv, err := rsa.GenerateKey(rand.Reader, rsaBits)
	if err != nil {
		log.Fatalf("failed to generate private key: %s", err)
	}
	notBefore := time.Now()
	notAfter := time.Date(2049, 12, 31, 2, 59, 59, 0, time.UTC)

	template := x509.Certificate{
		SerialNumber: new(big.Int).SetInt64(0),
		Subject: pkix.Name{
			Organization: []string{organization},
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

	template.IsCA = true
	template.KeyUsage |= x509.KeyUsageCertSign

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

