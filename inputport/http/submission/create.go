package submission

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	sub_s "github.com/LuchaComics/cps-backend/app/submission/datastore"
	"github.com/LuchaComics/cps-backend/utils/errorx"
)

func UnmarshalCreateRequest(ctx context.Context, r *http.Request) (*sub_s.Submission, error, int) {
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
	isValid, errStr := ValidateCreateRequest(&requestData)
	if isValid == false {
		return nil, errors.New(errStr), http.StatusBadRequest
	}

	return &requestData, nil, http.StatusOK
}

func ValidateCreateRequest(dirtyData *sub_s.Submission) (bool, string) {
	e := make(map[string]string)

	if dirtyData.ServiceType == 0 {
		e["service_type"] = "missing value"
	}
	if dirtyData.SeriesTitle == "" {
		e["series_title"] = "missing value"
	}
	if dirtyData.IssueVol == "" {
		e["issue_vol"] = "missing value"
	}
	if dirtyData.IssueNo == "" {
		e["issue_no"] = "missing value"
	}
	if dirtyData.IssueCoverDate == "" {
		e["issue_cover_date"] = "missing value"
	}
	if dirtyData.CreasesFinding == 0 {
		e["creases_finding"] = "missing value"
	}
	if dirtyData.TearsFinding == 0 {
		e["tears_finding"] = "missing value"
	}
	if dirtyData.MissingPartsFinding == 0 {
		e["missing_parts_finding"] = "missing value"
	}
	if dirtyData.StainsFinding == 0 {
		e["stains_finding"] = "missing value"
	}
	if dirtyData.DistortionFinding == 0 {
		e["distortion_finding"] = "missing value"
	}
	if dirtyData.PaperQualityFinding == 0 {
		e["paper_quality_finding"] = "missing value"
	}
	if dirtyData.SpineFinding == 0 {
		e["spine_finding"] = "missing value"
	}
	if dirtyData.OtherFinding != 0 {
		if dirtyData.OtherFindingText == "" {
			e["other_finding_text"] = "missing value"
		}
	}
	if dirtyData.OverallLetterGrade == "" {
		e["overall_letter_grade"] = "missing value"
	}
	if dirtyData.UserID == "" {
		e["user_id"] = "missing value"
	}
	if dirtyData.UserFirstName == "" {
		e["user_first_name"] = "missing value"
	}
	if dirtyData.UserLastName == "" {
		e["user_last_name"] = "missing value"
	}
	if dirtyData.UserCompanyName == "" {
		e["user_company_name"] = "missing value"
	}
	if dirtyData.UserSignature == "" {
		e["user_signature"] = "missing value"
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

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	requestData, err, errStatusCode := UnmarshalCreateRequest(ctx, r)
	if err != nil {
		http.Error(w, err.Error(), errStatusCode)
		return
	}

	err = h.Controller.Create(ctx, requestData)
	if err != nil {
		errorx.ResponseError(w, err)
		return
	}

	MarshalCreateResponse(requestData, w)
}

func MarshalCreateResponse(res *sub_s.Submission, w http.ResponseWriter) {
	if err := json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
