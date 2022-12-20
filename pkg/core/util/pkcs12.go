package util

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"

	pkcs12 "software.sslmate.com/src/go-pkcs12"
)

func ToPEM(fileContent string, password string) ([]byte, []byte, error) {
	decoded, err := base64.StdEncoding.DecodeString(fileContent)
	if err != nil {
		decoded = []byte(fileContent)
	}

	privateKey, certificate, chain, err := pkcs12.DecodeChain(decoded, password)
	if err != nil {
		return nil, nil, fmt.Errorf("failure decoding pkcs12: %v", err)
	}

	x509Certs := []*x509.Certificate{certificate}
	x509Certs = append(x509Certs, chain...)

	certs := []byte{}
	for _, certPem := range x509Certs {
		certs = append(certs, pem.EncodeToMemory(
			&pem.Block{
				Type:  "CERTIFICATE",
				Bytes: certPem.Raw,
			},
		)...)
	}

	var privateKeyBytes []byte
	switch typed_privateKey := privateKey.(type) {
	case *rsa.PrivateKey:
		privateKeyBytes = x509.MarshalPKCS1PrivateKey(typed_privateKey)
	case *ecdsa.PrivateKey:
		privateKeyBytes, err = x509.MarshalECPrivateKey(typed_privateKey)
		if err != nil {
			return nil, nil, fmt.Errorf("unable to marshal ec private key: %v", err)
		}
	case string:
		privateKeyBytes, err = x509.MarshalPKCS8PrivateKey(typed_privateKey)
		if err != nil {
			return nil, nil, fmt.Errorf("unable to marshal pkcs8 private key: %v", err)
		}
	default:
		return nil, nil, fmt.Errorf("unhandled private key type: %T", privateKey)
	}

	key := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: privateKeyBytes,
		},
	)

	return key, certs, nil
}
