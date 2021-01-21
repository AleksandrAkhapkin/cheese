package service

import (
	"fmt"
	"github.com/lilekov-studio/supreme-cheese/internal/clients/yoomoney"
	"github.com/lilekov-studio/supreme-cheese/internal/supreme-cheese/service/mail"
	"github.com/lilekov-studio/supreme-cheese/internal/supreme-cheese/types"
	"github.com/lilekov-studio/supreme-cheese/pkg/infrastruct"
	"github.com/lilekov-studio/supreme-cheese/pkg/logger"
	"github.com/pkg/errors"
)

func (s *Service) AddCheeseName(newName string) error {

	if err := s.p.AddCheeseName(newName); err != nil {
		logger.LogError(errors.Wrap(err, "err with s.p.AddCheeseName"))
		return infrastruct.ErrorInternalServerError
	}

	go s.reReadBills()

	return nil
}

func (s *Service) reReadBills() {

	//находим массив всех чеков со статусом "NOT FIND CHEESE"
	arrayOfWaitBillsID, err := s.p.GetArrayNotFindCheeseBills()
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with s.p.GetArrayValidWaitBills in reReadBills"))
		return
	}

	for i := range arrayOfWaitBillsID {
		//получаем массив всех позиций чеков по айди чека
		arrayOfWaitBillsPosition, err := s.p.GetPositionByBillID(arrayOfWaitBillsID[i])
		if err != nil {
			logger.LogError(errors.Wrap(err, "err with s.p.GetPositionByBillID in reReadBills"))
			return
		}
		//ищем сыр
		haveCheese, err := s.finderCheeseInBill(arrayOfWaitBillsID[i], arrayOfWaitBillsPosition)
		if err != nil {
			logger.LogError(errors.Wrap(err, "err with s.p.finderCheeseInBill in CheckBill"))
			return
		}

		//если сыр есть меняем статус
		if haveCheese {
			if err = s.p.UpdateInnerStatusBill(arrayOfWaitBillsID[i], types.StatusBillValid); err != nil {
				logger.LogError(errors.Wrap(err, "err with s.p.UpdateInnerStatusBill in CheckBill"))
				return
			}
			if err = s.p.UpdateStatusBillForUser(arrayOfWaitBillsID[i], types.StatusBillValid); err != nil {
				logger.LogError(errors.Wrap(err, "err with s.p.UpdateInnerStatusBill in CheckBill"))
				return
			}

			go s.PrizeLogic(arrayOfWaitBillsID[i])

			userID, err := s.p.GetUserIDByBillID(arrayOfWaitBillsID[i])
			if err != nil {
				logger.LogError(errors.Wrap(err, "err with GetUserIDByBillID in InvalidBill"))
				return
			}

			user, err := s.p.GetUserByID(userID)
			if err != nil {
				logger.LogError(errors.Wrap(err, "err with GetUserByID in InvalidBill"))

			}

			r := mail.NewRequest([]string{user.Email}, s.email)
			if err := r.Send("Проверка чека", "Ваш чек был успешно проверен!"); err != nil {
				logger.LogError(errors.Wrap(err, "err with send Email"))
				return
			}
		}
	}
}

