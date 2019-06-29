package magictoken

import (
	"../keys"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"log"
	"time"
)

//TO DO: Implement scope and token validation as middleware for api routes.

func encrypt(ghToken *string, pubKey *rsa.PublicKey) ([]byte, error) {
	encryptedToken, err := rsa.EncryptOAEP(sha512.New(), rand.Reader, pubKey, []byte(*ghToken), []byte(""))
	if err != nil {
		log.Fatalf("Encryption failed: %s", err)
	}
	return encryptedToken, err
}

func decrypt(encryptedToken string, privKey *rsa.PrivateKey) ([]byte, error) {
	decryptedToken, err := rsa.DecryptOAEP(sha512.New(), rand.Reader, privKey, []byte(encryptedToken), []byte(""))
	if err != nil {
		log.Fatalf("Decryption failed: %s", err)
	}
	return decryptedToken, err
}

func Create(ghToken string, scopes []string, ourKeys *keys.Keys) (string, error) {
	ctToken, _ := encrypt(&ghToken, ourKeys.PubKey)
	encodedCT := base64.StdEncoding.EncodeToString(ctToken)

	issuedAt := time.Now()
	expiresAt := issuedAt.Add(time.Hour * 24 * 365)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iat":          issuedAt.Unix(),
		"exp":          expiresAt.Unix(),
		"github_token": encodedCT,
		"scopes":       scopes,
	})

	tokenString, err := jwtToken.SignedString(ourKeys.PrivKey)

	return tokenString, err
}

func Verify(tokenString string, ourKeys *keys.Keys) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		//not working but will look into
		//	if _, ok := token.Method.(*jwt.SigningMethodRS256); !ok {
		//		return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		//	}
		return ourKeys.PubKey, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		encodedCT := claims["github_token"].(string)
		decodedCT, err := base64.StdEncoding.DecodeString(encodedCT)

		if err != nil {
			fmt.Errorf("Error with base64 decoding: %s", err)
		}

		ptToken, err := decrypt(string(decodedCT), ourKeys.PrivKey)
		return string(ptToken), err
	}

	return "NOT VALID TOKEN", err
}
