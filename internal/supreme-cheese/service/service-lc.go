package service

import (
	"database/sql"
	"github.com/hashicorp/go-uuid"
	"github.com/lilekov-studio/supreme-cheese/internal/supreme-cheese/types"
	"github.com/lilekov-studio/supreme-cheese/pkg/infrastruct"
	"github.com/lilekov-studio/supreme-cheese/pkg/logger"
	"github.com/pkg/errors"
	"strings"
)

func (s *Service) RegisterUser(user *types.UserRegister, ip string) error {

	user.Role = types.RoleUser
	replaceSpace(user)
	user.Email = strings.ToLower(user.Email)
	var err error

	user.City, err = s.findCityByURL(ip)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with check IP city"))
	}

	//проверяем наличие в базе по емейлу
	userInDB, err := s.p.GetUserByEmail(user.Email)
	if err != nil {
		//если нету в базе - ок
		if err != sql.ErrNoRows {
			logger.LogError(errors.Wrap(err, "err with GetUserByEmail in RegisterUser"))
			return infrastruct.ErrorInternalServerError
		}
		//если есть - смотри подтвержден ли емейл
	} else if userInDB.ConfirmEmail {
		//если подтвержден - возвращаем ошибку
		return infrastruct.ErrorEmailIsExist
	} else {
		//если не подтвержден, удаялем строчку
		if err = s.p.DeleteUserByID(userInDB.ID); err != nil {
			logger.LogError(errors.Wrap(err, "err with DeleteUserByID in RegisterUser"))
			return infrastruct.ErrorInternalServerError
		}
	}

	//проверяем наличие в базе по фону
	userByPhone, err := s.p.GetUserByPhone(user.Phone)
	if err != nil {
		//если ошибка не ноуровс
		if err != sql.ErrNoRows {
			//логируем и завершаем
			logger.LogError(errors.Wrap(err, "err with GetUserIDByPhone in RegisterUser"))
			return infrastruct.ErrorInternalServerError
		}
		//если пользователь с таким телефоном зарегистрирован проверяем подтверждение почты
	} else if userByPhone.ConfirmEmail {
		//если подтверждена возвращаем ошибку
		return infrastruct.ErrorPhoneIsExist
	} else {
		//если не подтвержден, удаялем строчку Л
		if err = s.p.DeleteUserByID(userByPhone.ID); err != nil {
			logger.LogError(errors.Wrap(err, "err with DeleteUserByID in RegisterUser"))
			return infrastruct.ErrorInternalServerError
		}
	}

	//генерируем пароль для отправки на почту
	tmpPass, err := uuid.GenerateUUID()
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with GenerateUUID in sendPassToEmail"))
		return infrastruct.ErrorInternalServerError
	}
	tmpPass = tmpPass[:8]
	tmpPass = strings.ToUpper(tmpPass)
	user.GeneratePass = tmpPass

	//создаем юзера с паролем
	id, err := s.p.CreateUser(user)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with CreateUser in RegisterUser"))
		return infrastruct.ErrorInternalServerError
	}

	//генерируем ссылку для входа
	urlKey, err := uuid.GenerateUUID()
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with GenerateUUID in RegisterUser"))
		return infrastruct.ErrorInternalServerError
	}

	//проверяем наличие ссылки в базе
	have, err := s.p.CheckDBConfURL(user.Email)
	if err != nil {
		if err != sql.ErrNoRows {
			logger.LogError(errors.Wrap(err, "err with CheckDBConfPhone in RegisterUser"))
			return infrastruct.ErrorInternalServerError
		}
	}
	if have {
		//если ссылка есть - удаляем
		if err := s.p.DeleteConfURL(user.Email); err != nil {
			logger.LogError(errors.Wrap(err, "err with DeleteConfPhone in RegisterUser"))
			return infrastruct.ErrorInternalServerError
		}
	}

	//добавляем ссылку в базу
	if err := s.p.WriteTokenForEmail(int(id), user.Email, urlKey); err != nil {
		logger.LogError(errors.Wrap(err, "err with s.p.WriteTokenForEmail in RegisterUser"))
		return infrastruct.ErrorInternalServerError
	}

	if err = s.sendPassToEmail(user.Email, user.GeneratePass, urlKey); err != nil {
		return infrastruct.ErrorInternalServerError
	}

	return nil
}

func (s *Service) RecoverPassword(rec *types.RecoverPass) error {

	rec.Email = strings.ToLower(rec.Email)
	rec.Email = strings.ReplaceAll(rec.Email, " ", "")

	//проверяем наличие в базе по емейлу
	userInDB, err := s.p.GetUserByEmail(rec.Email)
	if err != nil {
		//если нету в базе - пользователю не говорим, но сворачиваем движуху
		if err != sql.ErrNoRows {
			logger.LogError(errors.Wrap(err, "err with GetUserByEmail in RegisterUser"))
			return infrastruct.ErrorInternalServerError
		}
		return nil
	}

	//генерируем пароль
	tmpPass, err := uuid.GenerateUUID()
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with GenerateUUID in sendPassToEmail"))
		return infrastruct.ErrorInternalServerError
	}
	tmpPass = tmpPass[:8]
	tmpPass = strings.ToUpper(tmpPass)
	rec.GeneratePass = tmpPass

	//изменяем юзера
	err = s.p.PutUserPassByID(userInDB.ID, rec.GeneratePass)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with CreateUser in RegisterUser"))
		return infrastruct.ErrorInternalServerError
	}

	if err = s.sendPassToEmail(rec.Email, rec.GeneratePass, ""); err != nil {
		return infrastruct.ErrorInternalServerError
	}

	return nil
}

