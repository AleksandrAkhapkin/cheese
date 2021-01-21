package service

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/lilekov-studio/supreme-cheese/internal/supreme-cheese/service/mail"
	"github.com/lilekov-studio/supreme-cheese/internal/supreme-cheese/types"
	"github.com/lilekov-studio/supreme-cheese/pkg/infrastruct"
	"github.com/lilekov-studio/supreme-cheese/pkg/logger"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"net/http"
)

func (s *Service) ChoiceCertificate(userID int, choiceCertificate int) error {

	billID, err := s.p.GetWinnerCertificateByUserID(userID)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with s.p.GetWinnerCertificateByUserID in ChoiceCertificate"))
		return infrastruct.ErrorInternalServerError
	}
	err = s.p.UpdateWinnerCertStatus(billID, "Отправляется")
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with s.p.UpdateWinnerCertStatus in ChoiceCertificate"))
		return infrastruct.ErrorInternalServerError
	}

	prizeShopName := ""
	switch choiceCertificate {
	case 14073:
		prizeShopName = "La Corvette"
	case 14578:
		prizeShopName = "Яндекс.Афиша"
	case 11610:
		prizeShopName = "Перекресток"
	case 13370:
		prizeShopName = "Rendez-Vous"
	case 800:
		prizeShopName = "Decanter"
	case 13882:
		prizeShopName = "Ticketland"
	case 13515:
		prizeShopName = "Togas"
	case 11956:
		prizeShopName = "Lamoda"
	default:
		logger.LogInfo("Ошибка в номере сертификата")
		return infrastruct.ErrorInternalServerError
	}

	testModeForHash := "" //"&testmode=1" //todo PROD заменить на пустую строку
	testMode := ""        //"%26testmode=1"      //todo PROD заменить на пустую строку

	shaString := fmt.Sprintf("makeOrderproduct_id=%d&face=%d&delivery_type=%s%s%s", choiceCertificate, types.GifteryNominal, "download", testModeForHash, s.secretGiftery)
	h := sha256.New()
	h.Write([]byte(shaString))
	sig := fmt.Sprintf("sig=%x", h.Sum(nil))
	data := fmt.Sprintf("data=product_id=%d%%26face=%d%%26delivery_type=%s%s", choiceCertificate, types.GifteryNominal, "download", testMode)
	gifteryURL := fmt.Sprintf("https://ssl-api.giftery.ru/?id=15211&cmd=makeOrder&%s&%s&out=json", data, sig)

	//логируем реквест
	reqID, err := s.p.AddLogRequestGiftery(gifteryURL, userID, billID, types.GifteryTypeMethodMakeOrder)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with AddLogCheckBillResp in ChoiceCertificate for logger"))
		return infrastruct.ErrorInternalServerError
	}

	req, err := http.NewRequest("GET", gifteryURL, nil)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with http.NewReques in ChoiceCertificate"))
		return infrastruct.ErrorInternalServerError
	}
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with client.Do(req) in ChoiceCertificate"))
		return infrastruct.ErrorInternalServerError
	}
	defer res.Body.Close()

	//логируем респонс
	respLogerBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with read body in ChoiceCertificate for logger"))
	}

	if err := s.p.AddLogResponseGiftery(string(respLogerBytes), userID, billID, reqID, types.GifteryTypeMethodMakeOrder); err != nil {
		logger.LogError(errors.Wrap(err, "err with AddLogCheckBillResp in ChoiceCertificate for logger"))
	}

	gifteryResp := types.Giftery{}
	if err := json.Unmarshal(respLogerBytes, &gifteryResp); err != nil {
		logger.LogError(errors.Wrap(err, "err with json.Unmarshal in ChoiceCertificate"))
	}

	if gifteryResp.Status == "error" {
		logger.LogError(errors.New("error with GIFTERY - code: " + strconv.Itoa(gifteryResp.GifteryError.Error.Code) + " text: " + gifteryResp.GifteryError.Error.TextError))
		err = s.p.UpdateWinnerCertStatus(billID, fmt.Sprintf("Заказ на сертификат %s создан", prizeShopName))
		if err != nil {
			logger.LogError(errors.Wrap(err, "err with s.p.UpdateWinnerCertStatus in ChoiceCertificate"))
			return infrastruct.ErrorInternalServerError
		}

		return nil
	}

	err = s.p.UpdateWinnerCertStatus(billID, types.PrizeStatusAccept)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with s.p.UpdateWinnerCertStatus in ChoiceCertificate"))
		return infrastruct.ErrorInternalServerError
	}

	go s.DownloadCertificate(gifteryResp.GifteryOK.Data.ID, userID, billID)

	return nil
}

