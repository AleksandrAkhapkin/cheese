package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/lilekov-studio/supreme-cheese/internal/supreme-cheese/types"
	"github.com/lilekov-studio/supreme-cheese/pkg/infrastruct"
	"net/http"
)

func (h *Handlers) RegisterUser(w http.ResponseWriter, r *http.Request) {

	user := types.UserRegister{}
	var err error

	if err = json.NewDecoder(r.Body).Decode(&user); err != nil {
		apiErrorEncode(w, infrastruct.ErrorBadRequest)
		return
	}

	if user.Email == "" || user.FirstName == "" || user.LastName == "" || user.Phone == "" {
		apiErrorEncode(w, infrastruct.ErrorBadRequest)
		return
	}

	//ip := r.Header.Get("X-Real-IP")
	//	if ip == "" {
	ip := r.RemoteAddr
	//	}
	user.City = "Default"

	err = h.srv.RegisterUser(&user, ip)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}
}

func (h *Handlers) RecoverPassword(w http.ResponseWriter, r *http.Request) {

	rec := types.RecoverPass{}
	var err error

	if err = json.NewDecoder(r.Body).Decode(&rec); err != nil {
		apiErrorEncode(w, infrastruct.ErrorBadRequest)
		return
	}

	if rec.Email == "" {
		apiErrorEncode(w, infrastruct.ErrorBadRequest)
		return
	}

	err = h.srv.RecoverPassword(&rec)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}
}

func (h *Handlers) RegisterPhone(w http.ResponseWriter, r *http.Request) {

	phone := types.RegisterPhone{}
	var err error

	if err = json.NewDecoder(r.Body).Decode(&phone); err != nil {
		apiErrorEncode(w, infrastruct.ErrorBadRequest)
		return
	}

	claims, err := infrastruct.GetClaimsByRequest(r, h.secretKey)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}

	if phone.Phone == "" {
		apiErrorEncode(w, infrastruct.ErrorBadRequest)
		return
	}

	err = h.srv.RegisterPhone(&phone, claims.UserID)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}
}

func (h *Handlers) ConfirmationPhone(w http.ResponseWriter, r *http.Request) {

	user := types.ConfirmationPhone{}
	var err error

	if err = json.NewDecoder(r.Body).Decode(&user); err != nil {
		apiErrorEncode(w, infrastruct.ErrorBadRequest)
		return
	}

	if err := h.srv.ConfirmationPhone(&user); err != nil {
		apiErrorEncode(w, err)
		return
	}
}

func (h *Handlers) Auth(w http.ResponseWriter, r *http.Request) {

	auth := types.Authorize{}
	var err error

	if err = json.NewDecoder(r.Body).Decode(&auth); err != nil {
		apiErrorEncode(w, infrastruct.ErrorBadRequest)
		return
	}

	token, err := h.srv.Authorize(&auth)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}

	apiResponseEncoder(w, token)
}

func (h *Handlers) AuthByToken(w http.ResponseWriter, r *http.Request) {

	authToken := types.Token{}
	var err error

	if err = json.NewDecoder(r.Body).Decode(&authToken); err != nil {
		apiErrorEncode(w, infrastruct.ErrorBadRequest)
		return
	}

	token, err := h.srv.AuthByToken(&authToken)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}

	apiResponseEncoder(w, token)
}

func (h *Handlers) GetCabinetByClaims(w http.ResponseWriter, r *http.Request) {

	claims, err := infrastruct.GetClaimsByRequest(r, h.secretKey)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}
	user, err := h.srv.GetUserByID(claims.UserID)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}

	apiResponseEncoder(w, user)
}

func (h *Handlers) PutCabinetByClaims(w http.ResponseWriter, r *http.Request) {

	user := types.PutCabinet{}
	var err error

	claims, err := infrastruct.GetClaimsByRequest(r, h.secretKey)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}
	user.ID = claims.UserID

	if err = json.NewDecoder(r.Body).Decode(&user); err != nil {
		apiErrorEncode(w, infrastruct.ErrorBadRequest)
		return
	}

	if err = h.srv.PutCabinet(&user); err != nil {
		apiErrorEncode(w, err)
		return
	}
}

func (h *Handlers) UnsubscribeMail(w http.ResponseWriter, r *http.Request) {

	if err := h.srv.UnsubscribeMail(mux.Vars(r)["email"]); err != nil {
		apiErrorEncode(w, err)
		return
	}

	apiResponseEncoder(w, "Вы успешно отписались от рассылки!")
}
