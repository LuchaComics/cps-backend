package comicsub

import (
	"encoding/json"
	"net/http"
	"strconv"

	sub_s "github.com/LuchaComics/cps-backend/app/comicsub/datastore"
	"github.com/LuchaComics/cps-backend/utils/httperror"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Initialize the list filter with base results and then override them with the URL parameters.
	f := &sub_s.ComicSubmissionListFilter{
		Cursor:          primitive.NilObjectID,
		PageSize:        25,
		SortField:       "_id",
		SortOrder:       1, // 1=ascending | -1=descending
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

	cursor := query.Get("cursor")
	if cursor != "" {
		cursor, err := primitive.ObjectIDFromHex(cursor)
		if err != nil {
			httperror.ResponseError(w, err)
			return
		}
		f.Cursor = cursor
	}

	pageSize := query.Get("page_size")
	if pageSize != "" {
		pageSize, _ := strconv.ParseInt(pageSize, 10, 64)
		if pageSize == 0 || pageSize > 250 {
			pageSize = 250
		}
		f.PageSize = pageSize
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
