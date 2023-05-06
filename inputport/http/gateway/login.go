package gateway

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	gateway_s "github.com/LuchaComics/cps-backend/app/gateway/datastore"
	"github.com/LuchaComics/cps-backend/utils/errorx"
)

type LoginRequestIDO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
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

	res, err := h.Controller.Login(ctx, requestData.Email, requestData.Password)
	if err != nil {
		errorx.ResponseError(w, err)
		return
	}
	MarshalLoginResponse(res, w)
}

func MarshalLoginResponse(responseData *gateway_s.LoginResponseIDO, w http.ResponseWriter) {
	if err := json.NewEncoder(w).Encode(&responseData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
