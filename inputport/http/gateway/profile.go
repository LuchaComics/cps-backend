package gateway

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	user_s "github.com/LuchaComics/cps-backend/app/user/datastore"
	"github.com/LuchaComics/cps-backend/utils/httperror"
)

func (h *Handler) Profile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	profile, err := h.Controller.Profile(ctx)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}
	MarshalProfileResponse(profile, w)
}

func MarshalProfileResponse(responseData *user_s.User, w http.ResponseWriter) {
	if err := json.NewEncoder(w).Encode(&responseData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type ProfileUpdateRequestIDO struct {
	FirstName                 string `bson:"first_name" json:"first_name"`
	LastName                  string `bson:"last_name" json:"last_name"`
	Email                     string `json:"email"`
	CompanyName               string `bson:"company_name,omitempty" json:"company_name,omitempty"`
	Phone                     string `bson:"phone,omitempty" json:"phone,omitempty"`
	Country                   string `bson:"country,omitempty" json:"country,omitempty"`
	Region                    string `bson:"region,omitempty" json:"region,omitempty"`
	City                      string `bson:"city,omitempty" json:"city,omitempty"`
	PostalCode                string `bson:"postal_code,omitempty" json:"postal_code,omitempty"`
	AddressLine1              string `bson:"address_line_1,omitempty" json:"address_line_1,omitempty"`
	AddressLine2              string `bson:"address_line_2,omitempty" json:"address_line_2,omitempty"`
	HowDidYouHearAboutUs      int8   `bson:"how_did_you_hear_about_us,omitempty" json:"how_did_you_hear_about_us,omitempty"`
	HowDidYouHearAboutUsOther string `bson:"how_did_you_hear_about_us_other,omitempty" json:"how_did_you_hear_about_us_other,omitempty"`
	AgreePromotionsEmail      bool   `bson:"agree_promotions_email,omitempty" json:"agree_promotions_email,omitempty"`
}

func UnmarshalProfileUpdateRequest(ctx context.Context, r *http.Request) (*user_s.User, error) {
	// Initialize our array which will store all the results from the remote server.
	var requestData ProfileUpdateRequestIDO

	defer r.Body.Close()

	// Read the JSON string and convert it into our golang stuct else we need
	// to send a `400 Bad Request` errror message back to the client,
	err := json.NewDecoder(r.Body).Decode(&requestData) // [1]
	if err != nil {
		return nil, httperror.NewForSingleField(http.StatusBadRequest, "non_field_error", "payload structure is wrong")
	}

	// Defensive Code: For security purposes we need to remove all whitespaces from the email and lower the characters.
	requestData.Email = strings.ToLower(requestData.Email)
	requestData.Email = strings.ReplaceAll(requestData.Email, " ", "")

	// Perform our validation and return validation error on any issues detected.
	if err = ValidateProfileUpdateRequest(&requestData); err != nil {
		return nil, err
	}

	// Convert to the user collection.
	return &user_s.User{
		FirstName:                 requestData.Email,
		LastName:                  requestData.Email,
		Email:                     requestData.Email,
		Phone:                     requestData.Phone,
		Country:                   requestData.Country,
		Region:                    requestData.Region,
		City:                      requestData.City,
		PostalCode:                requestData.PostalCode,
		AddressLine1:              requestData.AddressLine1,
		AddressLine2:              requestData.AddressLine2,
		HowDidYouHearAboutUs:      requestData.HowDidYouHearAboutUs,
		HowDidYouHearAboutUsOther: requestData.HowDidYouHearAboutUsOther,
		AgreePromotionsEmail:      requestData.AgreePromotionsEmail,
	}, nil
}

func ValidateProfileUpdateRequest(dirtyData *ProfileUpdateRequestIDO) error {
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
	if dirtyData.AddressLine1 == "" {
		e["address_line_1"] = "missing value"
	}

	if len(e) != 0 {
		return httperror.NewForBadRequest(&e)
	}
	return nil
}

func (h *Handler) ProfileUpdate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	data, err := UnmarshalProfileUpdateRequest(ctx, r)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	if err := h.Controller.ProfileUpdate(ctx, data); err != nil {
		httperror.ResponseError(w, err)
		return
	}

	// Get the request
	h.Profile(w, r)
}
