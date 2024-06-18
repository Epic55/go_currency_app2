package httphandler

import (
	"Epic55/go_currency_app2/internal/models"
	"Epic55/go_currency_app2/internal/repository"
	"Epic55/go_currency_app2/internal/service"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type Handler struct {
	R    *repository.Repository
	Cnfg *models.Config
}

func NewHandler(repo *repository.Repository, config *models.Config) *Handler {
	if repo == nil {
		fmt.Println("Failed to initialize the repo")
	}

	return &Handler{
		R:    repo,
		Cnfg: config,
	}
}

func DateFormat(date string) (string, error) {
	parsedDate, err := time.Parse("02.01.2006", date)
	if err != nil {
		return "", err
	}
	formattedDate := parsedDate.Format("2006-01-02")
	return formattedDate, nil
}

func (h *Handler) RespondWithError(w http.ResponseWriter, status int, errorMsg string, err error) {
	http.Error(w, errorMsg, status)
	fmt.Println(errorMsg+": ", err)
}

func (h *Handler) SaveCurrencyHandler(w http.ResponseWriter, r *http.Request, ctx context.Context) {
	vars := mux.Vars(r)
	date := vars["date"]

	formattedDate, err := DateFormat(date)
	if err != nil {
		h.RespondWithError(w, http.StatusBadRequest, "Failed to parse and format the date", err)
		return
	}
	var service = service.NewService()

	go h.R.InsertData(*service.GetData1(ctx, date, h.Cnfg.APIURL), formattedDate)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"success": true}`))
	fmt.Println("Success: true")
}

func (h *Handler) GetCurrencyHandler(w http.ResponseWriter, r *http.Request, ctx context.Context) {
	vars := mux.Vars(r)
	date := vars["date"]
	code := vars["code"]

	formattedDate, err := DateFormat(date)
	if err != nil {
		h.RespondWithError(w, http.StatusBadRequest, "Failed to parse the date", err)
		return
	}

	data, err := h.R.GetData(ctx, formattedDate, code)
	if err != nil {
		h.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve date 1", err)
		return
	}
	fmt.Println("Data was showed")
	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func (h *Handler) DeleteCurrencyHandler(w http.ResponseWriter, r *http.Request, ctx context.Context) {
	vars := mux.Vars(r)
	date := vars["date"]
	code := vars["code"]

	formattedDate, err := DateFormat(date)
	if err != nil {
		h.RespondWithError(w, http.StatusBadRequest, "Failed to parse the date", err)
		return
	}

	data, err := h.R.DeleteData(ctx, formattedDate, code)
	if err != nil {
		h.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve date 2", err)
		return
	}
	fmt.Println("Data was removed")
	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func (h *Handler) StartScheduler(ctx context.Context) {
	date := time.Now().Format("02.01.2006")
	formattedDate, err := DateFormat(date)
	if err != nil {
		fmt.Println("Cannot parse the Data")
	}
	h.R.HourTick(date, formattedDate, ctx, h.Cnfg.APIURL)
}
