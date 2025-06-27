package appschema

import "crypto/rsa"

type CertificateKeys struct {
	PublicKeyPem   *rsa.PublicKey
	PrivateKey     *rsa.PrivateKey
}