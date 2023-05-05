package gateway

import (
	"net/http"
)

// Version returns the server version. Developers note, to see result you can run in your terminal `curl http://localhost:8000/api/v1/version`.
func (h *Handler) Version(w http.ResponseWriter, r *http.Request) {
	result := "v1.0"
	w.Write([]byte(result))
}
