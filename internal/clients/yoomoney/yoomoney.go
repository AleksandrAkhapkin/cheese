package yoomoney

import (
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/lilekov-studio/supreme-cheese/internal/clients/mysql"
	"github.com/lilekov-studio/supreme-cheese/internal/supreme-cheese/types"
	"github.com/lilekov-studio/supreme-cheese/pkg/logger"
	"github.com/pkg/errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	urlRequestStart  = "https://yoomoney.ru/api/request-payment"
	urlRequestFinish = "https://yoomoney.ru/api/process-payment"
)

var authorization = ""

func NewYoomoney(bearer string) error {

	authorization = bearer

	return nil
}

func SendCashbackYoomoney(number string, amount int, p *mysql.MySQL) error {
	number = strings.TrimSpace(number)
	number = strings.TrimPrefix(number, "+")
	number = strings.TrimPrefix(number, "8")
	ind := strings.Index(number, "9")
	if ind < 0 {
		return fmt.Errorf("номер неверного формата, номер: %s", number)
	}
	if !strings.HasPrefix(number, "7") {
		number = "7" + number
	}

	err := payment(number, amount, p)
	if err != nil {
		return errors.Wrap(err, "err with payment")
	}
	return nil
}

func payment(number string, amount int, p *mysql.MySQL) error {

	user, err := p.GetUserByPhone(number)
	if err != nil {
		return errors.Wrap(err, "err with GetUserByPhone in payment (yoomoney)")
	}

	//создаем оффер
	am := strconv.Itoa(amount)
	am = am + ".00"
	data := url.Values{}
	data.Set("amount", am)
	data.Set("pattern_id", "phone-topup")
	data.Set("phone-number", number)
	r, err := http.NewRequest("POST", urlRequestStart, strings.NewReader(data.Encode()))
	if err != nil {
		return errors.Wrap(err, "err with NewRequest in payment")
	}
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Authorization", authorization)
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	res, err := http.DefaultClient.Do(r)
	if err != nil {
		return errors.Wrap(err, "err with DefaultClient (создание оффера) in payment (yoomoney)")
	}
	if res.StatusCode == http.StatusUnauthorized {
		return errors.Wrap(err, "err with DefaultClient.Do (создание оффера) in payment (yoomoney)")
	}

	logStart := types.PayRes{}
	if err = json.NewDecoder(res.Body).Decode(&logStart); err != nil {
		return errors.Wrap(err, "err with NewDecoder (logStart) in payment (yoomoney)")
	}

	if err = p.YoomoneyLoggerStart(&logStart, number, am, user.ID); err != nil {
		return errors.Wrap(err, "err with YoomoneyLoggerStart in in payment (yoomoney)")
	}
	//проводим оплату
	data = url.Values{}
	data.Set("request_id", logStart.Request_id)
	r, err = http.NewRequest("POST", urlRequestFinish, strings.NewReader(data.Encode()))
	if err != nil {
		return errors.Wrap(err, "err with NewRequest in payment")
	}
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Authorization", authorization)
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	res, err = http.DefaultClient.Do(r)
	if err != nil {
		return err
	}
	if res.StatusCode == http.StatusUnauthorized {
		return errors.Wrap(err, "err with DefaultClient.Do in payment (yoomoney)")
	}

	logFinal := types.PayedRes{}
	if err = json.NewDecoder(res.Body).Decode(&logFinal); err != nil {
		return errors.Wrap(err, "err with NewDecoder(logFinal) in payment (yoomoney)")
	}

	if err = p.YoomoneyLoggerFinal(&logFinal, number, am, user.ID, logStart.Request_id); err != nil {
		return errors.Wrap(err, "err with YoomoneyLoggerFinal in in payment (yoomoney)")
	}

	if logFinal.Error != "" {
		if logFinal.Error == "not_enough_funds" {
			logger.LogInfo("err with payment(YOOMONEY): ДЕНЬГИ ЗАКОНЧИЛИСЬ - РАСЧЕХЛЯЙ КРЕДИТКУ ЖЕНЕ4КА")
			return fmt.Errorf("NoMoneyInYoomoney")
		}
		return errors.Wrap(errors.Errorf(logFinal.Error), "err with send money in payment(YOOMONEY)")
	}

	return nil
}
