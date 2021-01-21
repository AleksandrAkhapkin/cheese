package service

import (
	"database/sql"
	"github.com/lilekov-studio/supreme-cheese/internal/clients/yoomoney"
	"github.com/lilekov-studio/supreme-cheese/internal/supreme-cheese/types"
	"github.com/lilekov-studio/supreme-cheese/pkg/infrastruct"
	"github.com/lilekov-studio/supreme-cheese/pkg/logger"
	"github.com/pkg/errors"
	"time"
)

func (s *Service) PrizeLogic(billID int) {

	err := s.cashBack(billID)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with s.cashBack in PrizeLogic"))
	}

	userID, err := s.p.GetUserIDByBillID(billID)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with GetUserIDByBillID in PrizeLogic"))
		return
	}

	haveWinCertificate, err := s.p.HaveWinCertificate(userID)
	if err != nil {
		if err != sql.ErrNoRows {
			logger.LogError(errors.Wrap(err, "err with s.p.HaveWinCertificate in PrizeLogic"))
			return
		}
		haveWinCertificate = false
	}

	if haveWinCertificate {
		if err := s.p.PUTWinDayStatus(billID, types.WinDayStatusDouble); err != nil {
			logger.LogError(errors.Wrap(err, "err with PUTWinDayStatus in PrizeLogic"))
			return
		}
		return
	}

	if err := s.p.PUTWinDayStatus(billID, types.WinDayStatusWait); err != nil {
		logger.LogError(errors.Wrap(err, "err with PUTWinDayStatus in PrizeLogic"))
		return
	}

	return
}

func (s *Service) cashBack(billID int) error {

	//быстрая надстройка
	//увеличенный кэшбек
	//проверяем надо ли удвоить кэшбек

	//получаем список всех кому надо удвоить
	bigCashUsers, err := s.p.GetBigCashUsers()
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with GetBigCashUsers in cashBack"))
		return infrastruct.ErrorInternalServerError
	}

	userID, err := s.p.GetUserIDByBillID(billID)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with GetUserIDByBillID in cashBack"))
		return infrastruct.ErrorInternalServerError
	}

	//получаем телефон пользователя которому проводится выплта
	user, err := s.GetUserByID(userID)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with GetUserByID in cashBack"))
		return infrastruct.ErrorInternalServerError
	}

	//ренджим массив номеров с удвоеным кэшбеком
	for _, values := range bigCashUsers {
		//если номер пользователя совпадает с номером из списка удвоения - переходим в увеличенную выплату
		if user.Phone == values {
			if err := s.sendBigCash(billID, user.Phone); err != nil {
				return err
			}
			return nil
		}
	}

	err = s.p.CashBackByBillID(billID)
	if err != nil {
		if err != sql.ErrNoRows {
			logger.LogError(errors.Wrap(err, "err with CashBackByBillID in cashBack"))
			return infrastruct.ErrorInternalServerError
		}
	} else {
		return errors.Errorf("%s", "По данному чеку была выпалата кэшбека")
	}

	perekCashBack, err := s.p.GetBoolPerekPrizeType(billID)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with GetBoolPerekPrizeType in cashBack"))
		return infrastruct.ErrorInternalServerError
	}

	if perekCashBack {
		if err := s.p.WriteWinner(userID, billID, types.PrizePerek, types.PrizeStatusPerek); err != nil {
			logger.LogError(errors.Wrap(err, "err with UpdateWinnerCashBackStatus in cashBack (perek)"))
			return infrastruct.ErrorInternalServerError
		}
	} else {
		if err := s.p.WriteWinner(userID, billID, types.PrizePhone, types.PrizeStatusWaitSend); err != nil {
			logger.LogError(errors.Wrap(err, "err with UpdateWinnerCashBackStatus in cashBack (perek)"))
			return infrastruct.ErrorInternalServerError
		}
	}

	countInProject, err := s.p.CountCashBackInProject(userID)
	if err != nil {
		if err != sql.ErrNoRows {
			logger.LogError(errors.Wrap(err, "err with CountCashBackInProject in cashBack"))
			return infrastruct.ErrorInternalServerError
		}
		countInProject = 0
	}

	if countInProject > types.ProjectLimitCashBack {
		if err := s.p.UpdateWinnerCashBackStatus(billID, types.PrizeStatusLimitInPromo); err != nil {
			logger.LogError(errors.Wrap(err, "err with UpdateWinnerCertStatus in cashBack"))
			return infrastruct.ErrorInternalServerError
		}
		return nil
	}

	if countInProject > types.DaysLimitCashBack {
		timeThreeCashAgo, err := s.p.GetTimeThreeCashAgo(userID)
		if err != nil {
			logger.LogError(errors.Wrap(err, "err with GetTimeThreeCashAgo in cashBack"))
			return infrastruct.ErrorInternalServerError
		}

		timeCashBack, err := time.Parse("2006-01-02 15:04:05", timeThreeCashAgo)
		if err != nil {
			logger.LogError(errors.Wrap(err, "err with time.Parse in cashBack"))
			return infrastruct.ErrorInternalServerError
		}

		timeCashBack = timeCashBack.Add(time.Hour * -3)
		dayAgoTime := time.Now().AddDate(0, 0, -1)

		//если чек был менее суток назад - меняеем статус и сворачиваем
		if dayAgoTime.Before(timeCashBack) {
			if err := s.p.UpdateWinnerCashBackStatus(billID, types.PrizeStatusLimitInDay); err != nil {
				logger.LogError(errors.Wrap(err, "err with UpdateWinnerCertStatus in cashBack"))
				return infrastruct.ErrorInternalServerError
			}
			return nil
		}
	}

	if !perekCashBack {
		user, err := s.p.GetUserByID(userID)
		if err != nil {
			logger.LogError(errors.Wrap(err, "err with GetUserByID in cashBack"))
			return infrastruct.ErrorInternalServerError
		}

		if err = yoomoney.SendCashbackYoomoney(user.Phone, types.CashBackSum, s.p); err != nil {
			logger.LogError(err)
			return infrastruct.ErrorInternalServerError
		}

		//if err := qiwi.SendCashback(user.Phone, types.CashBackSum); err != nil {
		//	logger.LogError(errors.Wrap(err, "err with SendCashback in cashBack "))
		//	return infrastruct.ErrorInternalServerError
		//}

		if err := s.p.UpdateWinnerCashBackStatus(billID, types.PrizeStatusPhone); err != nil {
			logger.LogError(errors.Wrap(err, "err with UpdateWinnerCertStatus in cashBack"))
			return infrastruct.ErrorInternalServerError
		}
	}

	return nil
}

