package handlers

import (
	"PetHotel/services"
	"log/slog"
	"net/http"
)

type BoxHandler struct {
	Service services.BoxService
	slogger *slog.Logger
}

func NewBoxHandler(boxService services.BoxService, slogger *slog.Logger) BoxHandler {
	return BoxHandler{Service: boxService, slogger: slogger}
}
func (bh BoxHandler) GetBoxesView(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("All boxes"))
}
func (bh BoxHandler) CreateBoxView(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("create box view"))
}

func (bh BoxHandler) CreateBoxPost(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("create box POST"))
	http.Error(w, "Not implemented yet", http.StatusNotImplemented)
}

func (bh BoxHandler) GetBoxUpdateView(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("GetBoxUpdateView"))
	http.Error(w, "Not implemented yet", http.StatusNotImplemented)
}

func (bh BoxHandler) GetBoxUpdatePut(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("GetBoxUpdatePut"))
	http.Error(w, "Not implemented yet", http.StatusNotImplemented)
}

func (bh BoxHandler) BoxDelete(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("BoxDelete"))
	http.Error(w, "Not implemented yet", http.StatusNotImplemented)
}
