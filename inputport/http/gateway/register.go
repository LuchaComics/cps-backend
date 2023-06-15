package gateway

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	gateway_s "github.com/LuchaComics/cps-backend/app/gateway/datastore"
	"github.com/LuchaComics/cps-backend/utils/httperror"
)

func UnmarshalRegisterRequest(ctx context.Context, r *http.Request) (*gateway_s.RegisterRequestIDO, error) {
	// Initialize our array which will store all the results from the remote server.
	var requestData gateway_s.RegisterRequestIDO

	defer r.Body.Close()

	// Read the JSON string and convert it into our golang stuct else we need
	// to send a `400 Bad Request` errror message back to the client,
	err := json.NewDecoder(r.Body).Decode(&requestData) // [1]
	if err != nil {
		log.Println("Register | UnmarshalRegisterRequest | NewDecoder/Decode | err:", err)
		return nil, httperror.NewForSingleField(http.StatusBadRequest, "non_field_error", "payload structure is wrong")
	}

	// Defensive Code: For security purposes we need to remove all whitespaces from the email and lower the characters.
	requestData.Email = strings.ToLower(requestData.Email)
	requestData.Email = strings.ReplaceAll(requestData.Email, " ", "")

	// Perform our validation and return validation error on any issues detected.
	if err := ValidateRegisterRequest(&requestData); err != nil {
		return nil, err
	}

	return &requestData, nil
}

func ValidateRegisterRequest(dirtyData *gateway_s.RegisterRequestIDO) error {
	e := make(map[string]string)

	if dirtyData.FirstName == "" {
		e["first_name"] = "missing value"
	}
	if dirtyData.LastName == "" {
		e["last_name"] = "missing value"
	}
	if dirtyData.Email == "" {
		e["email"] = "missing value"
	}
	if len(dirtyData.Email) > 255 {
		e["email"] = "too long"
	}
	if dirtyData.Password == "" {
		e["password"] = "missing value"
	}
	if dirtyData.PasswordRepeated == "" {
		e["password_repeated"] = "missing value"
	}
	if dirtyData.PasswordRepeated != dirtyData.Password {
		e["password"] = "does not match"
		e["password_repeated"] = "does not match"
	}
	if dirtyData.CompanyName == "" {
		e["company_name"] = "missing value"
	}
	if dirtyData.Phone == "" {
		e["phone"] = "missing value"
	}
	if dirtyData.Country == "" {
		e["country"] = "missing value"
	}
	if dirtyData.Region == "" {
		e["region"] = "missing value"
	}
	if dirtyData.City == "" {
		e["city"] = "missing value"
	}
	if dirtyData.PostalCode == "" {
		e["postal_code"] = "missing value"
	}
	if dirtyData.Password == "" {
		e["password"] = "missing value"
	}
	if dirtyData.AddressLine1 == "" {
		e["address_line_1"] = "missing value"
	}
	// if dirtyData.HowDidYouHearAboutUs == 0 {
	// 	e["how_did_you_hear_about_us"] = "missing value"
	// }
	if dirtyData.AgreeTOS == false {
		e["agree_tos"] = "you must agree to the terms before proceeding"
	}
	if dirtyData.HowDidYouHearAboutUs > 7 || dirtyData.HowDidYouHearAboutUs < 1 {
		e["how_did_you_hear_about_us"] = "missing value"
	} else {
		if dirtyData.HowDidYouHearAboutUs == 1 && dirtyData.HowDidYouHearAboutUsOther == "" {
			e["how_did_you_hear_about_us_other"] = "missing value"
		}
	}

	if len(e) != 0 {
		return httperror.NewForBadRequest(&e)
	}
	return nil
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	data, err := UnmarshalRegisterRequest(ctx, r)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}
	if err := h.Controller.Register(ctx, data); err != nil {
		httperror.ResponseError(w, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
