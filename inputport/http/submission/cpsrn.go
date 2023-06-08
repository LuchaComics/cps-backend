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
	SpecialNotesLine1                  string    `bson:"special_notes_line_1" json:"special_notes_line_1"`
	SpecialNotesLine2                  string    `bson:"special_notes_line_2" json:"special_notes_line_2"`
	SpecialNotesLine3                  string    `bson:"special_notes_line_3" json:"special_notes_line_3"`
	SpecialNotesLine4                  string    `bson:"special_notes_line_4" json:"special_notes_line_4"`
	SpecialNotesLine5                  string    `bson:"special_notes_line_5" json:"special_notes_line_5"`
	GradingNotesLine1                  string    `bson:"grading_notes_line_1" json:"grading_notes_line_1"`
	GradingNotesLine2                  string    `bson:"grading_notes_line_2" json:"grading_notes_line_2"`
	GradingNotesLine3                  string    `bson:"grading_notes_line_3" json:"grading_notes_line_3"`
	GradingNotesLine4                  string    `bson:"grading_notes_line_4" json:"grading_notes_line_4"`
	GradingNotesLine5                  string    `bson:"grading_notes_line_5" json:"grading_notes_line_5"`
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
		SpecialNotesLine1:                  s.SpecialNotesLine1,
		SpecialNotesLine2:                  s.SpecialNotesLine2,
		SpecialNotesLine3:                  s.SpecialNotesLine3,
		SpecialNotesLine4:                  s.SpecialNotesLine4,
		SpecialNotesLine5:                  s.SpecialNotesLine5,
		GradingNotesLine1:                  s.GradingNotesLine1,
		GradingNotesLine2:                  s.GradingNotesLine2,
		GradingNotesLine3:                  s.GradingNotesLine3,
		GradingNotesLine4:                  s.GradingNotesLine4,
		GradingNotesLine5:                  s.GradingNotesLine5,
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
