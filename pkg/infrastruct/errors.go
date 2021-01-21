package infrastruct

import "net/http"

type CustomError struct {
	msg  string
	Code int
}

func NewError(msg string, code int) *CustomError {
	return &CustomError{
		msg:  msg,
		Code: code,
	}
}

func (c *CustomError) Error() string {
	return c.msg
}

var (
	ErrorEmailIsExist           = NewError("Email уже зарегистрирован", http.StatusBadRequest)
	ErrorPhoneIsExist           = NewError("Телефон уже зарегистрирован", http.StatusBadRequest)
	ErrorInternalServerError    = NewError("Внутренняя ошибка сервера", http.StatusInternalServerError)
	ErrorBadRequest             = NewError("Плохие входные данные запроса", http.StatusBadRequest)
	ErrorJWTIsBroken            = NewError("jwt испорчен", http.StatusForbidden)
	ErrorPermissionDenied       = NewError("У вас недостаточно прав", http.StatusForbidden)
	ErrorEmailOrPasswordNotFind = NewError("Неверный емейл или пароль", http.StatusForbidden)
	ErrorCodeIsIncorrect        = NewError("Неверный код", http.StatusForbidden)
	ErrorNotConfirmedPhone      = NewError("Телефон не подтвержден", http.StatusConflict)
	ErrorConfirmedPhone         = NewError("Вы не можете изменить номер телефона после его подтверждения", http.StatusConflict)
	ErrorOldPasswordsDoNotMatch = NewError("Вы указали неверный старый пароль", http.StatusForbidden)

	ErrorBillUnknownError = NewError("Ошибка проверки чека", http.StatusInternalServerError)
	ErrorBillDoubleError  = NewError("Данный чек был зарегистрирован ранее", http.StatusInternalServerError)

	ErrorMuchSms = NewError("Ошибка отправки СМС, свяжитесь с нами для подтверждения номера", http.StatusForbidden)
	ErrorGiftery = NewError("Ошибка при заказе сертификата, попробуйте позднее", http.StatusInternalServerError)
)
