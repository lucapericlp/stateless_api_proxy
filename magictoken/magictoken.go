package magictoken

import (
	"../keys"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"log"
	"regexp"
	"strings"
	"time"
)

type ProxyToken struct {
	GithubToken string
	Scopes      []string
	Iat         int64
	Exp         int64
}

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
		"github_token": encodedCT,
		"scopes":       scopes,
		"iat":          issuedAt.Unix(),
		"exp":          expiresAt.Unix(),
	})

	tokenString, err := jwtToken.SignedString(ourKeys.PrivKey)

	return tokenString, err
}

func Verify(tokenString string, ourKeys *keys.Keys) (*ProxyToken, error) {
	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		//not working but will look into
		//	if _, ok := token.Method.(*jwt.SigningMethodRS256); !ok {
		//		return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		//	}
		return ourKeys.PubKey, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		//token still valid? jwt claim seems to convert int64 to float64 which requires the type conversion here and at end of this statement.
		if int64(claims["exp"].(float64)) < time.Now().Unix() {
			return &ProxyToken{}, errors.New("EXPIRED TOKEN")
		}

		encodedCT := claims["github_token"].(string)
		decodedCT, err := base64.StdEncoding.DecodeString(encodedCT)

		if err != nil {
			fmt.Errorf("Error with base64 decoding: %s", err)
		}

		ptToken, err := decrypt(string(decodedCT), ourKeys.PrivKey)

		scopesInterface := claims["scopes"].([]interface{})
		scopes := make([]string, len(scopesInterface))
		for i, v := range scopesInterface {
			scopes[i] = fmt.Sprint(v)
		}

		proxyToken := &ProxyToken{
			GithubToken: string(ptToken),
			Scopes:      scopes,
			Iat:         int64(claims["iat"].(float64)),
			Exp:         int64(claims["exp"].(float64)),
		}

		return proxyToken, err
	}

	return &ProxyToken{}, errors.New("INVALID TOKEN")
}

func (p *ProxyToken) ValidateRequest(method string, path string) bool {
	validated := false
	for _, scope := range p.Scopes {
		scopeSplit := strings.Split(scope, " ")
		allowedMethod, allowedPath := scopeSplit[0], scopeSplit[1]

		if method != allowedMethod && allowedMethod != "*" {
			continue
		}

		if !strings.HasPrefix(path, "/") {
			prep := []string{"/"}
			path = strings.Join(append(prep, path), "")
		}

		re := regexp.MustCompile(allowedPath)
		if re.Match([]byte(path)) {
			validated = true
			break
		}
	}
	return validated
}