func (s *Service) HandWriteBill(billResp *types.HandWriteBill) error {

	userID, err := s.p.GetUserIDByBillID(billResp.BillID)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with s.p.GetUserIDByBillID in HandWriteBill"))
		return err
	}

	//отправляем файл на проверку
	fullBill, err := s.HandWriteGoToFNS(userID, billResp)
	if err != nil {
		//записываем статус ошибки
		if err := s.p.UpdateInnerStatusBill(billResp.BillID, types.InnerStatusBillErrorInHandWrite); err != nil {
			return errors.Wrap(err, "err with s.p.UpdateInnerStatusBill in HandWriteBill")

		}
		if err = s.p.UpdateStatusBillForUser(billResp.BillID, types.StatusBillWait); err != nil {
			return errors.Wrap(err, "err with s.p.UpdateInnerStatusBill in HandWriteBill")
		}
		return err
	}

	//если чек получен - обрабатываем его
	switch fullBill.Code {
	case 0:
		{
			if err = s.p.UpdateInnerStatusBill(billResp.BillID, types.InnerStatusBillScanButInvalid); err != nil {
				return errors.Wrap(err, "err with s.p.UpdateInnerStatusBill in HandWriteBill")
			}
			if err = s.p.UpdateStatusBillForUser(billResp.BillID, types.StatusBillWait); err != nil {
				return errors.Wrap(err, "err with s.p.UpdateInnerStatusBill in HandWriteBill")
			}

			return nil
		}
	case 1:
		{
			if err = s.p.WriteGoodResponseBill(userID, billResp.BillID, fullBill); err != nil {
				return errors.Wrap(err, "err with s.p.WriteGoodResponseBill in HandWriteBill")
			}

			err = s.checkDoubleBillForAdmin(fullBill)
			switch err {
			case infrastruct.ErrorBillDoubleError:
				{
					if err = s.p.UpdateInnerStatusBill(billResp.BillID, types.StatusBillDouble); err != nil {
						return errors.Wrap(err, "err with s.p.UpdateInnerStatusBill in HandWriteBill")
					}
					if err = s.p.UpdateStatusBillForUser(billResp.BillID, types.StatusBillDouble); err != nil {
						return errors.Wrap(err, "err with s.p.UpdateStatusBillForUser in HandWriteBill")
					}
					return nil
				}
			case nil:
				{
					if err = s.p.WriteBillPosition(userID, billResp.BillID, fullBill.Data.JSON.Items); err != nil {
						return errors.Wrap(err, "err with s.p.WriteBillPosition in HandWriteBill")
					}

					if err = s.p.UpdateInnerStatusBill(billResp.BillID, types.InnerStatusBillNotFindFACKINGCheese); err != nil {
						return errors.Wrap(err, "err with s.p.UpdateInnerStatusBill in HandWriteBill")
					}

					if err = s.p.UpdateStatusBillForUser(billResp.BillID, types.StatusBillWait); err != nil {
						return errors.Wrap(err, "err with s.p.UpdateInnerStatusBill in HandWriteBill")
					}

					haveCheese, err := s.finderCheeseInBill(billResp.BillID, fullBill.Data.JSON.Items)
					if err != nil {
						return errors.Wrap(err, "err with s.p.finderCheeseInBill in HandWriteBill")
					}

					if haveCheese {
						if err = s.p.UpdateInnerStatusBill(billResp.BillID, types.StatusBillValid); err != nil {
							return errors.Wrap(err, "err with s.p.UpdateInnerStatusBill in HandWriteBill")
						}
						if err = s.p.UpdateStatusBillForUser(billResp.BillID, types.StatusBillValid); err != nil {
							return errors.Wrap(err, "err with s.p.UpdateInnerStatusBill in HandWriteBill")
						}

						go s.PrizeLogic(billResp.BillID)

						userID, err := s.p.GetUserIDByBillID(billResp.BillID)
						if err != nil {
							return errors.Wrap(err, "err with s.p.GetUserIDByBillID in HandWriteBill")
						}

						user, err := s.p.GetUserByID(userID)
						if err != nil {
							return errors.Wrap(err, "err with s.p.GetUserByID in HandWriteBill")
						}

						r := mail.NewRequest([]string{user.Email}, s.email)
						if err := r.Send("Проверка чека", "Ваш чек был успешно проверен!"); err != nil {
							logger.LogError(errors.Wrap(err, "err with send Email"))
							return errors.Wrap(err, "err with r.Send in HandWriteBill")
						}
					}
					return nil
				}
			default:
				{
					return err
				}
			}
		}
	default:
		if err = s.p.UpdateInnerStatusBill(billResp.BillID, types.InnerStatusBillScanButInvalid); err != nil {
			logger.LogError(errors.Wrap(err, "err with s.p.UpdateInnerStatusBill in CheckBill"))
			return err
		}

		return nil
	}
}

func (s *Service) SendCashback(send *types.SendCashback) error {

	return yoomoney.SendCashbackYoomoney(send.Phone, send.Amount, s.p)
}

func (s *Service) InvalidBill(invalidBill *types.InvalidBill) error {

	if err := s.p.UpdateInnerStatusBill(invalidBill.BillID, types.StatusBillInvalid); err != nil {
		return errors.Wrap(err, "err with s.p.UpdateInnerStatusBill in InvalidBill")

	}
	if err := s.p.UpdateStatusBillForUser(invalidBill.BillID, types.StatusBillInvalid); err != nil {
		return errors.Wrap(err, "err with s.p.UpdateInnerStatusBill in InvalidBill")
	}

	userID, err := s.p.GetUserIDByBillID(invalidBill.BillID)
	if err != nil {
		return errors.Wrap(err, "err with GetUserIDByBillID in InvalidBill")
	}

	user, err := s.p.GetUserByID(userID)
	if err != nil {
		return errors.Wrap(err, "err with GetUserByID in InvalidBill")

	}
	body := fmt.Sprintf("К сожалению, ваш чек был отклонен по причине того, что %s", invalidBill.Text)
	r := mail.NewRequest([]string{user.Email}, s.email)
	if err := r.Send("Проверка чека", body); err != nil {
		logger.LogError(errors.Wrap(err, "err with send Email"))
		return infrastruct.ErrorInternalServerError
	}

	return nil
}
