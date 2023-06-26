package comicsub

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	sub_c "github.com/LuchaComics/cps-backend/app/comicsub/controller"
	sub_s "github.com/LuchaComics/cps-backend/app/comicsub/datastore"
	"github.com/LuchaComics/cps-backend/utils/httperror"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func UnmarshalUpdateRequest(ctx context.Context, r *http.Request) (*sub_c.ComicSubmissionUpdateRequestIDO, error) {
	// Initialize our array which will store all the results from the remote server.
	var requestData sub_c.ComicSubmissionUpdateRequestIDO

	defer r.Body.Close()

	// Read the JSON string and convert it into our golang stuct else we need
	// to send a `400 Bad Request` errror message back to the client,
	err := json.NewDecoder(r.Body).Decode(&requestData) // [1]
	if err != nil {
		log.Println(err)
		return nil, httperror.NewForSingleField(http.StatusBadRequest, "non_field_error", "payload structure is wrong")
	}

	// Perform our validation and return validation error on any issues detected.
	if err := ValidateUpdateRequest(&requestData); err != nil {
		return nil, err
	}

	return &requestData, nil
}

func ValidateUpdateRequest(dirtyData *sub_c.ComicSubmissionUpdateRequestIDO) error {
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
	if dirtyData.SpecialNotes != "" && len(dirtyData.SpecialNotes) > 638 {
		e["special_notes"] = "over 638 characters"
	}

	// Process optional validation for `Grading Notes`.
	if dirtyData.GradingNotes != "" && len(dirtyData.GradingNotes) > 638 {
		e["grading_notes"] = "over 638 characters"
	}
	if dirtyData.Status == 0 {
		e["status"] = "missing choice"
	}
	if dirtyData.ServiceType == 0 {
		e["service_type"] = "missing choice"
	}
	if dirtyData.OrganizationID.IsZero() {
		e["organization_id"] = "missing choice"
	}

	if len(e) != 0 {
		return httperror.NewForBadRequest(&e)
	}
	return nil
}

func (h *Handler) UpdateByID(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()

	d, err := UnmarshalUpdateRequest(ctx, r)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	d.ID, err = primitive.ObjectIDFromHex(id)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	submission, err := h.Controller.UpdateByID(ctx, d)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	log.Println("--->", submission)

	MarshalUpdateResponse(submission, w)
}

func MarshalUpdateResponse(res *sub_s.ComicSubmission, w http.ResponseWriter) {
	if err := json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
