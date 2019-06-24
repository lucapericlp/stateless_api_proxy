package magictoken

import (
	"../keys"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"fmt"
	"log"
)

type MagicToken struct {
	Title string
}

func encrypt(ghToken *string, pubKey *rsa.PublicKey) string {
	//fmt.Println(*ourKeys.PubKey, *ourKeys.PrivKey)
	encryptedToken, err := rsa.EncryptOAEP(sha512.New(), rand.Reader, pubKey, []byte(*ghToken), []byte(""))
	if err != nil {
		log.Fatalf("Encryption failed: %s", err)
	}
	return string(encryptedToken)
}

func decrypt(encryptedToken *string, privKey *rsa.PrivateKey) string {
	decryptedToken, err := rsa.DecryptOAEP(sha512.New(), rand.Reader, privKey, []byte(*encryptedToken), []byte(""))
	if err != nil {
		log.Fatalf("Decryption failed: %s", err)
	}
	return string(decryptedToken)
}

func Create(ghToken string, scopes [2]string) *MagicToken {
	ourKeys := keys.LoadKeys()
	ctToken := encrypt(&ghToken, ourKeys.PubKey)
	ptToken := decrypt(&ctToken, ourKeys.PrivKey)
	fmt.Println(ctToken, "\n", ptToken)
	return &MagicToken{
		Title: "Test",
	}
}
