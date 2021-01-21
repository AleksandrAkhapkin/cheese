package service

import (
	"fmt"
	"github.com/lilekov-studio/supreme-cheese/internal/supreme-cheese/service/mail"
	"github.com/lilekov-studio/supreme-cheese/internal/supreme-cheese/types"
	"github.com/lilekov-studio/supreme-cheese/pkg/infrastruct"
	"github.com/lilekov-studio/supreme-cheese/pkg/logger"
	"github.com/pkg/errors"
	"net/url"
)

func (s *Service) sendPassToEmail(email, pass, uuid string) error {

	var body string
	if uuid == "" {
		body = fmt.Sprintf("Ваш пароль для входа на сайт: %s\n\n", pass)
		r := mail.NewRequest([]string{email}, s.email)
		if err := r.Send("Восстановление пароля", body); err != nil {
			logger.LogError(errors.Wrap(err, "err with send Email"))
			return infrastruct.ErrorInternalServerError
		}
	} else {
		body = fmt.Sprintf("Ваш пароль для входа на сайт: %s\n\nСсылка для авторизации: %s", pass, fmt.Sprintf(types.HomePage+"/"+types.TokenForAuthByURL+uuid))
		r := mail.NewRequest([]string{email}, s.email)
		if err := r.Send("Регистрация на сайте", body); err != nil {
			logger.LogError(errors.Wrap(err, "err with send Email"))
			return infrastruct.ErrorInternalServerError
		}
	}

	return nil
}

func (s *Service) SendEmailOS(os *types.OS, userInDB *types.User) error {

	var message string
	if userInDB.ID == 0 {
		message = fmt.Sprintf("Новое обращение №%d!\n\nИмя - %s\nКонтактный телефон - %s\nКонтактный email - %s\nОбращение: \"%s\"",
			os.MessageID, os.Name, os.Phone, os.Email, os.Text)
	} else {
		message = fmt.Sprintf("Новое обращение №%d!\n\nИмя - %s\nКонтактный телефон - %s\nКонтактный email - %s\nОбращение: \"%s\"\n\n\nДанный пользователь зарегистрирован на сайте, со следующими данными:\nID - %d\nИмя - %s\nФамилия - %s\nКонтактный телефон - %s\nКонтактный email - %s",
			os.MessageID, os.Name, os.Phone, os.Email, os.Text, userInDB.ID, userInDB.FirstName, userInDB.LastName, userInDB.Phone, userInDB.Email)
	}
	logger.LogInfo(url.QueryEscape(message))
	logger.LogOSForSasha(url.QueryEscape(message))
	r := mail.NewRequest([]string{"support@supreme-cheese.ru"}, s.email)
	if err := r.Send("Новое обращение к службе поддержки", message); err != nil {
		logger.LogError(errors.Wrap(err, "err with send Email"))
		return infrastruct.ErrorInternalServerError
	}

	return nil
}

func (s *Service) UnsubscribeMail(email string) error {

	if err := s.p.SetUnsubscribe(email); err != nil {
		logger.LogError(errors.Wrapf(err, "err while SetUnsubscribe, email: $s", email))
		return infrastruct.ErrorInternalServerError
	}

	return nil
}
