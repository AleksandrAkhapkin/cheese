package mysql

import (
	"github.com/lilekov-studio/supreme-cheese/internal/supreme-cheese/types"
)

func (p MySQL) CheckDBConfURL(email string) (bool, error) {

	if err := p.db.QueryRow("SELECT email FROM conf_url WHERE email = ?", email).Scan(&email); err != nil {
		return false, err
	}

	return true, nil
}

func (p *MySQL) DeleteConfURL(email string) error {

	if _, err := p.db.Exec("DELETE FROM conf_url WHERE email = ?", email); err != nil {
		return err
	}

	return nil
}

func (p *MySQL) AddCodeForConfEmail(email, code string) error {

	if _, err := p.db.Exec("INSERT INTO conf_email (email, code) VALUES (?, ?)", email, code); err != nil {
		return err
	}

	return nil
}

func (p *MySQL) CheckCodeInConfEmail(ch *types.ConfirmationEmail) error {
	if err := p.db.QueryRow("SELECT email, code FROM conf_email WHERE email = ? AND code = ?",
		ch.Email, ch.Code).Scan(&ch.Email, &ch.Code); err != nil {
		return err
	}

	return nil
}

func (p MySQL) CheckDBConfPhone(phone string) (bool, error) {

	if err := p.db.QueryRow("SELECT phone FROM conf_phone WHERE phone = ?", phone).Scan(&phone); err != nil {
		return false, err
	}

	return true, nil
}

func (p *MySQL) DeleteConfPhone(phone string) error {

	if _, err := p.db.Exec("DELETE FROM conf_phone WHERE phone = ?", phone); err != nil {
		return err
	}

	return nil
}

func (p *MySQL) AddCodeForConfPhone(phone, code string, id int) error {

	_, err := p.db.Exec("INSERT INTO conf_phone (phone, code, user_id) VALUES (?, ?, ?)", phone, code, id)
	if err != nil {
		return err
	}

	return nil
}

func (p *MySQL) CheckCodeInConfPhone(ch *types.ConfirmationPhone) error {

	if err := p.db.QueryRow("SELECT phone, code FROM conf_phone WHERE phone = ? AND code = ?",
		ch.Phone, ch.Code).Scan(&ch.Phone, &ch.Code); err != nil {
		return err
	}
	return nil
}

func (p *MySQL) CheckConfirmationPhone(id int) (*types.CheckConfirmationPhone, error) {

	confirm := &types.CheckConfirmationPhone{}
	if err := p.db.QueryRow("SELECT confirm_phone, phone FROM users WHERE user_id = ?", id).
		Scan(&confirm.ConfirmPhone, &confirm.Phone); err != nil {
		return nil, err
	}

	return confirm, nil
}

func (p *MySQL) UpdateStatusPhone(phone string) error {

	if _, err := p.db.Exec("UPDATE users SET confirm_phone = true WHERE phone = ?", phone); err != nil {
		return err
	}

	return nil
}

func (p *MySQL) WriteTokenForEmail(id int, email string, uuid string) error {

	_, err := p.db.Exec("INSERT INTO conf_url (user_id, email, token) VALUES (?, ?, ?)", id, email, uuid)
	if err != nil {
		return err
	}

	return nil
}

func (p *MySQL) GetUserEmailByTokenInConfURL(uuid string) (string, error) {

	var email string

	err := p.db.QueryRow("SELECT email FROM conf_url WHERE token = ?", uuid).Scan(&email)
	if err != nil {
		return "", err
	}

	return email, nil
}
