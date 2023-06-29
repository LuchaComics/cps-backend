package attachment

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	sub_c "github.com/LuchaComics/cps-backend/app/attachment/controller"
	sub_s "github.com/LuchaComics/cps-backend/app/attachment/datastore"
	"github.com/LuchaComics/cps-backend/utils/httperror"
)

func UnmarshalUpdateRequest(ctx context.Context, r *http.Request) (*sub_c.AttachmentUpdateRequestIDO, error) {
	// Initialize our array which will store all the results from the remote server.
	var requestData sub_c.AttachmentUpdateRequestIDO

	defer r.Body.Close()

	// Read the JSON string and convert it into our golang stuct else we need
	// to send a `400 Bad Request` errror message back to the client,
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		log.Println(err)
		return nil, httperror.NewForSingleField(http.StatusBadRequest, "non_field_error", "payload structure is wrong")
	}

	return &requestData, nil
}

func (h *Handler) UpdateByID(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()

	data, err := UnmarshalUpdateRequest(ctx, r)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	attachment, err := h.Controller.UpdateByID(ctx, data)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	MarshalUpdateResponse(attachment, w)
}

func MarshalUpdateResponse(res *sub_s.Attachment, w http.ResponseWriter) {
	if err := json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
