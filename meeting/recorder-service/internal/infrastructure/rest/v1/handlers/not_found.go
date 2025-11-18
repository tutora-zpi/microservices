package handlers

import (
	"net/http"
)

// NotFound implements Handler.
func (h *handlerImpl) NotFound(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/api/v1/docs", http.StatusSeeOther)
}
