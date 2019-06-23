package magictoken

import (
	"bufio"
	"crypto"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

type MagicToken struct {
	Title string
}

type Keys struct {
	pubKey  *crypto.PublicKey
	privKey *rsa.PrivateKey
}

func loadKeys() *Keys {
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
	privKey := parseResult.(*rsa.PrivateKey)
	pubKey := privKey.Public()
	//fmt.Printf("%T %T", privKey, pubKey)
	return &Keys{
		pubKey:  &pubKey,
		privKey: privKey,
	}
}

func (m *MagicToken) Encrypt() {
	keys := loadKeys()
	fmt.Println(*keys.pubKey, *keys.privKey)
}

func Create() *MagicToken {
	return &MagicToken{
		Title: "Test",
	}
}
