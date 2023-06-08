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
	if dirtyData.IssueCoverYear == "" {
		e["issue_cover_year"] = "missing value"
	}
	if dirtyData.IssueCoverMonth == 0 {
		e["issue_cover_month"] = "missing value"
	}
	if dirtyData.PublisherName < 1 || dirtyData.PublisherName > 9 {
		e["publisher_name"] = "missing choice"
	} else if dirtyData.PublisherName == 1 && dirtyData.PublisherNameOther == "" {
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
	if dirtyData.GradingScale <= 0 || dirtyData.GradingScale > 3 {
		e["grading_scale"] = "missing choice"
	} else {
		if dirtyData.OverallLetterGrade == "" && dirtyData.GradingScale == sub_s.LetterGradeScale {
			e["overall_letter_grade"] = "missing value"
		}
		if dirtyData.OverallNumberGrade <= 0 && dirtyData.OverallNumberGrade > 10 && dirtyData.GradingScale == sub_s.NumberGradeScale {
			e["overall_number_grade"] = "missing value"
		}
		if dirtyData.CpsPercentageGrade < 5 && dirtyData.CpsPercentageGrade > 100 && dirtyData.GradingScale == sub_s.CPSPercentageGradingScale {
			e["cps_percentage_grade"] = "missing value"
		}
	}
	if dirtyData.ShowsSignsOfTamperingOrRestoration != sub_s.YesItShowsSignsOfTamperingOrRestoration && dirtyData.ShowsSignsOfTamperingOrRestoration != sub_s.NoItDoesNotShowsSignsOfTamperingOrRestoration {
		e["shows_signs_of_tampering_or_restoration"] = "missing value"
	}

	// Process optional validation
	if dirtyData.SpecialNotesLine1 != "" && len(dirtyData.SpecialNotesLine1) > 35 {
		e["special_notes_line_1"] = "over 35 characters"
	}
	if dirtyData.SpecialNotesLine2 != "" && len(dirtyData.SpecialNotesLine2) > 35 {
		e["special_notes_line_2"] = "over 35 characters"
	}
	if dirtyData.SpecialNotesLine3 != "" && len(dirtyData.SpecialNotesLine3) > 35 {
		e["special_notes_line_3"] = "over 35 characters"
	}
	if dirtyData.SpecialNotesLine4 != "" && len(dirtyData.SpecialNotesLine4) > 35 {
		e["special_notes_line_4"] = "over 35 characters"
	}
	if dirtyData.SpecialNotesLine5 != "" && len(dirtyData.SpecialNotesLine5) > 35 {
		e["special_notes_line_5"] = "over 35 characters"
	}
	if dirtyData.GradingNotesLine1 != "" && len(dirtyData.GradingNotesLine1) > 35 {
		e["grading_notes_line_1"] = "over 35 characters"
	}
	if dirtyData.GradingNotesLine2 != "" && len(dirtyData.GradingNotesLine2) > 35 {
		e["grading_notes_line_2"] = "over 35 characters"
	}
	if dirtyData.GradingNotesLine3 != "" && len(dirtyData.GradingNotesLine3) > 35 {
		e["grading_notes_line_3"] = "over 35 characters"
	}
	if dirtyData.GradingNotesLine4 != "" && len(dirtyData.GradingNotesLine4) > 35 {
		e["grading_notes_line_4"] = "over 35 characters"
	}
	if dirtyData.GradingNotesLine5 != "" && len(dirtyData.GradingNotesLine5) > 35 {
		e["grading_notes_line_5"] = "over 35 characters"
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

	submission, err := h.Controller.UpdateByID(ctx, data)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	MarshalUpdateResponse(submission, w)
}

func MarshalUpdateResponse(res *sub_s.Submission, w http.ResponseWriter) {
	if err := json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
