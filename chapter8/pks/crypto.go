package pks

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"os"
)

// The default key size 3072 bits; this is estimated by the NIST to be
// equivalent to a 128-bit AES key.
const KeySize = 3072

var defaultLabel = []byte{}

// GenerateKey generates an RSA key with the default key size.
func GenerateKey() (key *rsa.PrivateKey, err error) {
	return rsa.GenerateKey(rand.Reader, KeySize)
}

// Sign generates a digital signature of the message passed in.
func Sign(prv *rsa.PrivateKey, m []byte) (sig []byte, err error) {
	h := sha256.New()
	h.Write(m)
	d := h.Sum(nil)
	sig, err = rsa.SignPSS(rand.Reader, prv, crypto.SHA256, d, nil)
	return
}

// Verify verifies an RSA digital signature for the given public key.
func Verify(pub *rsa.PublicKey, m, sig []byte) (err error) {
	h := sha256.New()
	h.Write(m)
	d := h.Sum(nil)
	return rsa.VerifyPSS(pub, crypto.SHA256, d, sig, nil)
}

// ExportPrivateKey writes the RSA private key to a file in DER format.
func ExportPrivateKey(prv *rsa.PrivateKey, filename string) (err error) {
	cert := x509.MarshalPKCS1PrivateKey(prv)
	err = ioutil.WriteFile(filename, cert, 0600)
	return
}

// ExportPrivatePEM writes the RSA private key to a file in PEM format.
func ExportPrivatePEM(prv *rsa.PrivateKey, filename string) (err error) {
	cert := x509.MarshalPKCS1PrivateKey(prv)
	blk := new(pem.Block)
	blk.Type = "RSA PRIVATE KEY"
	blk.Bytes = cert
	out, err := os.Create(filename)
	if err == nil {
		err = pem.Encode(out, blk)
	}
	return
}

// ExportPublicKey writes the RSA public key to a file in DER format.
func ExportPublicKey(pub *rsa.PublicKey, filename string) (err error) {
	cert, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		return
	}
	err = ioutil.WriteFile(filename, cert, 0644)
	return
}

// ExportPublicPEM writes the public key to a file in PEM format.
func ExportPublicPEM(pub *rsa.PublicKey, filename string) (err error) {
	cert, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		return
	}
	blk := new(pem.Block)
	blk.Type = "RSA PUBLIC KEY"
	blk.Bytes = cert
	out, err := os.Create(filename)
	if err == nil {
		err = pem.Encode(out, blk)
	}
	return
}

// ImportPrivateKey reads an RSA private key in DER format from a file.
func ImportPrivateKey(filename string) (prv *rsa.PrivateKey, err error) {
	cert, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}

	prv, err = x509.ParsePKCS1PrivateKey(cert)
	if err != nil {
		return
	}
	return
}

// ImportPublicKey reads an RSA public key in DER format from a file.
func ImportPublicKey(filename string) (pub *rsa.PublicKey, err error) {
	cert, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}

	in, err := x509.ParsePKIXPublicKey(cert)
	if err != nil {
		return nil, err
	}
	pub = in.(*rsa.PublicKey)
	return
}

// ImportPEM imports an RSA key from a file. It works with both public and
// private keys.
func ImportPEM(filename string) (prv *rsa.PrivateKey, pub *rsa.PublicKey, err error) {
	cert, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}

	for {
		var blk *pem.Block
		blk, cert = pem.Decode(cert)
		if blk == nil {
			break
		}
		switch blk.Type {
		case "RSA PRIVATE KEY":
			prv, err = x509.ParsePKCS1PrivateKey(blk.Bytes)
			return
		case "RSA PUBLIC KEY":
			var in interface{}
			in, err = x509.ParsePKIXPublicKey(blk.Bytes)
			if err != nil {
				return
			}
			pub = in.(*rsa.PublicKey)
			return
		}
		if cert == nil || len(cert) == 0 {
			break
		}
	}
	return
}