func (s *Service) DownloadCertificate(idSert, userID, billID int) {
	time.Sleep(time.Second * 120)

	shaString := fmt.Sprintf("getCertificate{\"queue_id\":%d}%s", idSert, s.secretGiftery)
	h := sha256.New()
	h.Write([]byte(shaString))
	sig := fmt.Sprintf("sig=%x", h.Sum(nil))
	data := fmt.Sprintf("&data=%%7B%%22queue_id%%22%%3A%d%%7D", idSert)
	gifteryURL := fmt.Sprintf("https://ssl-api.giftery.ru/?cmd=getCertificate&id=15211&in=json&out=json%s&%s", data, sig)

	//логируем реквест
	reqID, err := s.p.AddLogRequestGiftery(gifteryURL, userID, billID, types.GifteryTypeMethodGetCertificate)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with AddLogCheckBillResp in DownloadCertificate for logger"))
		return
	}

	req, err := http.NewRequest("GET", gifteryURL, nil)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with http.NewReques in DownloadCertificate"))
		return
	}
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with client.Do(req) in DownloadCertificate"))
		return
	}
	defer res.Body.Close()

	//логируем респонс
	gifteryResp := types.Giftery{}
	if err := json.NewDecoder(res.Body).Decode(&gifteryResp); err != nil {
		logger.LogError(errors.Wrap(err, "err with json.NewDecoder in DownloadCertificate"))
		return
	}

	if err := s.p.AddLogResponseGiftery(fmt.Sprintf("%+v, %+v", gifteryResp.Status, gifteryResp.Error), userID, billID, reqID, types.GifteryTypeMethodGetCertificate); err != nil {
		logger.LogError(errors.Wrap(err, "err with AddLogResponseGiftery in DownloadCertificate for logger"))
	}

	if gifteryResp.Status == "error" {
		logger.LogError(errors.New("error with GIFTERY - code: " + strconv.Itoa(gifteryResp.GifteryError.Error.Code) + " text: " + gifteryResp.GifteryError.Error.TextError))
		return
	}

	dec, err := base64.StdEncoding.DecodeString(gifteryResp.GifteryOK.Data.Certificate)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with base64.StdEncoding.DecodeString in DownloadCertificate"))
		return
	}

	if err := os.Mkdir(filepath.Join(s.pathForSert, fmt.Sprintf("%d", userID)), 0777); err != nil {
		if !os.IsExist(err) {
			logger.LogError(errors.Wrap(err, "err with os.Mkdir in DownloadCertificate"))
			return
		}
	}

	path := filepath.Join(s.pathForSert, fmt.Sprintf("/%d", userID))
	file, err := os.Create(fmt.Sprintf("%s/Certificate.pdf", path))
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with os.Create in DownloadCertificate"))
		return
	}
	defer file.Close()

	if _, err := file.Write(dec); err != nil {
		logger.LogError(errors.Wrap(err, "err with file.Write in DownloadCertificate"))
		return
	}
	if err := file.Sync(); err != nil {
		logger.LogError(errors.Wrap(err, "err with file.Sync in DownloadCertificate"))
		return
	}

	user, err := s.p.GetUserByID(userID)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with s.p.GetUserByID in DownloadCertificate"))
		return
	}

	r := mail.NewRequest([]string{user.Email}, s.email)
	if err := r.SendCert("Подарок от SUPREME-CHEESE!", user.ID, s.pathForSert); err != nil {
		logger.LogError(errors.Wrap(err, "err with send Email DownloadCertificate"))
		return
	}

	err = s.p.UpdateWinnerCertStatus(billID, types.PrizeStatusSendEmail)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with s.p.UpdateWinnerCertStatus in ChoiceCertificate"))
		return
	}

	return
}
