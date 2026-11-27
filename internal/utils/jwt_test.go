package utils

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/spf13/viper"
)

func newTestJWT(secret string) *JWTHelper {
	v := viper.New()
	v.Set("JWT_SECRET", secret)
	return NewJWT(v)
}

func TestJWT_GenerateAndDecode(t *testing.T) {
	helper := newTestJWT("test-secret")
	id := uuid.New()

	token, err := helper.GenerateAccessToken(id, "user@example.com")
	if err != nil {
		t.Fatalf("GenerateAccessToken returned error: %v", err)
	}

	claims, err := helper.DecodeAccessToken(token)
	if err != nil {
		t.Fatalf("DecodeAccessToken returned error: %v", err)
	}
	if claims["sub"] != id.String() {
		t.Errorf("sub = %v, want %s", claims["sub"], id)
	}
	if claims["email"] != "user@example.com" {
		t.Errorf("email = %v, want user@example.com", claims["email"])
	}
	if claims["iss"] != "wallet-service" || claims["aud"] != "wallet-service" {
		t.Errorf("unexpected iss/aud: %v / %v", claims["iss"], claims["aud"])
	}
}

func TestJWT_GenerateWithEmptySecretFails(t *testing.T) {
	if _, err := newTestJWT("").GenerateAccessToken(uuid.New(), "user@example.com"); err == nil {
		t.Fatal("expected an error when the secret is empty")
	}
}

func TestJWT_DecodeRejectsTokenSignedWithAnotherSecret(t *testing.T) {
	token, err := newTestJWT("secret-a").GenerateAccessToken(uuid.New(), "user@example.com")
	if err != nil {
		t.Fatalf("GenerateAccessToken returned error: %v", err)
	}
	if _, err := newTestJWT("secret-b").DecodeAccessToken(token); err == nil {
		t.Fatal("expected an error for a token signed with a different secret")
	}
}

func TestJWT_DecodeRejectsExpiredToken(t *testing.T) {
	secret := "test-secret"
	expired := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": uuid.New().String(),
		"exp": time.Now().Add(-time.Minute).Unix(),
	})
	token, err := expired.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("signing test token: %v", err)
	}
	if _, err := newTestJWT(secret).DecodeAccessToken(token); err == nil {
		t.Fatal("expected an error for an expired token")
	}
}

func TestJWT_DecodeRejectsNoneAlgorithm(t *testing.T) {
	unsigned := jwt.NewWithClaims(jwt.SigningMethodNone, jwt.MapClaims{
		"sub": uuid.New().String(),
		"exp": time.Now().Add(time.Hour).Unix(),
	})
	token, err := unsigned.SignedString(jwt.UnsafeAllowNoneSignatureType)
	if err != nil {
		t.Fatalf("signing test token: %v", err)
	}
	if _, err := newTestJWT("test-secret").DecodeAccessToken(token); err == nil {
		t.Fatal("expected an error for a token using the none algorithm")
	}
}

func TestError_ImplementsHTTPError(t *testing.T) {
	cause := jwt.ErrTokenExpired
	err := Error("Wallet not found", 404, cause)

	httpErr, ok := err.(HTTPError)
	if !ok {
		t.Fatal("Error() should return a value implementing HTTPError")
	}
	if httpErr.Status() != 404 || httpErr.Message() != "Wallet not found" || httpErr.Unwrap() != cause {
		t.Errorf("unexpected HTTPError values: status=%d message=%q unwrap=%v", httpErr.Status(), httpErr.Message(), httpErr.Unwrap())
	}
}
