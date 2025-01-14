package auth

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/izumii.cxde/blog-api/config"
	"github.com/izumii.cxde/blog-api/types"
)

func ParseJWTRequest(r *http.Request) (int64, error) {
	c, err := r.Cookie("token")
	if err != nil {
		return 0, err
	}
	return ValidateJWTToken(c.Value)
}

/*
GenerateJWTToken generates a JWT token
@params: u(types.User) user info to generate the token
*/
func GenerateJWTToken(u types.User) (string, error) {
	expiration := time.Second * time.Duration(config.Envs.JWTExpiration)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":        u.ID,
		"expiredAt": expiration,
	})
	t, err := token.SignedString([]byte(config.Envs.JWTSecret))
	if err != nil {
		return "", err
	}
	return t, nil
}

/*
validate the token from the request.
@params: token(string) the token from the request
*/
func ValidateJWTToken(token string) (int64, error) {
	// Parse the token and provide the signing key for verification
	t, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		// Ensure the token is signed with the expected method (HS256 in this case)
		if t.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(config.Envs.JWTSecret), nil
	})

	if err != nil {
		return 0, fmt.Errorf("error parsing token: %v", err)
	}

	// Ensure the claims are valid and of the expected type
	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok || !t.Valid {
		return 0, fmt.Errorf("invalid token")
	}

	// Extract user ID from the claims
	userId, ok := claims["id"].(int64)
	if !ok {
		return 0, fmt.Errorf("invalid or missing value in token")
	}

	return userId, nil
}