func (s *Service) sendBigCash(billID int, phone string) error {

	userID, err := s.p.GetUserIDByBillID(billID)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with GetUserIDByBillID in sendBigCash"))
		return infrastruct.ErrorInternalServerError
	}

	err = s.p.CashBackByBillID(billID)
	if err != nil {
		if err != sql.ErrNoRows {
			logger.LogError(errors.Wrap(err, "err with CashBackByBillID in sendBigCash"))
			return infrastruct.ErrorInternalServerError
		}
	} else {
		return errors.Errorf("%s", "По данному чеку была выпалата кэшбека")
	}

	perekCashBack, err := s.p.GetBoolPerekPrizeType(billID)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with GetBoolPerekPrizeType in sendBigCash"))
		return infrastruct.ErrorInternalServerError
	}

	if perekCashBack {
		if err := s.p.WriteWinner(userID, billID, types.PrizePerekBIG, types.PrizeStatusPerek); err != nil {
			logger.LogError(errors.Wrap(err, "err with UpdateWinnerCashBackStatus in sendBigCash (perek)"))
			return infrastruct.ErrorInternalServerError
		}
	} else {
		if err := s.p.WriteWinner(userID, billID, types.PrizePhoneBIG, types.PrizeStatusWaitSend); err != nil {
			logger.LogError(errors.Wrap(err, "err with UpdateWinnerCashBackStatus in sendBigCash (perek)"))
			return infrastruct.ErrorInternalServerError
		}
	}

	countInProject, err := s.p.CountCashBackInProject(userID)
	if err != nil {
		if err != sql.ErrNoRows {
			logger.LogError(errors.Wrap(err, "err with CountCashBackInProject in sendBigCash"))
			return infrastruct.ErrorInternalServerError
		}
		countInProject = 0
	}

	if countInProject > types.ProjectLimitCashBack {
		if err := s.p.UpdateWinnerCashBackStatus(billID, types.PrizeStatusLimitInPromo); err != nil {
			logger.LogError(errors.Wrap(err, "err with UpdateWinnerCertStatus in sendBigCash"))
			return infrastruct.ErrorInternalServerError
		}
		return nil
	}

	if countInProject > types.DaysLimitCashBack {
		timeThreeCashAgo, err := s.p.GetTimeThreeCashAgo(userID)
		if err != nil {
			logger.LogError(errors.Wrap(err, "err with GetTimeThreeCashAgo in sendBigCash"))
			return infrastruct.ErrorInternalServerError
		}

		timeCashBack, err := time.Parse("2006-01-02 15:04:05", timeThreeCashAgo)
		if err != nil {
			logger.LogError(errors.Wrap(err, "err with time.Parse in sendBigCash"))
			return infrastruct.ErrorInternalServerError
		}

		timeCashBack = timeCashBack.Add(time.Hour * -3)
		dayAgoTime := time.Now().AddDate(0, 0, -1)

		//если чек был менее суток назад - меняеем статус и сворачиваем
		if dayAgoTime.Before(timeCashBack) {
			if err := s.p.UpdateWinnerCashBackStatus(billID, types.PrizeStatusLimitInDay); err != nil {
				logger.LogError(errors.Wrap(err, "err with UpdateWinnerCertStatus in sendBigCash"))
				return infrastruct.ErrorInternalServerError
			}
			return nil
		}
	}

	if !perekCashBack {
		user, err := s.p.GetUserByID(userID)
		if err != nil {
			logger.LogError(errors.Wrap(err, "err with GetUserByID in sendBigCash"))
			return infrastruct.ErrorInternalServerError
		}

		if err = yoomoney.SendCashbackYoomoney(user.Phone, types.CashBackSumBIG, s.p); err != nil {
			logger.LogError(err)
			return infrastruct.ErrorInternalServerError
		}

		//if err := qiwi.SendCashback(user.Phone, types.CashBackSumBIG); err != nil {
		//	logger.LogError(errors.Wrap(err, "err with SendCashback in sendBigCash "))
		//	return infrastruct.ErrorInternalServerError
		//}

		if err := s.p.UpdateWinnerCashBackStatus(billID, types.PrizeStatusPhone); err != nil {
			logger.LogError(errors.Wrap(err, "err with UpdateWinnerCertStatus in sendBigCash"))
			return infrastruct.ErrorInternalServerError
		}
	}

	err = s.p.DelBigWinUser(phone)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with DelBigWinUser"))
		return infrastruct.ErrorInternalServerError
	}

	return nil
}
