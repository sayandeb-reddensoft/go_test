package utils

import (
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/nelsonin-research-org/clenz-auth/globals"
	"github.com/nelsonin-research-org/clenz-auth/models/appschema"
)

func LoadCertificateKeys() error {
	publicKey, err := os.ReadFile(os.Getenv("PUBLIC_PEM_PATH"))
	if err != nil {
		return err
	}

	privateKey, err := os.ReadFile(os.Getenv("PRIVATE_PEM_PATH"))
	if err != nil {
		return err
	}

	PrivatePem, err := jwt.ParseRSAPrivateKeyFromPEM(privateKey)
	if err != nil {
		return err
	}

	PublicPem, err := jwt.ParseRSAPublicKeyFromPEM(publicKey)
	if err != nil {
		return err
	}

	globals.AppKeys = appschema.CertificateKeys{
		PublicKeyPem: PublicPem,
		PrivateKey:   PrivatePem,
	}

	return nil
}
