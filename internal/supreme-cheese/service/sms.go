package service

import (
	"bytes"
	"encoding/json"
	"github.com/lilekov-studio/supreme-cheese/internal/supreme-cheese/types"
	"github.com/lilekov-studio/supreme-cheese/pkg/logger"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
)

func (s *Service) SendSMS(mes []types.SMSMessages) ([]byte, error) {
	smsBody := types.Sms{
		Login:    s.smsLogin,
		Password: s.smsPassword,
		Messages: mes,
	}

	jsonStr, err := json.Marshal(smsBody)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with Marshal in SendSMS"))
		return nil, err
	}

	req, err := http.NewRequest("POST", types.SMSURL, bytes.NewBuffer(jsonStr))
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with NewRequest in SendSMS"))
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	//логируем респонс
	respLogerBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with read body in ChoiceCertificate for logger"))
	}

	defer resp.Body.Close()

	return respLogerBytes, nil
}
