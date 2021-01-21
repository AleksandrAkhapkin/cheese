package mysql

import (
	"github.com/lilekov-studio/supreme-cheese/internal/supreme-cheese/types"
	"github.com/pkg/errors"
)

func (p *MySQL) CreateUser(user *types.UserRegister) (int64, error) {

	var id int64

	tx, err := p.db.Begin()
	if err != nil {
		return 0, errors.Wrap(err, "err with Begin")
	}

	res, err := tx.Exec("INSERT INTO users (phone, email, first_name, last_name, role, sex, city, age, "+
		"unencrypted_pass, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, NOW())",
		user.Phone, user.Email, user.FirstName, user.LastName, user.Role, user.Sex, user.City, user.Age,
		user.GeneratePass)
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

func (p *MySQL) PutUserNameByID(user *types.PutCabinet) error {

	if _, err := p.db.Exec("UPDATE users SET first_name = ?, last_name = ?, updated_at = NOW() WHERE user_id = ?",
		user.FirstName, user.LastName, user.ID); err != nil {
		return err
	}

	return nil
}

func (p *MySQL) PutUserPassByID(id int, newPass string) error {

	if _, err := p.db.Exec("UPDATE users SET unencrypted_pass = ?, updated_at = NOW() WHERE user_id = ?",
		newPass, id); err != nil {
		return err
	}

	return nil
}

func (p *MySQL) UpdatePhone(phone string, id int) error {

	if _, err := p.db.Exec("UPDATE users SET phone = ?, updated_at = NOW() WHERE user_id = ?",
		phone, id); err != nil {
		return err
	}

	return nil
}

func (p *MySQL) AddPerekCard(id int, perekBool bool, perekCard string) error {

	if _, err := p.db.Exec("UPDATE users SET perek = ?, perek_card = ? WHERE user_id = ?",
		perekBool, perekCard, id); err != nil {
		return err
	}

	return nil
}

func (p *MySQL) ConfirmEmailStatus(email string) error {

	if _, err := p.db.Exec("UPDATE users SET confirm_email = true WHERE email = ?", email); err != nil {
		return err
	}

	return nil
}

func (p *MySQL) DeleteUserByID(id int) error {

	_, err := p.db.Exec("DELETE FROM users WHERE user_id = ?", id)
	if err != nil {
		return err
	}

	return nil
}

func (p *MySQL) SetUnsubscribe(email string) error {
	_, err := p.db.Exec("UPDATE users SET unsubscribe_email = true WHERE email = ?", email)
	if err != nil {
		return errors.Wrap(err, "err while Exec")
	}

	return nil
}
