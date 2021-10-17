package middleware

import (
	"YoutrackGChatBot/logging"
	"YoutrackGChatBot/settings"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
)

var VerifyGoogleAuth = func(f http.HandlerFunc) http.HandlerFunc {
	logger := logging.GetLogger()
	return func(w http.ResponseWriter, r *http.Request) {
		// Check the google token
		token, err := getAuthToken(r)

		if err != nil {
			logger.Println("Could not get authorization token")
			replyUnauthorized(w)
			return
		}

		settings, err := settings.GetSettings()

		if err != nil {
			logger.Println("Could not get settings")
			replyUnauthorized(w)
			return
		}

		if isValid, err := verifyToken(token, settings); isValid && err != nil {
			// If google auth passes, execute handler
			f(w, r)
		} else {
			logger.Println("Unauthorized access")
			replyUnauthorized(w)
			return
		}
	}
}

func replyUnauthorized(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	w.Write([]byte("<h1>Unauthorized access</h1>"))
}

func getAuthToken(req *http.Request) (*string, error) {
	var (
		token *string
		err   error
	)

	if val, ok := req.Header["Authorization"]; ok {
		if len(val) > 0 && len(strings.Split(val[0], " ")) > 1 {
			parts := strings.Split(val[0], " ")

			if parts[0] == "Bearer" {
				token = &parts[1]
			} else {
				err = errors.New("Invalid token type")
			}
		} else {
			err = errors.New("Invalid Authorization header")
		}
	} else {
		err = errors.New("Missing Authorization header")
	}

	return token, err
}

func verifyToken(token *string, s *settings.Settings) (bool, error) {
	claims := make(jwt.MapClaims)

	// Bring JWKS
	resp, err := http.Get(fmt.Sprintf("%s%s", s.PUBLIC_CERT_URL_PREFIX, s.GCHAT_ISSUER))
	if err != nil {
		return false, errors.New("Could not retrieve public certificate")
	}

	defer resp.Body.Close()

	var jsonObj map[string]interface{}

	// Decode JSON response
	err = json.NewDecoder(resp.Body).Decode(&jsonObj)

	if err != nil {
		return false, errors.New("Could not parse JSON response")
	}

	// jsonObj contains JSON response in a map
	decodedToken, err := jwt.ParseWithClaims(*token, claims, func(token *jwt.Token) (interface{}, error) {
		if kid, ok := token.Header["kid"].(string); ok {
			if key, ok := jsonObj[kid].(string); ok {
				return []byte(key), nil
			}
		}
		return nil, errors.New("No valid signing key")
	})

	if err != nil {
		return false, err
	}

	return decodedToken.Valid, nil
}
