package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

var Validate = validator.New()

// this function only use in method POST that contain req body
func ParseJSON[T any](r *http.Request, structPayload *T) error {
	if r.Body == nil {
		log.Fatal("missing req body")
	}
	// json.NewDecoder will fill the struct
	return json.NewDecoder(r.Body).Decode(&structPayload)
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	// json.NewDecoder will fill the struct
	return json.NewEncoder(w).Encode(v)
}

func WriteError(w http.ResponseWriter, status int, err error) {
	WriteJSON(w, status, map[string]string{"error": err.Error()})
}

func FormatValidationError(err error) map[string]string {
	errors := make(map[string]string)
	if ve, ok := err.(validator.ValidationErrors); ok {
		for _, fe := range ve {
			field := strings.ToLower(fe.Field())
			switch fe.Tag() {
			case "required":
				errors[field] = fmt.Sprintf("%s is required", field)
			case "min":
				errors[field] = fmt.Sprintf("%s must be at least %s characters", field, fe.Param())
			case "max":
				errors[field] = fmt.Sprintf("%s cannot be longer than %s characters", field, fe.Param())
			default:
				errors[field] = fmt.Sprintf("%s is not valid", field)
			}
		}
	}
	return errors
}

func HashPasswordBcrypt(password string) (string, error) {

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func ComparePasswordBcrypt(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
