package submission

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	sub_s "github.com/LuchaComics/cps-backend/app/submission/datastore"
	"github.com/LuchaComics/cps-backend/utils/httperror"
)

func UnmarshalCreateRequest(ctx context.Context, r *http.Request) (*sub_s.Submission, error) {
	// Initialize our array which will store all the results from the remote server.
	var requestData sub_s.Submission

	defer r.Body.Close()

	// Read the JSON string and convert it into our golang stuct else we need
	// to send a `400 Bad Request` errror message back to the client,
	err := json.NewDecoder(r.Body).Decode(&requestData) // [1]
	if err != nil {
		log.Println(err)
		return nil, httperror.NewForSingleField(http.StatusBadRequest, "non_field_error", "payload structure is wrong")
	}

	// Perform our validation and return validation error on any issues detected.
	if err := ValidateCreateRequest(&requestData); err != nil {
		return nil, err
	}
	return &requestData, nil
}

func ValidateCreateRequest(dirtyData *sub_s.Submission) error {
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
	if dirtyData.IssueCoverYear <= 0 {
		e["issue_cover_year"] = "missing value"
	}
	if dirtyData.IssueCoverMonth <= 0 {
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

	// Process optional validation for `Special Notes`.
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
	if dirtyData.SpecialNotesLine6 != "" && len(dirtyData.SpecialNotesLine6) > 35 {
		e["special_notes_line_6"] = "over 35 characters"
	}
	if dirtyData.SpecialNotesLine7 != "" && len(dirtyData.SpecialNotesLine7) > 35 {
		e["special_notes_line_7"] = "over 35 characters"
	}
	if dirtyData.SpecialNotesLine8 != "" && len(dirtyData.SpecialNotesLine8) > 35 {
		e["special_notes_line_8"] = "over 35 characters"
	}
	if dirtyData.SpecialNotesLine9 != "" && len(dirtyData.SpecialNotesLine9) > 35 {
		e["special_notes_line_9"] = "over 35 characters"
	}
	if dirtyData.SpecialNotesLine10 != "" && len(dirtyData.SpecialNotesLine10) > 35 {
		e["special_notes_line_10"] = "over 35 characters"
	}
	if dirtyData.SpecialNotesLine11 != "" && len(dirtyData.SpecialNotesLine11) > 35 {
		e["special_notes_line_11"] = "over 35 characters"
	}
	if dirtyData.SpecialNotesLine12 != "" && len(dirtyData.SpecialNotesLine12) > 35 {
		e["special_notes_line_12"] = "over 35 characters"
	}
	if dirtyData.SpecialNotesLine13 != "" && len(dirtyData.SpecialNotesLine13) > 35 {
		e["special_notes_line_13"] = "over 35 characters"
	}

	// Process optional validation for `Grading Notes`.
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
	if dirtyData.GradingNotesLine6 != "" && len(dirtyData.GradingNotesLine6) > 35 {
		e["grading_notes_line_6"] = "over 35 characters"
	}
	if dirtyData.GradingNotesLine7 != "" && len(dirtyData.GradingNotesLine7) > 35 {
		e["grading_notes_line_7"] = "over 35 characters"
	}
	if dirtyData.GradingNotesLine8 != "" && len(dirtyData.GradingNotesLine8) > 35 {
		e["grading_notes_line_8"] = "over 35 characters"
	}
	if dirtyData.GradingNotesLine9 != "" && len(dirtyData.GradingNotesLine9) > 35 {
		e["grading_notes_line_9"] = "over 35 characters"
	}
	if dirtyData.GradingNotesLine10 != "" && len(dirtyData.GradingNotesLine10) > 35 {
		e["grading_notes_line_10"] = "over 35 characters"
	}
	if dirtyData.GradingNotesLine11 != "" && len(dirtyData.GradingNotesLine11) > 35 {
		e["grading_notes_line_11"] = "over 35 characters"
	}
	if dirtyData.GradingNotesLine12 != "" && len(dirtyData.GradingNotesLine12) > 35 {
		e["grading_notes_line_12"] = "over 35 characters"
	}
	if dirtyData.GradingNotesLine13 != "" && len(dirtyData.GradingNotesLine13) > 35 {
		e["grading_notes_line_13"] = "over 35 characters"
	}

	if len(e) != 0 {
		return httperror.NewForBadRequest(&e)
	}
	return nil
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	data, err := UnmarshalCreateRequest(ctx, r)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	data, err = h.Controller.Create(ctx, data)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	MarshalCreateResponse(data, w)
}

func MarshalCreateResponse(res *sub_s.Submission, w http.ResponseWriter) {
	if err := json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