func (s *Service) RegisterPhone(phone *types.RegisterPhone, id int) error {

	//проверяем указанный телефон по базе
	userIDByDB, err := s.p.GetUserIDByPhone(phone.Phone)
	if err != nil && err != sql.ErrNoRows {
		return infrastruct.ErrorInternalServerError
	}

	//если айди пользователя с таким номером не совпадает с айди текущего юзера - говорим что номер уже зарегистрирован
	if userIDByDB != id && userIDByDB != 0 {
		return infrastruct.ErrorPhoneIsExist
	}

	//ищем старый номер по айди
	user, err := s.p.GetUserByID(id)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with s.p.GetUserByID in RegisterPhone"))
		return infrastruct.ErrorInternalServerError
	}

	if user.ConfirmPhone {
		return infrastruct.ErrorConfirmedPhone
	}

	//если номер пользователя изменился - перезаписываем
	if user.Phone != phone.Phone {
		if err = s.p.UpdatePhone(phone.Phone, id); err != nil {
			logger.LogError(errors.Wrap(err, "err with s.p.UpdatePhone in RegisterPhone"))
			return infrastruct.ErrorInternalServerError
		}
	}

	//отправляем код
	if err := s.makeAndSendConfirmationPhoneCode(id); err != nil {
		return err
	}

	return nil
}

func (s *Service) Authorize(auth *types.Authorize) (*types.Token, error) {

	auth.Email = strings.ToLower(auth.Email)
	auth.Email = strings.ReplaceAll(auth.Email, " ", "")
	auth.Pass = strings.ReplaceAll(auth.Pass, " ", "")

	user, err := s.p.GetUserByEmailHavePass(auth.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, infrastruct.ErrorEmailOrPasswordNotFind
		}
		logger.LogError(errors.Wrap(err, "err with GetUserByEmailHavePass in Authorize"))
		return nil, infrastruct.ErrorInternalServerError
	}

	if user.Pass != auth.Pass {
		return nil, infrastruct.ErrorEmailOrPasswordNotFind
	}

	token, err := infrastruct.GenerateJWT(user.ID, user.Role, s.secretKey)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with GenerateJWT in Authorize"))
		return nil, infrastruct.ErrorInternalServerError
	}

	if err = s.p.ConfirmEmailStatus(user.Email); err != nil {
		logger.LogError(errors.Wrap(err, "err with ConfirmEmailStatus in Authorize"))
		return nil, infrastruct.ErrorInternalServerError
	}

	if err := s.p.DeleteConfURL(user.Email); err != nil {
		if err != sql.ErrNoRows {
			logger.LogError(errors.Wrap(err, "err with s.p.DeleteConfEmail in Authorize"))
			return nil, infrastruct.ErrorInternalServerError
		}
	}

	return &types.Token{Token: token}, nil
}

func (s *Service) PutCabinet(user *types.PutCabinet) error {

	//проверяем наличие нового пароля
	if user.NewPass != "" {
		//находим страный пароль
		oldPassInDB, err := s.p.GetUserPassByID(user.ID)
		if err != nil {
			logger.LogError(errors.Wrap(err, "err with s.p.GetUserPassByID in PutCabinet"))
			return infrastruct.ErrorInternalServerError
		}
		//если старый из бызы и старый указанный юзером совпадают
		if oldPassInDB == user.OldPass {
			//если совпадают перезаписываем
			if err := s.p.PutUserPassByID(user.ID, user.NewPass); err != nil {
				logger.LogError(errors.Wrap(err, "err with PutUserNameByID in PutCabinet"))
				return infrastruct.ErrorInternalServerError
			}

		} else {
			return infrastruct.ErrorOldPasswordsDoNotMatch
		}
	}

	//перезаписываем имя
	//if err := s.p.PutUserNameByID(user); err != nil {
	//	logger.LogError(errors.Wrap(err, "err with PutUserNameByID in PutCabinet"))
	//	return infrastruct.ErrorInternalServerError
	//}

	userInDB, err := s.GetUserByID(user.ID)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with GetUserByID in PutCabinet"))
		return infrastruct.ErrorInternalServerError
	}

	//если номер не подтвержден - меняем
	if !userInDB.ConfirmPhone {
		if err := s.p.UpdatePhone(user.NewPhone, user.ID); err != nil {
			logger.LogError(errors.Wrap(err, "err with UpdatePhone in PutCabinet"))
			return infrastruct.ErrorInternalServerError
		}
	}

	return nil
}
