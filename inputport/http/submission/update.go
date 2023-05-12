package submission

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	sub_s "github.com/LuchaComics/cps-backend/app/submission/datastore"
	"github.com/LuchaComics/cps-backend/utils/httperror"
)

func UnmarshalUpdateRequest(ctx context.Context, r *http.Request) (*sub_s.Submission, error, int) {
	// Initialize our array which will store all the results from the remote server.
	var requestData sub_s.Submission

	defer r.Body.Close()

	// Read the JSON string and convert it into our golang stuct else we need
	// to send a `400 Bad Request` errror message back to the client,
	err := json.NewDecoder(r.Body).Decode(&requestData) // [1]
	if err != nil {
		return nil, err, http.StatusBadRequest
	}

	// Perform our validation and return validation error on any issues detected.
	isValid, errStr := ValidateUpdateRequest(&requestData)
	if isValid == false {
		return nil, errors.New(errStr), http.StatusBadRequest
	}

	return &requestData, nil, http.StatusOK
}

func ValidateUpdateRequest(dirtyData *sub_s.Submission) (bool, string) {
	e := make(map[string]string)

	if dirtyData.ServiceType == 0 {
		e["service_type"] = "missing value"
	}

	//TODO: Add more validation.

	if len(e) != 0 {
		b, err := json.Marshal(e)
		if err != nil { // Defensive code
			return false, err.Error()
		}
		return false, string(b)
	}
	return true, ""
}

func (h *Handler) UpdateBySubmissionID(w http.ResponseWriter, r *http.Request, submissionID string) {
	ctx := r.Context()

	requestData, err, errStatusCode := UnmarshalUpdateRequest(ctx, r)
	if err != nil {
		http.Error(w, err.Error(), errStatusCode)
		return
	}

	err = h.Controller.UpdateBySubmissionID(ctx, requestData)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	MarshalUpdateResponse(requestData, w)
}

func MarshalUpdateResponse(res *sub_s.Submission, w http.ResponseWriter) {
	if err := json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
