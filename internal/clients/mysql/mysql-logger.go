package mysql

import (
	"github.com/lilekov-studio/supreme-cheese/internal/supreme-cheese/types"
	"github.com/pkg/errors"
)

func (p *MySQL) LogerForOS(userID int, os *types.OS) (int64, error) {

	var id int64

	tx, err := p.db.Begin()
	if err != nil {
		return 0, errors.Wrap(err, "err with Begin")
	}

	res, err := tx.Exec("INSERT INTO os_logger (user_id, contact_name, contact_email, contact_phone, text) "+
		"VALUES (?, ?, ?, ?, ?)",
		userID, os.Name, os.Email, os.Phone, os.Text)
	if err != nil {
		tx.Rollback()
		return 0, errors.Wrap(err, "err with MySQL Exec in CreateUser")
	}

	id, err = res.LastInsertId()
	if err != nil {
		tx.Rollback()
		return 0, errors.Wrap(err, "err with MySQL LastInsertId in CreateUser")
	}

	if err = tx.Commit(); err != nil {
		tx.Rollback()
		return 0, errors.Wrap(err, "err with MySQL Commit in CreateUser")
	}

	return id, nil
}

func (p *MySQL) AddLogCheckBillReq(text string, userID, billID int) error {

	_, err := p.db.Exec("INSERT INTO logger_req_checkbill (logger, user_id,  bill_id) VALUES (?, ?, ?)",
		text, userID, billID)
	if err != nil {
		return err
	}

	return nil
}

func (p *MySQL) AddLogCheckBillResp(text string, userID, billID int) error {

	_, err := p.db.Exec("INSERT INTO logger_resp_checkbill (loger, user_id,  bill_id) VALUES (?, ?, ?)",
		text, userID, billID)
	if err != nil {
		return err
	}

	return nil
}

func (p *MySQL) AddLogResponseGiftery(text string, userID, billID, reqID int, method string) error {

	_, err := p.db.Exec("INSERT INTO giftery_log_response (response, user_id, bill_id, req_id, method) "+
		"VALUES (?, ?, ?, ?, ?)", text, userID, billID, reqID, method)
	if err != nil {
		return err
	}

	return nil
}

func (p *MySQL) AddLogRequestGiftery(text string, userID, billID int, method string) (int, error) {

	var id int64

	tx, err := p.db.Begin()
	if err != nil {
		return 0, errors.Wrap(err, "err with Begin")
	}

	res, err := p.db.Exec("INSERT INTO giftery_log_req (req, user_id, bill_id, method) VALUES (?, ?, ?, ?)",
		text, userID, billID, method)
	if err != nil {
		tx.Rollback()
		return 0, errors.Wrap(err, "err with MySQL Exec in AddLogRequestGiftery")
	}

	id, err = res.LastInsertId()
	if err != nil {
		tx.Rollback()
		return 0, errors.Wrap(err, "err with MySQL LastInsertId in AddLogRequestGiftery")
	}

	if err = tx.Commit(); err != nil {
		tx.Rollback()
		return 0, errors.Wrap(err, "err with MySQL Commit in AddLogRequestGiftery")
	}

	return int(id), nil

}

func (p *MySQL) LogSMS(phone string, userID int) (int, error) {

	var id int64

	tx, err := p.db.Begin()
	if err != nil {
		return 0, errors.Wrap(err, "err with Begin")
	}

	res, err := p.db.Exec("INSERT INTO log_sms (phone, user_id) VALUES (?, ?)", phone, userID)
	if err != nil {
		tx.Rollback()
		return 0, errors.Wrap(err, "err with MySQL Exec in LogSMS")
	}

	id, err = res.LastInsertId()
	if err != nil {
		tx.Rollback()
		return 0, errors.Wrap(err, "err with MySQL LastInsertId in LogSMS")
	}

	if err = tx.Commit(); err != nil {
		tx.Rollback()
		return 0, errors.Wrap(err, "err with MySQL Commit in LogSMS")
	}

	return int(id), nil

}

func (p *MySQL) AddLogRespSMS(text string, reqID, userID int) error {

	_, err := p.db.Exec("INSERT INTO log_sms_resp (resp, req_id, user_id) VALUES (?, ?, ?)", text, reqID, userID)
	if err != nil {
		return err
	}

	return nil
}

func (p *MySQL) YoomoneyLoggerStart(logStart *types.PayRes, number, amount string, userID int) error {

	_, err := p.db.Exec("INSERT INTO yoomoney_log (user_id, type, error, request_id, status, amount, phone) "+
		"VALUES (?, ?, ?, ?, ?, ?, ?)", userID, "start", logStart.Error, logStart.Request_id, logStart.Status, amount, number)
	if err != nil {
		return err
	}

	return nil
}

func (p *MySQL) YoomoneyLoggerFinal(logFinal *types.PayedRes, number, amount string, userID int, requestID string) error {

	_, err := p.db.Exec("INSERT INTO yoomoney_log (user_id, type, error, request_id, status, amount, payment_id, invoice_id, phone) "+
		"VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		userID, "finish", logFinal.Error, requestID, logFinal.Status, amount, logFinal.Payment_id, logFinal.Invoice_id, number)
	if err != nil {
		return err
	}

	return nil
}
