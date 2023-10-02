package handlers

import (
	"PetHotel/models"
	"PetHotel/responses"
	"PetHotel/services"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
)

type BoxHandler struct {
	Service        services.BoxService
	BookingService services.BookingService
	slogger        *slog.Logger
}

func NewBoxHandler(boxService services.BoxService, slogger *slog.Logger, bookingService services.BookingService) BoxHandler {
	return BoxHandler{Service: boxService, slogger: slogger, BookingService: bookingService}
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

type BoxResponse struct {
	ID            int
	Number        int
	Size          string
	Availabilites []*responses.AvailabilityResponse
	Bookings      []*models.Booking
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

	month := 24 * time.Hour * 30
	startOfMonth := time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, time.UTC)

	availabilities, err := bh.Service.GetAvailabilities(boxID, startOfMonth, startOfMonth.Add(month))
	if err != nil {
		bh.slogger.Error(err.Error())
		return
	}

	availabilities[0].StartTime.Format(time.RFC3339)

	bookings, err := bh.BookingService.GetBookings(boxID, startOfMonth, time.Now().Add(month))
	if err != nil {
		bh.slogger.Error(err.Error())
		return
	}

	t, err := template.ParseFiles(
		"./ui/html/base.html",
		"./ui/html/partials/navigation.html",
		"./ui/html/pages/box/view.html",
	)

	if err != nil {
		bh.slogger.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
	availabilitiesReponse := make([]*responses.AvailabilityResponse, 0, len(availabilities))

	for _, availability := range availabilities {
		availabilitiesReponse = append(availabilitiesReponse, &responses.AvailabilityResponse{
			ID:        availability.ID,
			BoxID:     availability.BoxID,
			StartTime: availability.StartTime.Format(time.RFC3339),
			EndTime:   availability.EndTime.Format(time.RFC3339),
			Price:     availability.Price,
		})
	}

	boxResponse := BoxResponse{
		ID:            boxID,
		Number:        box.Number,
		Size:          string(box.Size),
		Availabilites: availabilitiesReponse,
		Bookings:      bookings,
	}

	if err = t.ExecuteTemplate(w, "base", boxResponse); err != nil {
		bh.slogger.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (bh BoxHandler) GetBoxUpdateView(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("GetBoxUpdateView"))
	http.Error(w, "Not implemented yet", http.StatusNotImplemented)
}

func (bh BoxHandler) GetBoxUpdatePut(w http.ResponseWriter, r *http.Request) {

	boxID := chi.URLParam(r, "boxID")

	id, err := strconv.Atoi(boxID)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	r.ParseForm()

	// num := r.PostForm.Get("number")
	// size := r.PostForm.Get("size")

	start_times, err := fromArrayOfStringToTime(r.PostForm["start_time"])
	if err != nil || start_times == nil {
		fmt.Printf("erro!")
	}
	end_times, err := fromArrayOfStringToTime(r.PostForm["end_time"])
	if err != nil || end_times == nil {
		fmt.Printf("erro!")
	}
	prices, err := fromArrayOfStringToFloat64(r.PostForm["price"])
	if err != nil || prices == nil {
		fmt.Printf("erro!")
	}

	if len(start_times) != len(end_times) || len(start_times) != len(prices) {
		fmt.Printf("erro!")
	}

	availabilities := make([]models.Availability, 0)
	for i := range start_times {
		availabilities = append(availabilities, models.Availability{
			StartTime: *start_times[i],
			EndTime:   *end_times[i],
			Price:     *prices[i],
			BoxID:     id,
		})
	}

	isValid := CheckValidDates(availabilities)

	if !isValid {
		fmt.Println("Issue with dates")
		return
	}

	futureAvailabilities, err := bh.Service.GetFutureAvailabilities(id)
	if err != nil {
		fmt.Println("Issue with dates")
		return
	}

	hasOverlap := HasOverlap(availabilities, futureAvailabilities)
	if hasOverlap {
		fmt.Println("Issue with future availability dates")
		return
	}
	futureBookings, err := bh.BookingService.GetFutureBookings(id)
	if err != nil {
		fmt.Println("Issue with dates")
		return
	}

	hasOverlap = OverlapsFutureBookings(availabilities, futureBookings)
	if hasOverlap {
		fmt.Println("Issue with future bookings dates")
		return
	}

	err = bh.Service.AddAvailabilities(availabilities)
	if err != nil {
		fmt.Println("Issue with future bookings dates", err.Error())
		return
	}

	http.Error(w, "Not implemented yet", http.StatusNotImplemented)
}

func OverlapsFutureBookings(availabilities []models.Availability, futureBookings []*models.Booking) bool {
	for _, avail1 := range availabilities {
		for _, avail2 := range futureBookings {
			if !(avail1.EndTime.Before(avail2.CheckIn) || avail2.CheckOut.Before(avail1.StartTime)) {
				return true
			}
		}
	}
	return false
}

func fromArrayOfStringToTime(toConvert []string) ([]*time.Time, error) {
	times := make([]*time.Time, 0)

	for _, t := range toConvert {
		dt, err := time.Parse("2006-01-02", t)
		if err != nil {
			return nil, err
		}
		times = append(times, &dt)
	}

	return times, nil
}
func fromArrayOfStringToFloat64(toConvert []string) ([]*float64, error) {
	prices := make([]*float64, 0)

	for _, s := range toConvert {
		floatValue, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return nil, err
		}

		prices = append(prices, &floatValue)
	}

	return prices, nil
}

func (bh BoxHandler) BoxDelete(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("BoxDelete"))
	http.Error(w, "Not implemented yet", http.StatusNotImplemented)
}

// CheckValidDates checks if the start dates are at least on the current date or in the future.
func CheckValidDates(availabilities []models.Availability) bool {

	isValid := CheckDatesCorrectness(availabilities)

	overlaps := CheckOverlap(availabilities)

	isFuture := CheckIfDatesInPresentOrFuture(availabilities)

	return isValid && !overlaps && isFuture
}

func CheckOverlap(availabilities []models.Availability) bool {
	n := len(availabilities)

	// Iterate through each pair of availabilities and check for overlap
	for i := 0; i < n-1; i++ {
		for j := i + 1; j < n; j++ {
			// Check if the end time of one availability is before the start time of the other
			if availabilities[i].EndTime.Before(availabilities[j].StartTime) ||
				availabilities[j].EndTime.Before(availabilities[i].StartTime) {
				continue // No overlap, check the next pair
			} else {
				// Overlap detected
				return true
			}
		}
	}
	return false
}

func CheckIfDatesInPresentOrFuture(availabilities []models.Availability) bool {
	currentDate := time.Now()

	for _, avail := range availabilities {
		requestedDate := time.Date(avail.StartTime.Year(), avail.StartTime.Month(), avail.StartTime.Day(), 0, 0, 0, 0, time.UTC)
		date2 := time.Date(currentDate.Year(), currentDate.Month(), currentDate.Day(), 0, 0, 0, 0, time.UTC)

		if requestedDate.Before(date2) {
			return false
		}
	}
	return true
}

func CheckDatesCorrectness(availabilities []models.Availability) bool {

	for _, availability := range availabilities {
		if availability.EndTime.Before(availability.StartTime) {
			return false
		}
	}
	return true
}

func HasOverlap(availabilities []models.Availability, futureAvailabilities []*models.Availability) bool {
	for _, avail1 := range availabilities {
		for _, avail2 := range futureAvailabilities {
			if !(avail1.EndTime.Before(avail2.StartTime) || avail2.EndTime.Before(avail1.StartTime)) {
				return true
			}
		}
	}
	return false
}
