package privy

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// PrivyClaims a Go type for Privy JWTs
type PrivyClaims struct {
	AppID      string `json:"aud,omitempty"`
	Expiration uint64 `json:"exp,omitempty"`
	Issuer     string `json:"iss,omitempty"`
	UserId     string `json:"sub,omitempty"`
}

// Add these methods to implement jwt.Claims interface
func (c *PrivyClaims) GetExpirationTime() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(time.Unix(int64(c.Expiration), 0)), nil
}

func (c *PrivyClaims) GetIssuedAt() (*jwt.NumericDate, error) {
	return nil, nil
}

func (c *PrivyClaims) GetNotBefore() (*jwt.NumericDate, error) {
	return nil, nil
}

func (c *PrivyClaims) GetIssuer() (string, error) {
	return c.Issuer, nil
}

func (c *PrivyClaims) GetSubject() (string, error) {
	return c.UserId, nil
}

func (c *PrivyClaims) GetAudience() (jwt.ClaimStrings, error) {
	return jwt.ClaimStrings{c.AppID}, nil
}

// This method will be used to check the token's claims later
func (pc *PrivyClient) Valid(c *PrivyClaims) error {
	if c.AppID != pc.apiID {
		return errors.New("aud claim must be your Privy App ID")
	}
	if c.Issuer != "privy.io" {
		return errors.New("iss claim must be 'privy.io'")
	}
	if c.Expiration < uint64(time.Now().Unix()) {
		return errors.New("token is expired")
	}

	return nil
}

// This method will be used to load the verification key in the required format later
func (pc *PrivyClient) keyFunc(token *jwt.Token) (any, error) {
	if token.Method.Alg() != "ES256" {
		return nil, fmt.Errorf("unexpected JWT signing method=%v", token.Header["alg"])
	}
	// https://pkg.go.dev/github.com/dgrijalva/jwt-go#ParseECPublicKeyFromPEM
	return jwt.ParseECPublicKeyFromPEM(pc.verificationKey)
}

// ValidAccessToken valid privy access token, decode privy user DID
func (pc *PrivyClient) ValidAccessToken(accessToken string) (*PrivyClaims, error) {
	// Check the JWT signature and decode claims
	// https://pkg.go.dev/github.com/dgrijalva/jwt-go#ParseWithClaims
	token, err := jwt.ParseWithClaims(accessToken, &PrivyClaims{}, pc.keyFunc)
	if err != nil {
		return nil, fmt.Errorf("JWT signature is invalid: %w", err)
	}

	// Parse the JWT claims into your custom struct
	privyClaim, ok := token.Claims.(*PrivyClaims)
	if !ok {
		// fmt.Println("JWT does not have all the necessary claims.")
		return nil, errors.New("JWT does not have all the necessary claims")
	}

	// fmt.Println("privy claims: ", privyClaim)

	// Check the JWT claims
	err = pc.Valid(privyClaim)
	if err != nil {
		return nil, err
	}

	return privyClaim, nil
}
