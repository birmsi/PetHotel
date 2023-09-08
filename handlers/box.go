package handlers

import (
	"PetHotel/services"
	"net/http"
)

type BoxHandler struct {
	//TODO:logger para logar as cenas
	Service services.BoxService
}

func NewBoxHandler(boxService services.BoxService) BoxHandler {
	return BoxHandler{Service: boxService}
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
