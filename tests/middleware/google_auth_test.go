package middleware_test

import (
	"YoutrackGChatBot/middleware"
	"YoutrackGChatBot/settings"
	"fmt"
	"log"
	"os"
	"testing"

	"net/http"
	"net/http/httptest"

	"github.com/golang-jwt/jwt"
)

func TestGoogleAuthMiddlewareRequestWithoutAuthHeader(t *testing.T) {
	req, err := http.NewRequest("GET", "/test_endpoint", nil)

	if err != nil {
		t.Fatal(err)
	}

	// Create a response recorder to record the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(middleware.ApplyMiddleware(testHandler))

	// Serve HTTP from handler
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("Handler returned wrong status code. Got %v, want %v", status, http.StatusUnauthorized)
	}
}

func TestGoogleAuthMiddlewareRequestWithBadAuthHeader(t *testing.T) {
	req, err := http.NewRequest("GET", "/test_endpoint", nil)

	req.Header.Add("Authorization", "Bearer ofthebadtoken")

	if err != nil {
		t.Fatal(err)
	}

	// Create a response recorder to record the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(middleware.ApplyMiddleware(testHandler))

	// Serve HTTP from handler
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("Handler returned wrong status code. Got %v, want %v", status, http.StatusUnauthorized)
	} else {
		log.Printf("Status: %v, msg: %v", status, rr.Body)
	}
}

func testHandler(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hello, World!\n")
}

func TestValidToken(t *testing.T) {
	privKeyPath := "keys/test_key.key"
	certPath := "keys/test_key.cert"

	settings := settings.Settings{
		YOUTRACK_TOKEN:         "foo",
		GCHAT_ISSUER:           "foo",
		PUBLIC_CERT_URL_PREFIX: "foo",
		GCHAT_AUDIENCE:         "foobar",
	}

	claims := jwt.MapClaims{
		"aud": settings.GCHAT_AUDIENCE,
	}
	tokenString := generateTestJWT(privKeyPath, &claims)

	keyFunc := func(token *jwt.Token) (interface{}, error) {
		return getPubKeyFromCertFile(certPath)
	}

	valid, err := middleware.VerifyToken(tokenString, settings, keyFunc)

	if err != nil || valid == false {
		t.Errorf("ValidToken returned %t, %s\n", valid, err)
	}
}

func TestExtractPubFromX509(t *testing.T) {
	certPath := "keys/test_key.cert"
	pubKeyPath := "keys/test_key.pub"

	// Wrapper around readBytesFromFile and middleware.ExtractPubFromX509
	certPubKey, err := getPubKeyFromCertFile(certPath)

	if err != nil {
		t.Errorf("Error found: %s\n", err)
	}

	pubKeyBytes, err := readBytesFromFile(pubKeyPath)

	if err != nil {
		t.Errorf("Error found: %s\n", err)
	}

	pubKey, err := jwt.ParseRSAPublicKeyFromPEM(pubKeyBytes)

	if err != nil {
		t.Errorf("Error found: %s\n", err)
	}

	if !pubKey.Equal(certPubKey) {
		t.Errorf(`\
			Public key extracted from ExtractPubFromX509 does not match the expected result.
			Expected: %v.
			Got: %v.
		`, pubKey, certPubKey)
	}
}

func TestExtractPubFromX509ReturnsErrorOnWrongInput(t *testing.T) {
	if res, err := middleware.ExtractPubFromX509([]byte(`TheObviouslyWrongX509Format`)); err == nil {
		t.Errorf("Expected error on wrong input. Instead got %v\n", res)
	}
}

func generateTestJWT(privKey string, claims *jwt.MapClaims) string {
	// Read certificate into array of bytes
	keyBytes, err := readBytesFromFile(privKey)
	// Extract RSA Private key from file
	key, err := jwt.ParseRSAPrivateKeyFromPEM(keyBytes)

	if err != nil {
		log.Fatal(err)
	}

	// Create JWT and sign it
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, *claims)

	tokenString, err := token.SignedString(key)

	if err != nil {
		log.Fatal(err)
	}

	return tokenString
}

func readBytesFromFile(path string) ([]byte, error) {
	// Read certificate into array of bytes
	file, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	stats, err := file.Stat()

	if err != nil {
		return nil, err
	}

	size := stats.Size()

	keyBytes := make([]byte, size)

	_, err = file.Read(keyBytes)

	if err != nil {
		return nil, err
	}

	return keyBytes, nil
}

func getPubKeyFromCertFile(certPath string) (interface{}, error) {
	certBytes, err := readBytesFromFile(certPath)

	if err != nil {
		return nil, err
	}

	return middleware.ExtractPubFromX509(certBytes)
}
