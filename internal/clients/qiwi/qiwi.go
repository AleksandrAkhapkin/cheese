package qiwi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
	"strings"
	"time"
)

const (
	detectURL  = "https://qiwi.com/mobile/detect.action"
	paymentURL = "https://edge.qiwi.com/sinap/api/v2/terms/%s/payments"
)

var token = ""

func NewQiwi(bearer string) error {

	token = bearer

	return nil
}

type accounts struct {
	Accounts []struct {
		Balance struct {
			Amount float64 `json:"amount"`
		} `json:"balance"`
	} `json:"accounts"`
}

type codeReq struct {
	Message string `json:"message"`
}

type payReq struct {
	ID  string `json:"id"`
	Sum struct {
		Amount   float32 `json:"amount"`
		Currency string  `json:"currency"`
	} `json:"sum"`
	PaymentMethod struct {
		Type      string `json:"type"`
		AccountID string `json:"accountId"`
	} `json:"paymentMethod"`
	Fields struct {
		Account string `json:"account"`
	} `json:"fields"`
}

func SendCashback(number string, amount int) error {
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

	code, err := detect(number)
	if err != nil {
		return fmt.Errorf("этот номер недоступен для автопополнения: %s", number)
	}

	if len(code) > 7 {
		return fmt.Errorf("код номера превышен, code: %d, number: %s", code, number)
	}

	res, err := payment(number[1:], code, amount)
	if err != nil {
		return errors.Wrap(err, "err with payment")
	}

	if len(res) > 0 {
		return fmt.Errorf("код ответа payment превышен, res: %d, number: %s", res, number)
	}

	return nil
}

func payment(number, code string, amount int) (string, error) {
	p := payReq{
		ID: fmt.Sprintf("%d", time.Now().Unix()*1000),
		Sum: struct {
			Amount   float32 `json:"amount"`
			Currency string  `json:"currency"`
		}{
			Amount:   float32(amount),
			Currency: "643",
		},
		PaymentMethod: struct {
			Type      string `json:"type"`
			AccountID string `json:"accountId"`
		}{
			Type:      "Account",
			AccountID: "643",
		},
		Fields: struct {
			Account string `json:"account"`
		}{
			Account: number,
		},
	}

	b, err := json.Marshal(&p)
	if err != nil {
		return "", err
	}

	payURL := fmt.Sprintf(paymentURL, code)

	r, err := http.NewRequest("POST", payURL, bytes.NewReader(b))
	if err != nil {
		return "", err
	}

	r.Header.Set("Accept", "application/json")
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Authorization", token)

	res, err := http.DefaultClient.Do(r)
	if err != nil {
		return "", err
	}
	if res.StatusCode == http.StatusUnauthorized {
		return "", fmt.Errorf("Ключ доступа устарел или был заменен,\n пожалуйста впишите актуальный ключ в файл key.txt,\n если ключа у вас нет ключа, вы можете получить новый\n по ссылке https://qiwi.com/api \"Выпустить новый токен\"")
	}
	c := codeReq{}
	if err = json.NewDecoder(res.Body).Decode(&c); err != nil {
		return "", err
	}

	return c.Message, err
}

func detect(number string) (string, error) {

	body := fmt.Sprintf("phone=%s", number)
	r, err := http.NewRequest("POST", detectURL, strings.NewReader(body))
	if err != nil {
		return "", err
	}

	r.Header.Set("Accept", "application/json")
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(r)
	if err != nil {
		return "", err
	}

	c := codeReq{}
	if err = json.NewDecoder(res.Body).Decode(&c); err != nil {
		return "", err
	}

	return c.Message, nil
}
