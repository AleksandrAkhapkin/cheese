package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/lilekov-studio/supreme-cheese/internal/supreme-cheese/service"
	"github.com/lilekov-studio/supreme-cheese/internal/supreme-cheese/types"
	"github.com/lilekov-studio/supreme-cheese/internal/supreme-cheese/types/config"
	"github.com/lilekov-studio/supreme-cheese/pkg/infrastruct"
	"github.com/lilekov-studio/supreme-cheese/pkg/logger"
	"github.com/pkg/errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type Handlers struct {
	srv       *service.Service
	secretKey string
}

func NewHandlers(srv *service.Service, cnf *config.Config) *Handlers {

	return &Handlers{
		srv:       srv,
		secretKey: cnf.SecretKeyJWT,
	}
}

func (h *Handlers) Ping(w http.ResponseWriter, _ *http.Request) {
	_, _ = w.Write([]byte("pong v0912 1841"))
}

func (h *Handlers) OS(w http.ResponseWriter, r *http.Request) {

	os := types.OS{UserID: 0}

	claims, err := infrastruct.GetClaimsByRequest(r, h.secretKey)
	if err == nil {
		os.UserID = claims.UserID
	}

	if err = json.NewDecoder(r.Body).Decode(&os); err != nil {
		apiErrorEncode(w, infrastruct.ErrorBadRequest)
		return
	}

	err = h.srv.OS(&os)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}
}

func (h *Handlers) GetWinners(w http.ResponseWriter, _ *http.Request) {

	winners, err := h.srv.GetWinners()
	if err != nil {
		apiErrorEncode(w, err)
		return
	}

	apiResponseEncoder(w, winners)
}

func (h *Handlers) GetMyBills(w http.ResponseWriter, r *http.Request) {

	claims, err := infrastruct.GetClaimsByRequest(r, h.secretKey)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}

	bills, err := h.srv.GetMyBills(claims.UserID)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}

	apiResponseEncoder(w, bills)
}

func (h *Handlers) GetMyPrize(w http.ResponseWriter, r *http.Request) {

	claims, err := infrastruct.GetClaimsByRequest(r, h.secretKey)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}

	prizes, err := h.srv.GetMyPrizes(claims.UserID)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}

	apiResponseEncoder(w, prizes)
}

func apiErrorEncode(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if customError, ok := err.(*infrastruct.CustomError); ok {
		w.WriteHeader(customError.Code)
	}

	result := struct {
		Err string `json:"error"`
	}{
		Err: err.Error(),
	}

	if err = json.NewEncoder(w).Encode(result); err != nil {
		logger.LogError(errors.Wrap(err, "err with Encode in GetMyPrize"))
	}
}

func apiResponseEncoder(w http.ResponseWriter, res interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		logger.LogError(errors.Wrap(err, "err with apiResponseEncoder"))
	}
}

func (h *Handlers) MakeOrderForToken(w http.ResponseWriter, r *http.Request) {

	//отправляю запрос через постман
	errRed := r.FormValue("error")
	errRedDesc := r.FormValue("error_description")
	if errRed != "" {
		logger.LogError(fmt.Errorf("err with auth нщщьщтун, err:%s\n%s", errRed, errRedDesc))
	}
	code := r.FormValue("error")
	logger.LogInfo(code)

	//отправляю запрос с кодом что бы получить токен
	data := url.Values{}
	data.Set("code", code)
	data.Set("client_id", "38FCD146700D7AB0BACB46EFD24184CBEFA45A43509ABBD6B835A365D091F3BB")
	data.Set("grant_type", "authorization_code")
	data.Set("redirect_uri", "https://supreme-cheese.ru/yoomoney/getkey/token")
	r, err := http.NewRequest("POST", "https://yoomoney.ru/oauth/token", strings.NewReader(data.Encode()))
	if err != nil {
		errors.Wrap(err, "err with NewRequest in payment")
	}
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	_, _ = http.DefaultClient.Do(r)

}

func (h *Handlers) ReqToken(w http.ResponseWriter, r *http.Request) {

	type Req struct {
		Access_token string `json:"access_token"`
		Error        string `json:"error"`
	}
	req := Req{}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with NewDecoder in MakeOrderForToken"))
		return
	}
	logger.LogInfo(req.Access_token)
}
