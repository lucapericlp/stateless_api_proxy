package keys

import (
	"bufio"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"golang.org/x/crypto/ssh"
	"os"
)

type Keys struct {
	PubKey  *rsa.PublicKey
	PrivKey *rsa.PrivateKey
	Signer  *ssh.Signer
}

//TO DO: derive signer using std lib
func LoadKeys() *Keys {
	privateKeyFile, err := os.Open(os.Getenv("MAGICTOKEN_PRIVATE_KEY"))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	//read in private key, also get public key
	pemfileinfo, _ := privateKeyFile.Stat()
	pembytes := make([]byte, pemfileinfo.Size()) //int64

	buffer := bufio.NewReader(privateKeyFile)
	_, err = buffer.Read(pembytes)

	block, _ := pem.Decode(pembytes)
	parseResult, _ := x509.ParsePKCS8PrivateKey(block.Bytes)
	signer, _ := ssh.NewSignerFromKey(parseResult)
	privKey := parseResult.(*rsa.PrivateKey)
	pubKey := (privKey.Public()).(*rsa.PublicKey)
	return &Keys{
		PubKey:  pubKey,
		PrivKey: privKey,
		Signer:  &signer,
	}
}
