package organization

import (
	"encoding/json"
	"net/http"

	sub_s "github.com/LuchaComics/cps-backend/app/organization/datastore"
	"github.com/LuchaComics/cps-backend/utils/httperror"
)

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	f := &sub_s.OrganizationListFilter{
		PageSize:        10,
		LastID:          "",
		SortField:       "_id",
		ExcludeArchived: true,
	}

	m, err := h.Controller.ListByFilter(ctx, f)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	MarshalListResponse(m, w)
}

func MarshalListResponse(res *sub_s.OrganizationListResult, w http.ResponseWriter) {
	if err := json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
