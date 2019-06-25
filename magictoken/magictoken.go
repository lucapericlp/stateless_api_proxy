package magictoken

import (
	"../keys"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	//	"github.com/dgrijalva/jwt-go"
	"log"
	"time"
)

type MagicToken struct {
	iat          int64
	exp          int64
	github_token string
	scopes       [2]string
}

func encrypt(ghToken *string, pubKey *rsa.PublicKey) []byte {
	//fmt.Println(*ourKeys.PubKey, *ourKeys.PrivKey)
	encryptedToken, err := rsa.EncryptOAEP(sha512.New(), rand.Reader, pubKey, []byte(*ghToken), []byte(""))
	if err != nil {
		log.Fatalf("Encryption failed: %s", err)
	}
	return encryptedToken
}

func decrypt(encryptedToken *string, privKey *rsa.PrivateKey) []byte {
	decryptedToken, err := rsa.DecryptOAEP(sha512.New(), rand.Reader, privKey, []byte(*encryptedToken), []byte(""))
	if err != nil {
		log.Fatalf("Decryption failed: %s", err)
	}
	return decryptedToken
}

func Create(ghToken string, scopes [2]string) *MagicToken {
	ourKeys := keys.LoadKeys()
	ctToken := encrypt(&ghToken, ourKeys.PubKey)
	encodedCT := base64.StdEncoding.EncodeToString(ctToken)
	fmt.Println(encodedCT)

	issuedAt := time.Now()
	expiresAt := issuedAt.Add(time.Hour * 24 * 365)
	//ptToken := decrypt(&ctToken, ourKeys.PrivKey)
	//fmt.Println(ctToken, "\n", ptToken)
	ourToken := &MagicToken{
		iat:          issuedAt.Unix(),
		exp:          expiresAt.Unix(),
		github_token: encodedCT,
		scopes:       scopes,
	}

	//	jwtToken := jwt.NewWithClaims(jwt.SigningMethodRSA, jwt.MapClaims(*ourToken))

	return ourToken
}
