package handlers

import (
	"fmt"
	"github.com/lilekov-studio/supreme-cheese/pkg/infrastruct"
	"github.com/lilekov-studio/supreme-cheese/pkg/logger"
	"net/http"
)

func (h *Handlers) VerifyPhone(w http.ResponseWriter, r *http.Request) {

	claims, err := infrastruct.GetClaimsByRequest(r, h.secretKey)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}

	phone, err := h.srv.VerifyPhone(claims.UserID)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}

	apiResponseEncoder(w, phone)
}

func (h *Handlers) UploadBill(w http.ResponseWriter, r *http.Request) {

	claims, err := infrastruct.GetClaimsByRequest(r, h.secretKey)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}

	//Закомичено, так как акция завершилась

	//fileBody, _, err := r.FormFile("file")
	//if err != nil {
	//	apiErrorEncode(w, infrastruct.ErrorBadRequest)
	//	return
	//}
	//
	//perekBool, err := strconv.ParseBool(r.FormValue("perek_bool"))
	//if err != nil {
	//	apiErrorEncode(w, infrastruct.ErrorBadRequest)
	//	return
	//}
	//perecCard := r.FormValue("perek_card")
	//
	//file := types.UploadFile{
	//	UserID:          claims.UserID,
	//	Body:            fileBody,
	//	PerecrestokBool: perekBool,
	//	PerecrestokCard: perecCard,
	//}
	//
	//if err = h.srv.UploadFile(&file); err != nil {
	//	apiErrorEncode(w, err)
	//	return
	//}

	logger.LogInfo(fmt.Sprintf("попытка загрузить чек, ID %d", claims.UserID))
	return

}
