package submission

import (
	"encoding/json"
	"net/http"
	"time"

	sub_s "github.com/LuchaComics/cps-backend/app/submission/datastore"
	"github.com/LuchaComics/cps-backend/utils/httperror"
)

func (h *Handler) GetRegistryByCPSRN(w http.ResponseWriter, r *http.Request, cpsn string) {
	ctx := r.Context()
	m, err := h.Controller.GetByCPSRN(ctx, cpsn)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	MarshalRegistryResponse(m, w)
}

// date issued, title, volume, issue number, comic cover date, signs of restoration (yes/no), special notes, grading notes, overall grade

type RegistryReponse struct {
	CPSRN                              string    `bson:"cpsrn" json:"cpsrn"`
	SubmissionDate                     time.Time `bson:"submission_date" json:"submission_date"`
	Item                               string    `bson:"item" json:"item"` // Created by system.
	SeriesTitle                        string    `bson:"series_title" json:"series_title"`
	IssueVol                           string    `bson:"issue_vol" json:"issue_vol"`
	IssueNo                            string    `bson:"issue_no" json:"issue_no"`
	IssueCoverYear                     int64     `bson:"issue_cover_year" json:"issue_cover_year"`
	IssueCoverMonth                    int8      `bson:"issue_cover_month" json:"issue_cover_month"`
	PublisherName                      int8      `bson:"publisher_name" json:"publisher_name"`
	PublisherNameOther                 string    `bson:"publisher_name_other" json:"publisher_name_other"`
	SpecialNotes                       string    `bson:"special_notes" json:"special_notes"`
	GradingNotes                       string    `bson:"grading_notes" json:"grading_notes"`
	ShowsSignsOfTamperingOrRestoration int8      `bson:"shows_signs_of_tampering_or_restoration" json:"shows_signs_of_tampering_or_restoration"`
	GradingScale                       int8      `bson:"grading_scale" json:"grading_scale"`
	OverallLetterGrade                 string    `bson:"overall_letter_grade" json:"overall_letter_grade"`
	OverallNumberGrade                 float64   `bson:"overall_number_grade" json:"overall_number_grade"`
	CpsPercentageGrade                 float64   `bson:"cps_percentage_grade" json:"cps_percentage_grade"`
}

func MarshalRegistryResponse(s *sub_s.Submission, w http.ResponseWriter) {
	resp := &RegistryReponse{
		CPSRN:                              s.CPSRN,
		SubmissionDate:                     s.SubmissionDate,
		Item:                               s.Item,
		SeriesTitle:                        s.SeriesTitle,
		IssueVol:                           s.IssueVol,
		IssueNo:                            s.IssueNo,
		IssueCoverYear:                     s.IssueCoverYear,
		IssueCoverMonth:                    s.IssueCoverMonth,
		PublisherName:                      s.PublisherName,
		PublisherNameOther:                 s.PublisherNameOther,
		SpecialNotes:                       s.SpecialNotes,
		GradingNotes:                       s.GradingNotes,
		ShowsSignsOfTamperingOrRestoration: s.ShowsSignsOfTamperingOrRestoration,
		GradingScale:                       s.GradingScale,
		OverallLetterGrade:                 s.OverallLetterGrade,
		OverallNumberGrade:                 s.OverallNumberGrade,
		CpsPercentageGrade:                 s.CpsPercentageGrade,
	}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
