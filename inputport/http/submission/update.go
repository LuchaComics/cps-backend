package submission

import (
	"context"
	"encoding/json"
	"net/http"

	sub_s "github.com/LuchaComics/cps-backend/app/submission/datastore"
	"github.com/LuchaComics/cps-backend/utils/httperror"
)

func UnmarshalUpdateRequest(ctx context.Context, r *http.Request) (*sub_s.Submission, error) {
	// Initialize our array which will store all the results from the remote server.
	var requestData sub_s.Submission

	defer r.Body.Close()

	// Read the JSON string and convert it into our golang stuct else we need
	// to send a `400 Bad Request` errror message back to the client,
	err := json.NewDecoder(r.Body).Decode(&requestData) // [1]
	if err != nil {
		return nil, httperror.NewForSingleField(http.StatusBadRequest, "non_field_error", "payload structure is wrong")
	}

	// Perform our validation and return validation error on any issues detected.
	if err := ValidateUpdateRequest(&requestData); err != nil {
		return nil, err
	}

	return &requestData, nil
}

func ValidateUpdateRequest(dirtyData *sub_s.Submission) error {
	e := make(map[string]string)

	// if dirtyData.ServiceType == 0 {
	// 	e["service_type"] = "missing value"
	// }
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
	if dirtyData.PublisherName == "" {
		e["publisher_name"] = "missing choice"
	} else if dirtyData.PublisherName == "Other" && dirtyData.PublisherNameOther == "" {
		e["publisher_name_other"] = "missing choice"
	}
	if dirtyData.CreasesFinding == "" {
		e["creases_finding"] = "missing choice"
	}
	if dirtyData.TearsFinding == "" {
		e["tears_finding"] = "missing choice"
	}
	if dirtyData.MissingPartsFinding == "" {
		e["missing_parts_finding"] = "missing choice"
	}
	if dirtyData.StainsFinding == "" {
		e["stains_finding"] = "missing choice"
	}
	if dirtyData.DistortionFinding == "" {
		e["distortion_finding"] = "missing choice"
	}
	if dirtyData.PaperQualityFinding == "" {
		e["paper_quality_finding"] = "missing choice"
	}
	if dirtyData.SpineFinding == "" {
		e["spine_finding"] = "missing choice"
	}
	if dirtyData.CoverFinding == "" {
		e["cover_finding"] = "missing choice"
	}
	if dirtyData.OverallLetterGrade == "" {
		e["overall_letter_grade"] = "missing value"
	}
	if dirtyData.OverallGrade <= 0 {
		e["overall_grade"] = "missing value"
	}
	if dirtyData.OverallGrade > 10 {
		e["overall_grade"] = "value must be less then or equal to 10"
	}

	if dirtyData.ShowsSignsOfTamperingOrRestoration != sub_s.YesItShowsSignsOfTamperingOrRestoration && dirtyData.ShowsSignsOfTamperingOrRestoration != sub_s.NoItDoesNotShowsSignsOfTamperingOrRestoration {
		e["shows_signs_of_tampering_or_restoration"] = "missing value"
	}
	if len(e) != 0 {
		return httperror.NewForBadRequest(&e)
	}
	return nil
}

func (h *Handler) UpdateByID(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()

	data, err := UnmarshalUpdateRequest(ctx, r)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	err = h.Controller.UpdateByID(ctx, data)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	MarshalUpdateResponse(data, w)
}

func MarshalUpdateResponse(res *sub_s.Submission, w http.ResponseWriter) {
	if err := json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
