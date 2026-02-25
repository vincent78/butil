package tls

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"log"
	"math/big"
	"os"
	"time"

	"io/ioutil"
)

/*
sample golang minimalapp that
with a generated ca public key and file  (CA_crt.pem, CA_key.pem  (no passphrase))
will.
   1. create a kew rsa keypair and use the CA to sign that.
   2. create a private key, then a certificate signing request (CSR), then have the CA sign that
 generate a CA:
	https://github.com/salrashid123/squid_proxy#generating-new-ca
*/

const ()

var ()

func main() {

	log.Printf("load ca crt, key")

	certPEMBytes, err := ioutil.ReadFile("CA_crt.pem")
	if err != nil {
		log.Fatalf("err %v", err)
	}
	block, _ := pem.Decode(certPEMBytes)
	if block == nil {
		log.Fatalf("err %v", err)
	}
	ca, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		log.Fatalf("err %v", err)
	}

	log.Printf(ca.SerialNumber.String())

	keyPEMBytes, err := ioutil.ReadFile("CA_key.pem")
	if err != nil {
		log.Fatalf("err %v", err)
	}
	privPem, _ := pem.Decode(keyPEMBytes)
	parsedKey, err := x509.ParsePKCS1PrivateKey(privPem.Bytes)
	if err != nil {
		log.Fatalf("err %v", err)
	}

	// ************************

	// create keypair and sign with cakey
	var notBefore time.Time
	notBefore = time.Now()

	notAfter := notBefore.Add(time.Hour * 24 * 365)

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		log.Fatalf("Failed to generate serial number: %s", err)
	}

	cc := &x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization:       []string{"Acme Co"},
			OrganizationalUnit: []string{"Enterprise"},
			Locality:           []string{"Mountain View"},
			Province:           []string{"California"},
			Country:            []string{"US"},
			CommonName:         "server.domain.com",
		},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		DNSNames:              []string{"server.domain.com"},
		KeyUsage:              x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
		IsCA:                  false,
	}

	priv, _ := rsa.GenerateKey(rand.Reader, 2048)
	pub := &priv.PublicKey

	cert_b, err := x509.CreateCertificate(rand.Reader, cc, ca, pub, parsedKey)
	if err != nil {
		log.Fatalf("err %v", err)
	}
	certOut, err := os.Create("server.crt")
	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: cert_b})
	certOut.Close()
	log.Print("written cert.pem\n")

	// Private key
	keyOut, err := os.OpenFile("server.key", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})
	keyOut.Close()
	log.Print("written key.pem\n")

	/// *****************************  gen priv key and create csr

	server2_priv, _ := rsa.GenerateKey(rand.Reader, 2048)
	filename_csr := "server2.csr"
	var csrtemplate = x509.CertificateRequest{
		Subject: pkix.Name{
			Organization:       []string{"Acme Co"},
			OrganizationalUnit: []string{"Enterprise"},
			Locality:           []string{"Mountain View"},
			Province:           []string{"California"},
			Country:            []string{"US"},
			CommonName:         "server2.domain.com",
		},
		DNSNames: []string{"server2.domain.com"},
	}

	csrBytes, err := x509.CreateCertificateRequest(rand.Reader, &csrtemplate, server2_priv)
	if err != nil {
		log.Fatalf("Failed to create CSR: %s", err)
	}
	certOut, err = os.Create(filename_csr)
	if err != nil {
		log.Fatalf("Failed to open %s for writing: %s", filename_csr, err)
	}
	if err := pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE REQUEST", Bytes: csrBytes}); err != nil {
		log.Fatalf("Failed to write data to %s: %s", filename_csr, err)
	}
	if err := certOut.Close(); err != nil {
		log.Fatalf("Error closing %s  %s", filename_csr, err)
	}
	log.Printf("wrote %s\n", filename_csr)

	//  sign csr

	var csr_notBefore time.Time
	csr_notBefore = time.Now()

	csr_notAfter := csr_notBefore.Add(time.Hour * 24 * 365)

	csr_serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		log.Fatalf("Failed to generate serial number: %s", err)
	}

	ccr := &x509.Certificate{
		SerialNumber: csr_serialNumber,
		Subject: pkix.Name{
			Organization:       []string{"Acme Co"},
			OrganizationalUnit: []string{"Enterprise"},
			Locality:           []string{"Mountain View"},
			Province:           []string{"California"},
			Country:            []string{"US"},
			CommonName:         csrtemplate.Subject.CommonName,
		},
		NotBefore:             csr_notBefore,
		NotAfter:              csr_notAfter,
		DNSNames:              csrtemplate.DNSNames,
		KeyUsage:              x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
		IsCA:                  false,
	}

	csrPEMBytes, err := ioutil.ReadFile("server2.csr")
	if err != nil {
		log.Fatalf("err %v", err)
	}
	csrblock, _ := pem.Decode(csrPEMBytes)
	if block == nil {
		log.Fatalf("err %v", err)
	}

	csr, err := x509.ParseCertificateRequest(csrblock.Bytes)
	if err != nil {
		log.Fatalf("Error closing %s  %s", filename_csr, err)
	}

	cert_c, err := x509.CreateCertificate(rand.Reader, ccr, ca, csr.PublicKey, parsedKey)
	if err != nil {
		log.Fatalf("err %v", err)
	}

	newcert, err := x509.ParseCertificate(cert_c)
	if err != nil {
		log.Fatalf("err %v", err)
	}
	log.Println(newcert.Subject.CommonName)

}
