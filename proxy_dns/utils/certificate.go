package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"time"

	"google.golang.org/grpc/credentials"
)

func GenKeyPair() (rawCert, rawKey []byte, err error) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return
	}
	validFor := time.Hour * 24 * 365 * 10 // ten years
	notBefore := time.Now()
	notAfter := notBefore.Add(validFor)
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"GuardLink"},
		},
		NotBefore: notBefore,
		NotAfter:  notAfter,

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}
	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		return
	}

	rawCert = pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	rawKey = pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})

	return
}

func NewServerTLSFromString(serverCrt string, serverKey string, caCrt string) (credentials.TransportCredentials, error) {
	pair, err := tls.X509KeyPair([]byte(serverCrt), []byte(serverKey))
	if err != nil {
		return nil, err
	}

	// 创建一组根证书
	certPool := x509.NewCertPool()
	// 解析证书
	if ok := certPool.AppendCertsFromPEM([]byte(caCrt)); !ok {
		fmt.Println("AppendCertsFromPEM error")
		return nil, errors.New("AppendCertsFromPEM error")
	}

	return credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{pair},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    certPool,
	}), nil
}

func NewClientTLSFromString(serverCrt string, serverKey string, caCrt string, serverName string) (credentials.TransportCredentials, error) {
	//pair, err := tls.X509KeyPair([]byte(serverCrt), []byte(serverKey))
	//if err != nil {
	//	return nil, err
	//}

	// 创建一组根证书
	certPool := x509.NewCertPool()
	// 解析证书
	if ok := certPool.AppendCertsFromPEM([]byte(caCrt)); !ok {
		fmt.Println("AppendCertsFromPEM error")
		return nil, errors.New("AppendCertsFromPEM error")
	}

	return credentials.NewTLS(&tls.Config{
		//Certificates: []tls.Certificate{pair},
		//ClientAuth:   tls.RequireAndVerifyClientCert,
		ServerName: serverName,
		ClientCAs:  certPool,
	}), nil
}
