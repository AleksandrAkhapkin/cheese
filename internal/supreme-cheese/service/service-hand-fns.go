package service

import (
	"encoding/json"
	"fmt"
	"github.com/lilekov-studio/supreme-cheese/internal/supreme-cheese/types"
	"github.com/lilekov-studio/supreme-cheese/pkg/infrastruct"
	"github.com/lilekov-studio/supreme-cheese/pkg/logger"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func (s *Service) HandWriteGoToFNS(userID int, bill *types.HandWriteBill) (*types.CheckBillReq, error) {

	if err := s.p.AddLogCheckBillReq(fmt.Sprintf("fn=%s&fd=%s&fp=%s&n=1&s=%s&t=%s&qr=0&token=%s",
		bill.FN, bill.FD, bill.FP, bill.Sum, bill.Date, s.secretProverkaCheka), userID, bill.BillID); err != nil {
		logger.LogError(errors.Wrap(err, "err with AddLogCheckBillResp in GoToFNS for logger"))
	}

	req, err := http.NewRequest(http.MethodPost, types.CheckBillURL, strings.NewReader(
		fmt.Sprintf("fn=%s&fd=%s&fp=%s&n=1&s=%s&t=%s&qr=0&token=%s",
			bill.FN, bill.FD, bill.FP, bill.Sum, bill.Date, s.secretProverkaCheka)))
	if err != nil {
		return nil, errors.Wrap(err, "err with http.NewRequest in HandWriteGoToFNS")
	}

	form := url.Values{}
	req.PostForm = form
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	res, _ := http.DefaultClient.Do(req)

	//логируем ответ
	respLogerBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with read body in GoToFNS for logger"))
	}
	if err = s.p.AddLogCheckBillResp(string(respLogerBytes), userID, bill.BillID); err != nil {
		logger.LogError(errors.Wrap(err, "err with AddLogCheckBillResp in GoToFNS for logger"))
	}

	if res.StatusCode != 200 {
		return nil, infrastruct.ErrorBillUnknownError
	}

	fullBill := types.CheckBillReq{}
	if err := json.Unmarshal(respLogerBytes, &fullBill); err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("err with json.NewDecoder. Ответ от ПровеиЧеков: %s", string(respLogerBytes)))
	}

	return &fullBill, nil
}
