package handlers

import (
	"fmt"
	"errors"
	"time"
	"os"
	"strings"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/dgrijalva/jwt-go"
)

// GoogleClaims -
type GoogleClaims struct {
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	FirstName     string `json:"given_name"`
	LastName      string `json:"family_name"`
	jwt.StandardClaims
}

func GoogleLogin(c *fiber.Ctx) error {
	// Validate the JWT is valid
	if c.Get("Authorization") == "" {
		return c.Status(403).JSON(fiber.Map{"error":"Authorization not found"})
	} 
	claims, err := ValidateGoogleJWT(strings.Split(c.Get("Authorization"), " ")[1])
	if err != nil {
		return c.Status(403).JSON(fiber.Map{"error": err.Error(), "message": "Invalid google auth"})
	}
	if claims.Email != "marco.urriola@gmail.com" || claims.StandardClaims.Subject != "105528794879915282379" {
		return c.Status(403).JSON(fiber.Map{"message": "Unauthorized"})
	}
	// Go to next middleware:
	return c.Next()
  }

func getGooglePublicKey(keyID string) (string, error) {
	resp, err := http.Get("https://www.googleapis.com/oauth2/v1/certs")
	if err != nil {
		return "", err
	}
	dat, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	myResp := map[string]string{}
	err = json.Unmarshal(dat, &myResp)
	if err != nil {
		return "", err
	}
	key, ok := myResp[keyID]
	if !ok {
		return "", errors.New("key not found")
	}
	return key, nil
}

// ValidateGoogleJWT -
func ValidateGoogleJWT(tokenString string) (GoogleClaims, error) {
	claimsStruct := GoogleClaims{}
	token, err := jwt.ParseWithClaims(
		tokenString,
		&claimsStruct,
		func(token *jwt.Token) (interface{}, error) {
			pem, err := getGooglePublicKey(fmt.Sprintf("%s", token.Header["kid"]))
			if err != nil {
				return nil, err
			}
			key, err := jwt.ParseRSAPublicKeyFromPEM([]byte(pem))
			if err != nil {
				return nil, err
			}
			return key, nil
		},
	)
	if err != nil {
		return GoogleClaims{}, err
	}

	claims, ok := token.Claims.(*GoogleClaims)
	if !ok {
		return GoogleClaims{}, errors.New("Invalid Google JWT")
	}

	if claims.Issuer != "accounts.google.com" {
		return GoogleClaims{}, errors.New("iss is invalid")
	}

	if claims.Audience != os.Getenv("API_GOOGLE_CLIENT") {
		return GoogleClaims{}, errors.New("aud is invalid")
	}

	if claims.ExpiresAt < time.Now().UTC().Unix() {
		return GoogleClaims{}, errors.New("JWT is expired")
	}

	return *claims, nil
}