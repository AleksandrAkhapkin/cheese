package handlers

import (
	"encoding/json"
	"github.com/lilekov-studio/supreme-cheese/internal/supreme-cheese/types"
	"net/http"
)

func (h *Handlers) AddCheeseName(w http.ResponseWriter, r *http.Request) {

	addName := types.AddName{}
	if err := json.NewDecoder(r.Body).Decode(&addName); err != nil {
		apiErrorEncode(w, err)
		return
	}

	err := h.srv.AddCheeseName(addName.Name)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}
}

func (h *Handlers) WriteBill(w http.ResponseWriter, r *http.Request) {

	bill := types.HandWriteBill{}
	if err := json.NewDecoder(r.Body).Decode(&bill); err != nil {
		apiErrorEncode(w, err)
		return
	}

	if err := h.srv.HandWriteBill(&bill); err != nil {
		apiErrorEncode(w, err)
		return
	}
}

func (h *Handlers) InvalidBill(w http.ResponseWriter, r *http.Request) {

	invalidbill := types.InvalidBill{}
	if err := json.NewDecoder(r.Body).Decode(&invalidbill); err != nil {
		apiErrorEncode(w, err)
		return
	}

	err := h.srv.InvalidBill(&invalidbill)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}
}

func (h *Handlers) SendCashback(w http.ResponseWriter, r *http.Request) {
	send := types.SendCashback{}
	if err := json.NewDecoder(r.Body).Decode(&send); err != nil {
		apiErrorEncode(w, err)
		return
	}

	if err := h.srv.SendCashback(&send); err != nil {
		apiErrorEncode(w, err)
		return
	}

	apiResponseEncoder(w, "Саня, все прошло заебись")
}
