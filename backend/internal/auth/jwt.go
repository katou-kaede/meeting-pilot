package auth

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID               int64 `json:"user_id"`
	jwt.RegisteredClaims       // JWTの標準項目をまとめた構造体
}

// JWTを生成
func GenerateToken(userID int64) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", errors.New("JWT_SECRET is not set")
	}

	now := time.Now()

	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt: jwt.NewNumericDate(now), // 発行日時
			ExpiresAt: jwt.NewNumericDate( // 有効期限(24時間)
				now.Add(24 * time.Hour),
			),
		},
	}

	// 署名前のJWTオブジェクトを生成
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256, // 署名の方式
		claims,
	)

	// token に秘密鍵で署名し、最終的なJWT文字列を生成
	return token.SignedString([]byte(secret))
}

// JWTを検証
func ParseToken(tokenString string) (*Claims, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return nil, errors.New("JWT_SECRET is not set")
	}

	token, err := jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(token *jwt.Token) (any, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, errors.New(
					"unexpected signing method",
				)
			}

			return []byte(secret), nil
		},
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
