package handlers

import (
	"PetHotel/models"
	"PetHotel/services"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
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

type CreateBoxViewResponse struct {
	ErrorMessage string
	BoxSizes     []string
	Box          models.Box
}

func (bh BoxHandler) CreateBoxView(w http.ResponseWriter, r *http.Request) {

	response := CreateBoxViewResponse{
		BoxSizes: []string{string(models.SmallSize), string(models.MediumSize), string(models.LargeSize)},
		Box:      models.Box{},
	}

	bh.CreateBoxViewParser(w, r, response)

}

func (bh BoxHandler) CreateBoxPost(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	errors := make([]string, 0)

	num := r.PostForm.Get("number")
	size := r.PostForm.Get("size")

	//TODO: why not just return the error when it happens instead of adding to array?
	if len(num) == 0 {
		errors = append(errors, "Missing number")
	}
	if len(size) == 0 {
		errors = append(errors, "Missing size")
	}
	if len(errors) > 0 {
		bh.CreateBoxViewParser(w, r, CreateBoxViewResponse{
			BoxSizes:     []string{string(models.SmallSize), string(models.MediumSize), string(models.LargeSize)},
			Box:          models.Box{Number: 0, Size: models.Size(size)},
			ErrorMessage: strings.Join(errors, "."),
		})
	}

	newNumber, err := strconv.Atoi(num)
	if err != nil {
		errors = append(errors, "Invalid number")
		bh.CreateBoxViewParser(w, r, CreateBoxViewResponse{
			BoxSizes:     []string{string(models.SmallSize), string(models.MediumSize), string(models.LargeSize)},
			Box:          models.Box{Number: newNumber, Size: models.Size(size)},
			ErrorMessage: strings.Join(errors, "."),
		})
	}

	id, err := bh.Service.CreateBox(models.Box{Number: newNumber, Size: models.Size(size)})
	if err != nil {
		errors = append(errors, err.Error())
		bh.CreateBoxViewParser(w, r, CreateBoxViewResponse{
			BoxSizes:     []string{string(models.SmallSize), string(models.MediumSize), string(models.LargeSize)},
			Box:          models.Box{Number: newNumber, Size: models.Size(size)},
			ErrorMessage: strings.Join(errors, "."),
		})
	}

	http.Redirect(w, r, fmt.Sprintf("/box/%d/view", id), http.StatusSeeOther)

}

func (bh BoxHandler) CreateBoxViewParser(
	w http.ResponseWriter,
	r *http.Request,
	response CreateBoxViewResponse) {

	t, err := template.ParseFiles(
		"./ui/html/base.html",
		"./ui/html/partials/navigation.html",
		"./ui/html/pages/box/create.html",
	)

	if err != nil {
		bh.slogger.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	if err = t.ExecuteTemplate(w, "base", response); err != nil {
		bh.slogger.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

type GetBox struct {
	Box models.Box
}

// GetBoxView Does not allow for edit - only view mode :)
func (bh BoxHandler) GetBoxView(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "boxID")
	boxID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "BAD REQUEST", http.StatusBadRequest)
		return
	}

	box, err := bh.Service.GetBox(boxID)
	if err != nil {
		bh.slogger.Error(err.Error())
		http.Redirect(w, r, "/box", http.StatusSeeOther)
		return
	}

	//TODO: eventually add the bookings of and for this box
	// add the availabilities

	t, err := template.ParseFiles(
		"./ui/html/base.html",
		"./ui/html/partials/navigation.html",
		"./ui/html/pages/box/view.html",
	)

	if err != nil {
		bh.slogger.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	if err = t.ExecuteTemplate(w, "base", box); err != nil {
		bh.slogger.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
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
