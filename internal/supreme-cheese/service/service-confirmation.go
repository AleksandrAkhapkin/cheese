package service

import (
	"database/sql"
	"fmt"
	"github.com/lilekov-studio/supreme-cheese/internal/supreme-cheese/types"
	"github.com/lilekov-studio/supreme-cheese/pkg/infrastruct"
	"github.com/lilekov-studio/supreme-cheese/pkg/logger"
	"github.com/pkg/errors"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func (s *Service) ConfirmationPhone(conf *types.ConfirmationPhone) error {

	conf.Code = strings.ReplaceAll(conf.Code, " ", "")
	conf.Code = strings.ToUpper(conf.Code)

	if err := s.p.CheckCodeInConfPhone(conf); err != nil {
		if err == sql.ErrNoRows {
			return infrastruct.ErrorCodeIsIncorrect
		}
		logger.LogError(errors.Wrap(err, "err with CheckCodeInConfEmail in ConfirmationPhone"))
		return infrastruct.ErrorInternalServerError
	}

	if err := s.p.UpdateStatusPhone(conf.Phone); err != nil {
		logger.LogError(errors.Wrap(err, "err with s.p.UpdateStatusPhone in ConfirmationPhone"))
		return infrastruct.ErrorInternalServerError
	}

	if err := s.p.DeleteConfPhone(conf.Phone); err != nil {
		logger.LogError(errors.Wrap(err, "err with s.p.DeleteConfEmail in RegisterUserConfirmationEmail"))
		return infrastruct.ErrorInternalServerError
	}

	return nil
}

func (s *Service) AuthByToken(authToken *types.Token) (*types.Token, error) {

	//получаем емейл пользователя по токену
	email, err := s.p.GetUserEmailByTokenInConfURL(authToken.Token)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, infrastruct.ErrorCodeIsIncorrect
		}
		logger.LogError(errors.Wrap(err, "err with CheckCodeInConfEmail in AuthByToken"))
		return nil, infrastruct.ErrorInternalServerError
	}

	//если токен валидный - меняем статус емейла
	if err := s.p.ConfirmEmailStatus(email); err != nil {
		logger.LogError(errors.Wrap(err, "err with s.p.UpdateStatusPhone in AuthByToken"))
		return nil, infrastruct.ErrorInternalServerError
	}

	//удаляем запись с токеном
	if err := s.p.DeleteConfURL(email); err != nil {
		if err != sql.ErrNoRows {
			logger.LogError(errors.Wrap(err, "err with s.p.DeleteConfEmail in AuthByToken"))
			return nil, infrastruct.ErrorInternalServerError
		}
	}

	//получаем айди пользователя по емейлу
	user, err := s.p.GetUserByEmail(email)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with s.p.GetUserByEmail in AuthByToken"))
		return nil, infrastruct.ErrorInternalServerError
	}

	//Генерируем джвтэшку
	token, err := infrastruct.GenerateJWT(user.ID, types.RoleUser, s.secretKey)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with GenerateJWT in AuthByToken"))
		return nil, infrastruct.ErrorInternalServerError
	}

	return &types.Token{Token: token}, nil
}

func (s *Service) makeAndSendConfirmationPhoneCode(id int) error {

	user, err := s.p.GetUserByID(id)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with s.p.GetUserByID in makeAndSendConfirmationPhoneCode"))
		return infrastruct.ErrorInternalServerError
	}
	phone := user.Phone

	//проверяем есть ли уже сгенерированный код - если есть, удаляем и делаем новый
	have, err := s.p.CheckDBConfPhone(phone)
	if err != nil {
		if err != sql.ErrNoRows {
			logger.LogError(errors.Wrap(err, "err with CheckDBConfPhone in makeAndSendConfirmationPhoneCode"))
			return infrastruct.ErrorInternalServerError
		}
	}
	if have {
		if err := s.p.DeleteConfPhone(phone); err != nil {
			logger.LogError(errors.Wrap(err, "err with DeleteConfPhone in makeAndSendConfirmationPhoneCode"))
			return infrastruct.ErrorInternalServerError
		}
	}

	//генерируем код
	rand.Seed(time.Now().UnixNano())
	num := 999 + rand.Intn(9000)
	code := strconv.Itoa(num)

	if err := s.p.AddCodeForConfPhone(phone, code, id); err != nil {
		logger.LogError(errors.Wrap(err, "err with AddCodeForConfPhone in makeAndSendConfirmationPhoneCode"))
		return infrastruct.ErrorInternalServerError
	}

	mess := make([]types.SMSMessages, 0)
	mes := types.SMSMessages{Phone: phone, ClientID: id, Text: code}
	mess = append(mess, mes)

	reqID, err := s.p.LogSMS(phone, id)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with LogSMS in makeAndSendConfirmationPhoneCode"))
	}

	inlimit, err := s.inLimitSMS(id)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with inLimitSMS in makeAndSendConfirmationPhoneCode"))
		return infrastruct.ErrorInternalServerError
	}
	if !inlimit {
		logger.LogInfo(fmt.Sprintf("ПРЕВЫШЕН ЛИМИТ СМС ЮЗЕР - %d", id))
		return infrastruct.ErrorMuchSms
	}

	resp, err := s.SendSMS(mess)
	if err != nil {
		return infrastruct.ErrorInternalServerError
	}

	if err := s.p.AddLogRespSMS(string(resp), reqID, id); err != nil {
		logger.LogError(errors.Wrap(err, "err with AddLogRespSMS in makeAndSendConfirmationPhoneCode"))
	}

	return nil
}
