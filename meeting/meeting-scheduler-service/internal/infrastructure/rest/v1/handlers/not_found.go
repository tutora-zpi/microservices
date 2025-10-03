package handlers

import (
	"net/http"
)

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/api/v1/docs", http.StatusSeeOther)
}
