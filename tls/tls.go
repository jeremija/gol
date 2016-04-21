package tls

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
)

func ReadCaCerts(caCertsFilename string) *x509.CertPool {
	certs, err := ioutil.ReadFile(caCertsFilename)

	if err != nil {
		message := fmt.Sprintf(
			"Error reading CA certs file: %s - %s",
			caCertsFilename,
			err.Error(),
		)
		panic(message)
	}

	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(certs)
	return pool
}

func GetTlsConfig(caCertsFilename string) *tls.Config {
	if caCertsFilename == "" {
		return nil
	}
	return &tls.Config{
		RootCAs: ReadCaCerts(caCertsFilename),
	}
}
