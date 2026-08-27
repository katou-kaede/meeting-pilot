package auth

import (
	"os"
	"testing"
)

func TestGenerateAndParseToken(t *testing.T) {
	t.Setenv(
		"JWT_SECRET",
		"test-secret-key-for-jwt-unit-test",
	)

	const userID int64 = 10

	tokenString, err := GenerateToken(userID)
	if err != nil {
		t.Fatalf(
			"JWTの生成に失敗しました: %v",
			err,
		)
	}

	claims, err := ParseToken(tokenString)
	if err != nil {
		t.Fatalf(
			"JWTの検証に失敗しました: %v",
			err,
		)
	}

	if claims.UserID != userID {
		t.Errorf(
			"UserID = %d, want %d",
			claims.UserID,
			userID,
		)
	}

	if claims.ExpiresAt == nil {
		t.Error("有効期限が設定されていません")
	}

	if claims.IssuedAt == nil {
		t.Error("発行日時が設定されていません")
	}
}

func TestParseToken_InvalidToken(t *testing.T) {
	t.Setenv(
		"JWT_SECRET",
		"test-secret-key-for-jwt-unit-test",
	)

	_, err := ParseToken("invalid-token")

	if err == nil {
		t.Fatal(
			"不正なJWTでエラーを期待しましたがnilでした",
		)
	}
}

func TestParseToken_DifferentSecret(t *testing.T) {
	t.Setenv("JWT_SECRET", "first-secret-key")

	tokenString, err := GenerateToken(10)
	if err != nil {
		t.Fatalf(
			"JWTの生成に失敗しました: %v",
			err,
		)
	}

	// 生成時とは異なる秘密鍵に変更
	t.Setenv("JWT_SECRET", "second-secret-key")

	_, err = ParseToken(tokenString)

	if err == nil {
		t.Fatal(
			"異なる秘密鍵でエラーを期待しましたがnilでした",
		)
	}
}

func TestGenerateToken_SecretNotSet(t *testing.T) {
	originalSecret := os.Getenv("JWT_SECRET")
	t.Cleanup(func() {
		_ = os.Setenv("JWT_SECRET", originalSecret)
	})

	_ = os.Unsetenv("JWT_SECRET")

	_, err := GenerateToken(10)

	if err == nil {
		t.Fatal(
			"JWT_SECRET未設定でエラーを期待しましたがnilでした",
		)
	}
}