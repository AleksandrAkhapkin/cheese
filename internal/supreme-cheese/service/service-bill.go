package service

import (
	"database/sql"
	"fmt"
	"github.com/lilekov-studio/supreme-cheese/internal/supreme-cheese/types"
	"github.com/lilekov-studio/supreme-cheese/pkg/infrastruct"
	"github.com/lilekov-studio/supreme-cheese/pkg/logger"
	"github.com/pkg/errors"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func (s *Service) VerifyPhone(id int) (*types.CheckConfirmationPhone, error) {

	//проверяем подтверждение номера, если не подвержден - возвращаем ошибку
	confirmation, err := s.p.CheckConfirmationPhone(id)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with s.p.CheckConfirmationPhone in VerifyPhone"))
		return nil, infrastruct.ErrorInternalServerError
	}

	return confirmation, nil
}

func (s *Service) UploadFile(file *types.UploadFile) error {

	if file.PerecrestokBool {
		file.PerecrestokCard = strings.ReplaceAll(file.PerecrestokCard, " ", "")
		if err := s.p.AddPerekCard(file.UserID, file.PerecrestokBool, file.PerecrestokCard); err != nil {
			logger.LogError(errors.Wrap(err, "err with s.p.AddPerekCard in RegisterPhone"))
			return infrastruct.ErrorInternalServerError
		}
	} else {
		file.PerecrestokCard = ""
	}
	if file.PerecrestokCard == "" {
		file.PerecrestokBool = false
	}

	//проверяем подтверждение номера, если не подвержден возвращаем ошибку
	confirmation, err := s.p.CheckConfirmationPhone(file.UserID)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with s.p.CheckConfirmationPhone in UploadFile"))
		return infrastruct.ErrorInternalServerError
	}
	if !confirmation.ConfirmPhone {
		return infrastruct.ErrorNotConfirmedPhone
	}

	//записываем что пользователь отправил нам файл
	billID, err := s.p.WriteBillByUser(file.UserID, file.PerecrestokBool, file.PerecrestokCard)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with s.p.WriteBillByUser in ScanAndCheckBill"))
		return infrastruct.ErrorInternalServerError
	}

	if err := os.Mkdir(filepath.Join(s.pathForBill, fmt.Sprintf("%d", file.UserID)), 0777); err != nil {
		if !os.IsExist(err) {
			logger.LogError(errors.Wrap(err, "err with os.Mkdir in UploadFile"))
			return infrastruct.ErrorInternalServerError
		}
	}

	path := filepath.Join(s.pathForBill, fmt.Sprintf("%d/bill_id_%d", file.UserID, billID))
	dst, err := os.Create(filepath.Join(path))
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with UploadFile in create file"))
		return infrastruct.ErrorInternalServerError
	}

	_, err = io.Copy(dst, file.Body)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with UploadFile in io.Copy"))
		return infrastruct.ErrorInternalServerError
	}

	go s.ScanAndCheckBill(file.UserID, int(billID), path)

	return nil
}

func (s *Service) ScanAndCheckBill(userID, billID int, path string) {

	//отправлем файл на проверку
	s.CheckBill(userID, billID, path)
	return
}

func (s *Service) checkDoubleBill(fullBill *types.CheckBillReq) error {

	_, err := s.p.FindBillByFdFpFnInGoodResponse(fullBill.Data.JSON.FiscalDocumentNumber, fullBill.Data.JSON.FiscalSign, fullBill.Data.JSON.FiscalDriveNumber)
	if err != nil {
		if err != sql.ErrNoRows {
			logger.LogError(errors.Wrap(err, "err with s.p.FindBillByFdFpFnInGoodResponse in checkDoubleBill"))
			return infrastruct.ErrorInternalServerError
		}
		return nil
	}

	return infrastruct.ErrorBillDoubleError
}

func (s *Service) checkDoubleBillForAdmin(fullBill *types.CheckBillReq) error {

	count, err := s.p.CountBillByFdFpFnForAdmin(fullBill.Data.JSON.FiscalDocumentNumber, fullBill.Data.JSON.FiscalSign, fullBill.Data.JSON.FiscalDriveNumber)
	if err != nil {
		if err != sql.ErrNoRows {
			logger.LogError(errors.Wrap(err, "err with s.p.CountBillByFdFpFnForAdmin in checkDoubleBillForAdmin"))
			return infrastruct.ErrorInternalServerError
		}
		return nil
	}

	if count > 1 {
		return infrastruct.ErrorBillDoubleError
	}

	return nil
}
