/**
 * keygen.go - based on the keygen code in the TLS package adapted
 * to fit the demands of this ws.go project.
 */
package keygen

import (
    "../prompt"
    "os"
    "path"
    "crypto/rand"
    "crypto/rsa"
    "crypto/x509"
    "crypto/x509/pkix"
    "encoding/pem"
    "math/big"
    "time"
    "net"
    "strings"
    "fmt"
)

func Keygen(basedir string, cert_pem string, key_pem string) (string, string, error) {
    var (
        cert_filename = cert_pem
        key_filename = key_pem
        hostnames string
        ssl_path string
        organization string
        rsaBits int
        OK = false
    )

	hostnames = os.Getenv("HOSTNAME")
    ssl_path = path.Join(basedir)
    for OK == false {
        ssl_path = prompt.PromptString(fmt.Sprintf("Use %s for cert and key? (enter accepts the default) ", 
            ssl_path), ssl_path)
        cert_filename = prompt.PromptString("Use cert.pem for certificate file?", "cert.pem")
        key_filename = prompt.PromptString("Use key.pem for key file?", "key.pem")
        hostnames = prompt.PromptString(fmt.Sprintf("SSL certs for %s? (enter accepts default, use comma to separate hostnames)", hostnames), hostnames)
        if hostnames == "" {
            hostnames = "localhost"
        }
	    fmt.Printf("\n\n"+
		    " Cert: %s\n"+
		    "  Key: %s\n"+
		    " Host: %s\n"+
		    "\n\n",
            path.Join(ssl_path, cert_filename),
		    path.Join(ssl_path, key_filename),
		    hostnames)
        OK = prompt.YesNo("Is this correct?")
    }

    // FIXME see if directory exists first
    if ssl_path != "" {
        fmt.Printf("Creating %s\n", ssl_path)
        err := os.MkdirAll(ssl_path, 0770)
        if err != nil {
            return "", "", err 
        }
    }

    organization = "Acme Co."
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
			Organization: []string{organization},
		},
		NotBefore: notBefore,
		NotAfter:  notAfter,
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

	template.IsCA = true
	template.KeyUsage |= x509.KeyUsageCertSign

    fmt.Println("Generating x509 certs from template")
	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
        return "", "", err
	}

    certFilename := path.Join(ssl_path, cert_filename)
    keyFilename := path.Join(ssl_path, key_filename)
    fmt.Printf("Creating %s", certFilename)
	certOut, err := os.Create(certFilename)
	if err != nil {
        return "", "", err
	}
    fmt.Println("Encode as pem")
	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	certOut.Close()
	fmt.Printf("Wrote %s\n", certFilename)

    fmt.Printf("Creating %s", keyFilename)
	keyOut, err := os.OpenFile(keyFilename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return "", "", err
	}
	pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})
	keyOut.Close()
	fmt.Printf("Wrote %s\n", keyFilename)
	// We got this for so no errors
	return certFilename, keyFilename, nil
}

