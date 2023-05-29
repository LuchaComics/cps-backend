package customer

import (
	"encoding/json"
	"fmt"
	"net/http"

	sub_s "github.com/LuchaComics/cps-backend/app/user/datastore"
	"github.com/LuchaComics/cps-backend/utils/httperror"
)

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	fmt.Println("GET params were:", r.URL.Query())

	f := &sub_s.UserListFilter{
		// PageSize:  10,
		// LastID:    "",
		SortField:       "_id",
		ExcludeArchived: true,
	}

	// Apply search text if it exists in url parameter.
	searchKeyword := r.URL.Query().Get("search")
	if searchKeyword != "" {
		f.SearchText = searchKeyword
	}

	// Apply filters it exists in url parameter.
	firstName := r.URL.Query().Get("first_name")
	if firstName != "" {
		f.FirstName = firstName
	}
	lastName := r.URL.Query().Get("first_name")
	if lastName != "" {
		f.LastName = lastName
	}
	email := r.URL.Query().Get("email")
	if email != "" {
		f.Email = email
	}
	phone := r.URL.Query().Get("phone")
	if phone != "" {
		f.Phone = phone
	}

	// Perform our database operation.
	m, err := h.Controller.ListByFilter(ctx, f)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	MarshalListResponse(m, w)
}

func MarshalListResponse(res *sub_s.UserListResult, w http.ResponseWriter) {
	if err := json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
