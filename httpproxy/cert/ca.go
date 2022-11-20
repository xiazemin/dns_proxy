package cert

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"math/big"
	"net"
	"os"
	"time"
)

func GetCrt() (crt, key string) {
	max := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, _ := rand.Int(rand.Reader, max)
	subject := pkix.Name{
		Organization:       []string{"燕子李三"},
		OrganizationalUnit: []string{"Books"},
		CommonName:         "GO Web",
	}

	cwd, _ := os.Getwd()
	path := cwd + "/httpproxy/cert/"
	rootcertificate, rootPrivateKey, err := parseCerts(path+"server.pem", path+"server.key")
	fmt.Println(rootPrivateKey, err)

	rootTemplate := x509.Certificate{
		Version:      1,
		SerialNumber: serialNumber,
		Subject:      subject,
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(365 * time.Hour),
		KeyUsage:     x509.KeyUsageKeyEncipherment | x509.KeyUsageCertSign | x509.KeyUsageCRLSign | x509.KeyUsageDataEncipherment,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		IPAddresses:  []net.IP{net.ParseIP("127.0.0.1")},
		Extensions:   []pkix.Extension{},
		SubjectKeyId: []byte{1, 2, 3},
		DNSNames:     []string{"localhost", "*.xzm.com"},
	}
	pk, _ := rsa.GenerateKey(rand.Reader, 2048)
	derBytes, _ := x509.CreateCertificate(rand.Reader, &rootTemplate, rootcertificate, &pk.PublicKey, rootPrivateKey)

	certOut, _ := os.Create(path + "/cert.pem")
	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})

	certOut.Close()

	keyOut, _ := os.Create(path + "/key.pem")
	pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(pk)})
	keyOut.Close()

	// http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
	// 	io.WriteString(w, "Hello, TLS!\n")
	// })

	// // One can use generate_cert.go in crypto/tls to generate cert.pem and key.pem.
	// log.Printf("About to listen on 8443. Go to https://127.0.0.1:8443/")
	// err = http.ListenAndServeTLS(":8443", path+"/cert.pem", path+"/key.pem", nil)
	// log.Fatal(err)
	return path + "/cert.pem", path + "/key.pem"
}

func parseCerts(ROOTCertificate, ROOTPrivateKey string) (*x509.Certificate, *rsa.PrivateKey, error) {
	var rootPrivateKey *rsa.PrivateKey
	var rootcertificate *x509.Certificate

	buf, err := ioutil.ReadFile(ROOTCertificate)
	if err != nil {
		return nil, nil, err
	}
	p := &pem.Block{}
	p, buf = pem.Decode(buf)
	rootcertificate, err = x509.ParseCertificate(p.Bytes)
	if err != nil {
		return nil, nil, err
	}
	buf, err = ioutil.ReadFile(ROOTPrivateKey)
	if err != nil {
		return nil, nil, err
	}
	p, buf = pem.Decode(buf)
	rootPrivateKey, err = x509.ParsePKCS1PrivateKey(p.Bytes)
	if err != nil {
		return nil, nil, err
	}
	return rootcertificate, rootPrivateKey, nil
}
