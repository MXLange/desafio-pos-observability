package weatherapi

import (
	"encoding/json"
	"net/http"

	"github.com/MXLange/desafio-pos-observability/internal/zipcode"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

type requestPayload struct {
	CEP string `json:"cep"`
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var payload requestPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		zipcode.WriteError(w, http.StatusUnprocessableEntity, zipcode.ErrInvalidZipcode)
		return
	}

	cep, err := zipcode.Validate(payload.CEP)
	if err != nil {
		zipcode.WriteError(w, http.StatusUnprocessableEntity, zipcode.ErrInvalidZipcode)
		return
	}

	result, err := h.service.GetWeatherByCEP(r.Context(), cep)
	if err != nil {
		statusCode, message := zipcode.StatusFromError(err)
		zipcode.WriteError(w, statusCode, message)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(result)
}
