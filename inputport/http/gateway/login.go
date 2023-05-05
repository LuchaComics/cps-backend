package gateway

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	user_s "github.com/LuchaComics/cps-backend/app/user/datastore"
)

type LoginRequestIDO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponseIDO struct used to represent the system's response when the `login` POST request was a success.
type LoginResponseIDO struct {
	Email                  string    `json:"email"`
	AccessToken            string    `json:"access_token"`
	AccessTokenExpiryDate  time.Time `json:"access_token_expiry_date"`
	RefreshToken           string    `json:"refresh_token"`
	RefreshTokenExpiryDate time.Time `json:"refresh_token_expiry_date"`
}

func UnmarshalLoginRequest(ctx context.Context, r *http.Request) (*LoginRequestIDO, error, int) {
	// Initialize our array which will store all the results from the remote server.
	var requestData LoginRequestIDO

	defer r.Body.Close()

	// Read the JSON string and convert it into our golang stuct else we need
	// to send a `400 Bad Request` errror message back to the client,
	err := json.NewDecoder(r.Body).Decode(&requestData) // [1]
	if err != nil {
		return nil, err, http.StatusBadRequest
	}

	// Defensive Code: For security purposes we need to remove all whitespaces from the email and lower the characters.
	requestData.Email = strings.ToLower(requestData.Email)
	requestData.Email = strings.ReplaceAll(requestData.Email, " ", "")

	// Perform our validation and return validation error on any issues detected.
	isValid, errStr := ValidateLoginRequest(&requestData)
	if isValid == false {
		return nil, errors.New(errStr), http.StatusBadRequest
	}

	return &requestData, nil, http.StatusOK
}

func ValidateLoginRequest(dirtyData *LoginRequestIDO) (bool, string) {
	e := make(map[string]string)

	if dirtyData.Email == "" {
		e["email"] = "missing value"
	}
	if len(dirtyData.Email) > 255 {
		e["email"] = "too long"
	}
	if dirtyData.Password == "" {
		e["password"] = "missing value"
	}

	if len(e) != 0 {
		b, err := json.Marshal(e)
		if err != nil { // Defensive code
			return false, err.Error()
		}
		return false, string(b)
	}
	return true, ""
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	requestData, err, errStatusCode := UnmarshalLoginRequest(ctx, r)
	if err != nil {
		http.Error(w, err.Error(), errStatusCode)
		return
	}

	user, accessToken, accessTokenExpiryDate, refreshToken, refreshTokenExpiryDate, err := h.Controller.Login(ctx, requestData.Email, requestData.Password)
	if user == nil {
		http.Error(w, "{'non_field_error':'user does not exist'}", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusOK)
		return
	}

	MarshalLoginResponse(accessToken, accessTokenExpiryDate, refreshToken, refreshTokenExpiryDate, user, w)
}

func MarshalLoginResponse(accessToken string, accessTokenExpiryDate time.Time, refreshToken string, refreshTokenExpiryDate time.Time, u *user_s.User, w http.ResponseWriter) {
	// Generate our response.
	responseData := LoginResponseIDO{
		Email:                  u.Email,
		AccessToken:            accessToken,
		AccessTokenExpiryDate:  accessTokenExpiryDate,
		RefreshToken:           refreshToken,
		RefreshTokenExpiryDate: refreshTokenExpiryDate,
	}
	if err := json.NewEncoder(w).Encode(&responseData); err != nil {
		// log.Println("MarshalLoginResponse | Encode | err:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
