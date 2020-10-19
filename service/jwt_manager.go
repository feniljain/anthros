package service

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

//JWTManager is a JSOn web token manager
type JWTManager struct {
	secretKey     string
	tokenDuration time.Duration
}

//UserClaims is a custom JWT claims that contains some user's information
type UserClaims struct {
	jwt.StandardClaims
	Username string
	Role     string
}

//NewJWTManager returns a new JWT manager
func NewJWTManager(secretKey string, tokenDuration time.Duration) *JWTManager {
	return &JWTManager{
		secretKey:     secretKey,
		tokenDuration: tokenDuration,
	}
}

//Generate generates and signs new token for a user
func (manager *JWTManager) Generate(user *User) (string, error) {
	claims := UserClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(manager.tokenDuration).Unix(),
		},
		Username: user.Username,
		Role:     user.Role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(manager.secretKey))
}

//Verify verifies the access token string and return a user claim if the token is valid
func (manager *JWTManager) Verify(accessToken string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(
		accessToken,
		&UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, fmt.Errorf("Unexpected token signing method")
			}
			return []byte(manager.secretKey), nil
		},
	)

	if err != nil {
		return nil, fmt.Errorf("Invalid token: %v", err)
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return nil, fmt.Errorf("Invalid user claims: %v", err)
	}

	return claims, nil
}
