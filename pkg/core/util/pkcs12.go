package util

import (
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"

	"golang.org/x/crypto/pkcs12"
)

func ToPEM(fileContent string, password string) ([]byte, []byte, error) {
	decoded, err := base64.StdEncoding.DecodeString(fileContent)
	if err != nil {
		decoded = []byte(fileContent)
	}

	blocks, err := pkcs12.ToPEM(decoded, password)
	if err != nil {
		// TODO filename in error
		return nil, nil, fmt.Errorf("failure reading pkcs12 certificate: %v", err)
	}

	// TODO feels like there must be a more straightforward way :)
	pemData := []byte{}
	for _, b := range blocks {
		pemData = append(pemData, pem.EncodeToMemory(b)...)
	}

	cert, err := tls.X509KeyPair(pemData, pemData)
	if err != nil {
		return nil, nil, fmt.Errorf("failure getting certificate: %v", err)
	}

	certs := []byte{}
	for _, certPem := range cert.Certificate {
		certs = append(certs, pem.EncodeToMemory(
			&pem.Block{
				Type:  "CERTIFICATE",
				Bytes: certPem,
			},
		)...)
	}

	key := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(cert.PrivateKey.(*rsa.PrivateKey)), // TODO what if not PKCS1?
		},
	)

	return key, certs, nil
}
