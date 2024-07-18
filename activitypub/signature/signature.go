package signature

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-fed/httpsig"
)

func VerifySignature(publicKeyPem string, r *http.Request) bool {
	v, err := httpsig.NewVerifier(r)
	if err != nil {
		slog.Error("Failed to create verifier", "Error", err)
		return false
	}

	block, _ := pem.Decode([]byte(publicKeyPem))
	if block == nil {
		slog.Error("Failed to decode public key")
		return false
	}

	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)

	if err != nil {
		if strings.Contains(err.Error(), "ParsePKCS1PublicKey") {
			publicKey, err = x509.ParsePKCS1PublicKey(block.Bytes)
			if err != nil {
				slog.Error("Failed to load public key", "Error", err)
				return false
			}
		}
	}

	signature := r.Header.Get("Signature")
	if signature == "" {
		slog.Error("Signature header not found")
		return false
	}

	signatureSep := strings.Split(signature, ",")
	if len(signatureSep) < 2 {
		slog.Error("Invalid signature header", "Signature", signature)
		return false
	}
	signatureMap := make(map[string]string)
	for _, s := range signatureSep {
		pos := strings.Index(s, "=")
		if pos == -1 {
			slog.Error("Invalid signature header", "s", s)
			return false
		}
		key := s[:pos]
		value := s[pos+2 : len(s)-1]
		signatureMap[key] = value
	}

	if signatureMap["algorithm"] == "" {
		slog.Error("Algorithm not found in signature header")
		return false
	}

	err = v.Verify(publicKey, httpsig.Algorithm(signatureMap["algorithm"]))
	if err != nil {
		slog.Error("Failed to verify signature", "Error", err)
		return false
	}

	return true

}

func Signature(privateKeyPem string, keyId string, r *http.Request) error {
	algos := []httpsig.Algorithm{
		httpsig.RSA_SHA256,
	}

	headersToSign := []string{httpsig.RequestTarget} // "date",

	slog.Info("Request headers", "Headers", r.Header, "path", r.URL.Path)

	if r.Header.Get("Date") != "" {
		headersToSign = append(headersToSign, "date")
	}
	if r.Header.Get("Host") != "" {
		headersToSign = append(headersToSign, "host")
	}

	headersToSign = append(headersToSign, "digest")

	s, algo, err := httpsig.NewSigner(
		algos,
		httpsig.DigestSha256,
		headersToSign,
		httpsig.Signature,
		10)
	if err != nil {
		slog.Error("Failed to create signer", "Error", err)
		return err
	}
	slog.Info("Signature algorithm", "Algorithm", algo)

	block, _ := pem.Decode([]byte(privateKeyPem))
	if block == nil {
		slog.Error("Failed to decode private key", "privateKeyPem", privateKeyPem)
		return err
	}

	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		slog.Error("Failed to load private key", "Error", err)
		return err
	}
	data := []byte{}
	if r.Body != nil {
		data, err = io.ReadAll(r.Body)
		if err != nil {
			slog.Error("Failed to read request body", "Error", err)
			return err
		}
	}
	err = s.SignRequest(privateKey, keyId, r, data)
	if err != nil {
		slog.Error("Failed to sign request", "Error", err)
		return err
	}

	// replace hs2019 to rsa-sha256
	signature := r.Header.Get("Signature")
	signature = strings.Replace(signature, "hs2019", "rsa-sha256", 1)
	r.Header.Set("Signature", signature)
	slog.Info("Signature header", "Signature", signature)

	return nil
}

func Pubout(privateKeyPem []byte) []byte {

	// parse pem file
	block, _ := pem.Decode(privateKeyPem)
	if block == nil {
		slog.Error("Failed to decode public key")
		return nil
	}

	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		slog.Error("Failed to load public key", "Error", err, "pem", privateKeyPem)
		return nil
	}

	// encode to pem
	keyBytes := x509.MarshalPKCS1PublicKey(&key.(*rsa.PrivateKey).PublicKey)

	return pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: keyBytes,
	})
}

func GeneratePrivateKey() string {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		slog.Error("Failed to generate private key", "Error", err)
		return ""
	}

	keyBytes, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		slog.Error("Failed to marshal private key", "Error", err)
		return ""
	}

	return string(pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: keyBytes,
	}))
}

// func VerifyRequest(r *http.Request) error {
// 	sig := r.Header.Get("Signature")
// 	if sig == "" {
// 		slog.Warn("Signature header not found")
// 		return fmt.Errorf("Signature header not found")
// 	}

// 	// get keyid
// 	keyId := ""
// 	sigSep := strings.Split(sig, ",")
// 	for _, s := range sigSep {
// 		if strings.Contains(s, "keyId") {
// 			keyId = strings.Split(s, "=")[1]
// 			keyId = strings.Trim(keyId, "\"")
// 			break
// 		}
// 	}
// 	if keyId == "" {
// 		slog.Warn("keyId not found in signature header")
// 		return fmt.Errorf("keyId not found in signature header")
// 	}

// 	pem, err := GetKeyById(keyId)
// 	if err != nil {
// 		slog.Warn("Failed to get public key", "Error", err, "keyId", keyId)
// 		return err
// 	}

// 	if !VerifySignature(pem, r) {
// 		slog.Warn("Failed to verify signature")
// 		return fmt.Errorf("failed to verify signature")
// 	}
// 	return nil

// }
