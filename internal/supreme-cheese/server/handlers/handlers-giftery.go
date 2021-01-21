package handlers

import (
	"encoding/json"
	"github.com/lilekov-studio/supreme-cheese/internal/supreme-cheese/types"
	"github.com/lilekov-studio/supreme-cheese/pkg/infrastruct"
	"net/http"
)

func (h *Handlers) ChoiceCertificate(w http.ResponseWriter, r *http.Request) {

	claims, err := infrastruct.GetClaimsByRequest(r, h.secretKey)
	if err != nil {
		apiErrorEncode(w, err)
		return
	}

	choiceCertificate := types.ChoiceCertificate{}
	if err := json.NewDecoder(r.Body).Decode(&choiceCertificate); err != nil {
		apiErrorEncode(w, err)
		return
	}

	if err = h.srv.ChoiceCertificate(claims.UserID, choiceCertificate.ProductID); err != nil {
		apiErrorEncode(w, err)
		return
	}

}
