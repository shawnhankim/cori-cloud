package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/pem"
	"os"
)

var oidEmailAddress = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 9, 1}

func main() {
	keyBytes, _ := rsa.GenerateKey(rand.Reader, 1024)

	emailAddress := "test@example.com"
	subj := pkix.Name{
		CommonName:         "example.com",
		Country:            []string{"AU"},
		Province:           []string{"Some-State"},
		Locality:           []string{"MyCity"},
		Organization:       []string{"Company Ltd"},
		OrganizationalUnit: []string{"IT"},
		ExtraNames: []pkix.AttributeTypeAndValue{
			{
				Type: oidEmailAddress,
				Value: asn1.RawValue{
					Tag:   asn1.TagIA5String,
					Bytes: []byte(emailAddress),
				},
			},
		},
	}

	template := x509.CertificateRequest{
		Subject:            subj,
		SignatureAlgorithm: x509.SHA256WithRSA,
	}

	csrBytes, _ := x509.CreateCertificateRequest(rand.Reader, &template, keyBytes)
	pem.Encode(os.Stdout, &pem.Block{Type: "CERTIFICATE REQUEST", Bytes: csrBytes})

}
