package comicsub

import (
	"encoding/json"
	"net/http"

	sub_s "github.com/LuchaComics/cps-backend/app/comicsub/datastore"
	"github.com/LuchaComics/cps-backend/utils/httperror"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Initialize the list filter with base results and then override them with the URL parameters.
	f := &sub_s.ComicSubmissionListFilter{
		PageSize:        10,
		LastID:          "",
		SortField:       "_id",
		ExcludeArchived: true,
	}

	// Here is where you extract url parameters.
	query := r.URL.Query()
	organizationID := query.Get("organization_id")
	if organizationID != "" {
		organizationID, err := primitive.ObjectIDFromHex(organizationID)
		if err != nil {
			httperror.ResponseError(w, err)
			return
		}
		f.OrganizationID = organizationID
	}

	userID := query.Get("user_id")
	if userID != "" {
		userID, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			httperror.ResponseError(w, err)
			return
		}
		f.UserID = userID
	}

	// Fet
	m, err := h.Controller.ListByFilter(ctx, f)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	MarshalListResponse(m, w)
}

func MarshalListResponse(res *sub_s.ComicSubmissionListResult, w http.ResponseWriter) {
	if err := json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
