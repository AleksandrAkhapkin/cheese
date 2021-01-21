package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/lilekov-studio/supreme-cheese/internal/supreme-cheese/service/mail"
	"github.com/lilekov-studio/supreme-cheese/internal/supreme-cheese/types"
	"github.com/lilekov-studio/supreme-cheese/pkg/infrastruct"
	"github.com/lilekov-studio/supreme-cheese/pkg/logger"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func (s *Service) CheckBill(userID, billID int, path string) {

	//отправляем файл на проверку
	fullBill, err := s.GoToFNS(path, userID, billID)
	if err != nil {
		//если ошибка не на стороне проверки чека - логируем
		if err != infrastruct.ErrorBillUnknownError {
			logger.LogError(errors.Wrap(err, "err with s.GoToFNS in CheckBill"))
		}
		//если есть ошибка в целом - логируем чек
		if err := s.p.UpdateInnerStatusBill(billID, types.InnerStatusBillNotCanScan); err != nil {
			logger.LogError(errors.Wrap(err, "err with s.p.UpdateInnerStatusBill in CheckBill"))
			return
		}
		logger.Postphoto(s.pathForBill, userID, billID)
		return
	}

	//если чек получен - обрабатываем его
	switch fullBill.Code {
	case 0:
		{
			if err = s.p.UpdateInnerStatusBill(billID, types.InnerStatusBillScanButInvalid); err != nil {
				logger.LogError(errors.Wrap(err, "err with s.p.UpdateInnerStatusBill in CheckBill"))
				return
			}
			if err = s.p.UpdateStatusBillForUser(billID, types.StatusBillWait); err != nil {
				logger.LogError(errors.Wrap(err, "err with s.p.UpdateInnerStatusBill in CheckBill"))
				return
			}

			return
		}
	case 1:
		{
			//проверяем на дубль
			err = s.checkDoubleBill(fullBill)
			switch err {
			case infrastruct.ErrorBillDoubleError:
				{
					if err = s.p.WriteGoodResponseBill(userID, billID, fullBill); err != nil {
						logger.LogError(errors.Wrap(err, "err with s.p.WriteGoodResponseBill in CheckBill"))
						return
					}
					if err = s.p.UpdateInnerStatusBill(billID, types.StatusBillDouble); err != nil {
						logger.LogError(errors.Wrap(err, "err with s.p.UpdateInnerStatusBill in CheckBill"))
						return
					}
					if err = s.p.UpdateStatusBillForUser(billID, types.StatusBillDouble); err != nil {
						logger.LogError(errors.Wrap(err, "err with s.p.UpdateStatusBillForUser in CheckBill"))
						return
					}

					user, err := s.p.GetUserByID(userID)
					if err != nil {
						logger.LogError(errors.Wrap(err, "err with GetUserByID in InvalidBill"))

					}

					r := mail.NewRequest([]string{user.Email}, s.email)
					if err := r.Send("Проверка чека", "Мы отклонили чек, так как он был проверен ранее"); err != nil {
						logger.LogError(errors.Wrap(err, "err with send Email"))
						return
					}

					return
				}
			case nil:
				{
					checkTime, err := time.Parse("2006-01-02T15:04:05", fullBill.Data.JSON.Date)
					if err != nil {
						logger.LogError(errors.Wrap(err, "err with checkTime, err := time.Parse in CheckBill"))
						return
					}
					timeStart, err := time.Parse("2006-01-02T15:04:05", "2020-11-22T00:00:00")
					if err != nil {
						logger.LogError(errors.Wrap(err, "err with checkTime, err := time.Parse in CheckBill"))
						return
					}

					if checkTime.Before(timeStart) {
						if err = s.p.UpdateInnerStatusBill(billID, types.InnerStatusBillInvalidTimeBill); err != nil {
							logger.LogError(errors.Wrap(err, "err with s.p.UpdateInnerStatusBill in CheckBill"))
							return
						}
						if err = s.p.UpdateStatusBillForUser(billID, types.StatusBillWait); err != nil {
							logger.LogError(errors.Wrap(err, "err with s.p.UpdateStatusBillForUser in CheckBill"))
							return
						}
						return
					}
					if err = s.p.WriteGoodResponseBill(userID, billID, fullBill); err != nil {
						logger.LogError(errors.Wrap(err, "err with s.p.WriteGoodResponseBill in CheckBill"))
						return
					}

					if err = s.p.WriteBillPosition(userID, billID, fullBill.Data.JSON.Items); err != nil {
						logger.LogError(errors.Wrap(err, "err with s.p.WriteBillPosition in CheckBill"))
						return
					}

					if err = s.p.UpdateInnerStatusBill(billID, types.InnerStatusBillNotFindFACKINGCheese); err != nil {
						logger.LogError(errors.Wrap(err, "err with s.p.UpdateInnerStatusBill in CheckBill"))
						return
					}

					haveCheese, err := s.finderCheeseInBill(billID, fullBill.Data.JSON.Items)
					if err != nil {
						logger.LogError(errors.Wrap(err, "err with s.p.finderCheeseInBill in CheckBill"))
						return
					}

					if haveCheese {
						if err = s.p.UpdateInnerStatusBill(billID, types.StatusBillValid); err != nil {
							logger.LogError(errors.Wrap(err, "err with s.p.UpdateInnerStatusBill in CheckBill"))
							return
						}
						if err = s.p.UpdateStatusBillForUser(billID, types.StatusBillValid); err != nil {
							logger.LogError(errors.Wrap(err, "err with s.p.UpdateInnerStatusBill in CheckBill"))
							return
						}

						go s.PrizeLogic(billID)

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
					return
				}
			default:
				{
					return
				}
			}
		}
	default:
		if err = s.p.UpdateInnerStatusBill(billID, types.InnerStatusBillScanButInvalid); err != nil {
			logger.LogError(errors.Wrap(err, "err with s.p.UpdateInnerStatusBill in CheckBill"))
			return
		}

		return
	}
}

func (s *Service) GoToFNS(path string, userID, billID int) (*types.CheckBillReq, error) {

	file, _ := os.Open(path)
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("qrfile", filepath.Base(file.Name()))
	if err != nil {
		return nil, errors.Wrap(err, "err while CreateFormFile")
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return nil, errors.Wrap(err, "err while Copy")
	}

	field, err := writer.CreateFormField("token")
	if err != nil {
		return nil, errors.Wrap(err, "err while CreateFormField")
	}
	defer writer.Close()

	_, err = io.Copy(field, strings.NewReader(s.secretGiftery))
	if err != nil {
		return nil, errors.Wrap(err, "err while Copy")
	}

	r, err := http.NewRequest("POST", types.CheckBillURL, body)
	if err != nil {
		return nil, errors.Wrap(err, "err while NewRequest")
	}

	r.Header.Add("Content-Type", writer.FormDataContentType())
	res, err := http.DefaultClient.Do(r)
	if err != nil {
		return nil, errors.Wrap(err, "err while Do")
	}

	//логируем ответ
	respLogerBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with read body in GoToFNS for logger"))
	}
	if err = s.p.AddLogCheckBillResp(string(respLogerBytes), userID, billID); err != nil {
		logger.LogError(errors.Wrap(err, "err with AddLogCheckBillResp in GoToFNS for logger"))
	}

	if res.StatusCode != 200 {
		return nil, infrastruct.ErrorBillUnknownError
	}

	fullBill := types.CheckBillReq{}
	if err := json.Unmarshal(respLogerBytes, &fullBill); err != nil {
		return nil, errors.Wrap(err, "err with json.NewDecoder")
	}
	fmt.Printf("%v", fullBill)

	return &fullBill, nil
}
