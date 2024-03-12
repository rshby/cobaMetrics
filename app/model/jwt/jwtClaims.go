package jwt

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	Email            string                `json:"email,omitempty"`
	RegisteredClaims *jwt.RegisteredClaims `json:"registered_claims,omitempty"`
}

func (c *Claims) GetExpirationTime() (*jwt.NumericDate, error) {
	return c.RegisteredClaims.ExpiresAt, nil
}

func (c *Claims) GetIssuedAt() (*jwt.NumericDate, error) {
	return c.RegisteredClaims.IssuedAt, nil
}

func (c *Claims) GetNotBefore() (*jwt.NumericDate, error) {
	return c.RegisteredClaims.NotBefore, nil
}

func (c *Claims) GetIssuer() (string, error) {
	return c.RegisteredClaims.Issuer, nil
}

func (c *Claims) GetSubject() (string, error) {
	return c.RegisteredClaims.Subject, nil
}

func (c *Claims) GetAudience() (jwt.ClaimStrings, error) {
	return c.RegisteredClaims.Audience, nil
}
