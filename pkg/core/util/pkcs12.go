package util

import (
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"

	pkcs12 "software.sslmate.com/src/go-pkcs12"
)

func ToPEM(fileContent string, password string, chainDelimiter string) ([]byte, []byte, error) {
	decoded, err := base64.StdEncoding.DecodeString(fileContent)
	if err != nil {
		decoded = []byte(fileContent)
	}

	privateKey, certificate, chain, err := pkcs12.DecodeChain(decoded, password)
	if err != nil {
		return nil, nil, fmt.Errorf("failure decoding pkcs12: %v", err)
	}

	x509Certs, err := sorted(certificate, chain)
	if err != nil {
		return nil, nil, err
	}

	certs := []byte{}
	for i, certPem := range x509Certs {
		if i > 0 && chainDelimiter != "" {
			certs = append(certs, []byte(chainDelimiter)...)
		}
		certs = append(certs, pem.EncodeToMemory(
			&pem.Block{
				Type:  "CERTIFICATE",
				Bytes: certPem.Raw,
			},
		)...)
	}

	privateKeyBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to marshal pkcs8 private key: %v", err)
	}

	key := pem.EncodeToMemory(
		&pem.Block{
			Type:  "PRIVATE KEY",
			Bytes: privateKeyBytes,
		},
	)

	return key, certs, nil
}

func sorted(certificate *x509.Certificate, caCerts []*x509.Certificate) ([]*x509.Certificate, error) {
	res := []*x509.Certificate{certificate}
	if len(caCerts) == 0 {
		return res, nil
	} else {
		res = append(res, recurse(certificate, caCerts)...)
	}
	if len(res) != len(caCerts)+1 {
		return nil, fmt.Errorf("Unexpected CA chain")
	}
	return res, nil
}

func recurse(certificate *x509.Certificate, caCerts []*x509.Certificate) []*x509.Certificate {
	for _, caCert := range caCerts {
		if certificate.Issuer.CommonName == caCert.Subject.CommonName && caCert.Issuer.CommonName == caCert.Subject.CommonName {
			// root CA
			return []*x509.Certificate{caCert}
		}
		if certificate.Issuer.CommonName == caCert.Subject.CommonName {
			res := []*x509.Certificate{caCert}
			return append(res, recurse(caCert, caCerts)...)
		}
	}
	return []*x509.Certificate{}
}
