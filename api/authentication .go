package api

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pichuchen/hatsuaki/datastore/config"
)

func IssueJWT(username string) (string, error) {
	// 這邊我們會簽發一個 HS256 的 JWT

	// 這邊我們會使用 config 裡面的 secret 來簽發 JWT
	secret := config.GetLoginJWTSecret()
	if secret == "" {
		return "", fmt.Errorf("secret not found")
	}

	// 這邊我們會簽發一個 JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
	})
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil

}

func VerifyJWT(tokenString string) (string, error) {
	// 這邊我們會驗證一個 HS256 的 JWT

	// 這邊我們會使用 config 裡面的 secret 來驗證 JWT
	secret := config.GetLoginJWTSecret()
	if secret == "" {
		return "", fmt.Errorf("secret not found")
	}

	// 這邊我們會驗證 JWT
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 這邊我們只接受 HS256 的 token
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// 這邊我們會返回 secret 來驗證 token
		return []byte(secret), nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", fmt.Errorf("invalid token")
	}

	username, ok := claims["username"].(string)
	if !ok {
		return "", fmt.Errorf("username not found")
	}

	return username, nil
}
