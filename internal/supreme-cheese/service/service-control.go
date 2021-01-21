package service

import (
	"database/sql"
	"fmt"
	"github.com/lilekov-studio/supreme-cheese/internal/supreme-cheese/types"
	"github.com/lilekov-studio/supreme-cheese/pkg/infrastruct"
	"github.com/lilekov-studio/supreme-cheese/pkg/logger"
	"github.com/pkg/errors"
)

func (s *Service) inLimitSMS(userID int) (bool, error) {

	count, err := s.p.CountSMSByUser(userID)
	if err != nil {
		return false, errors.Wrap(err, "err in CountSMSByUser in inLimitSMS")
	}

	if count > types.LimitForSendSMS {
		return false, nil
	}

	return true, nil

}

func (s *Service) CheckCanGetSertificate(userID int) error {

	_, err := s.p.FindCertificateWinnerHaveChoiceCertificate(userID)
	if err != nil {
		if err != sql.ErrNoRows {
			logger.LogError(errors.Wrap(err, "err in CountSMSByUser in inLimitSMS"))
			return infrastruct.ErrorInternalServerError
		}
		logger.LogInfo(fmt.Sprintf("Пользователь с ID %d пытается получить сертификат который ему не положен!", userID))
		return infrastruct.ErrorPermissionDenied
	}

	return nil

}
