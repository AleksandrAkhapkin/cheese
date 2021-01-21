package mysql

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/lilekov-studio/supreme-cheese/internal/supreme-cheese/types"
	"github.com/lilekov-studio/supreme-cheese/pkg/logger"
	"github.com/pkg/errors"
	"time"
)

type MySQL struct {
	db *sql.DB
}

func NewMySQL(dsn string) (*MySQL, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, errors.Wrap(err, "err with Open DB")
	}
	db.SetConnMaxLifetime(time.Minute * 5)
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(10)
	if err = db.Ping(); err != nil {
		return nil, errors.Wrap(err, "err with ping DB")
	}

	go func(db *sql.DB) {
		for {
			time.Sleep(30 * time.Second)
			if err = db.Ping(); err != nil {
				logger.LogError(errors.Wrap(err, "err with db.Ping in NewMySQL"))
			}
		}
	}(db)
	return &MySQL{db}, nil
}

func (p *MySQL) Close() error {
	return p.db.Close()
}

func (p *MySQL) GetUserByEmail(email string) (*types.User, error) {

	user := types.User{Email: email}
	err := p.db.QueryRow("SELECT user_id, phone, first_name, last_name, role, confirm_phone, "+
		"perek, perek_card, city, age, sex, confirm_email FROM users WHERE email = ?", email).
		Scan(&user.ID, &user.Phone, &user.FirstName, &user.LastName, &user.Role, &user.ConfirmPhone,
			&user.PerekBool, &user.PerekCard, &user.City, &user.Age, &user.Sex, &user.ConfirmEmail)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (p *MySQL) GetUserIDByPhone(phone string) (int, error) {

	var id int
	err := p.db.QueryRow("SELECT user_id FROM users WHERE phone = ?", phone).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (p *MySQL) GetUserByPhone(phone string) (*types.User, error) {

	user := types.User{Email: phone}
	err := p.db.QueryRow("SELECT user_id, email, first_name, last_name, role, confirm_phone, "+
		"perek, perek_card, city, age, sex, confirm_email FROM users WHERE phone = ?", phone).
		Scan(&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.Role, &user.ConfirmPhone,
			&user.PerekBool, &user.PerekCard, &user.City, &user.Age, &user.Sex, &user.ConfirmEmail)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (p *MySQL) GetUserByEmailHavePass(email string) (*types.UserAndPass, error) {

	user := types.UserAndPass{Email: email}
	err := p.db.QueryRow("SELECT user_id, phone, first_name, last_name, role, unencrypted_pass FROM users "+
		"WHERE email = ?", email).Scan(&user.ID, &user.Phone, &user.FirstName, &user.LastName, &user.Role, &user.Pass)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (p *MySQL) GetUserPassByID(id int) (string, error) {

	var pass string
	err := p.db.QueryRow("SELECT unencrypted_pass FROM users WHERE user_id = ?", id).Scan(&pass)
	if err != nil {
		return "", err
	}

	return pass, nil
}

func (p *MySQL) GetUserByID(id int) (*types.User, error) {

	user := types.User{ID: id}
	err := p.db.QueryRow("SELECT phone, email, first_name, last_name, role, confirm_phone, perek, perek_card, "+
		"age, sex, confirm_email FROM users WHERE user_id = ?", id).
		Scan(&user.Phone, &user.Email, &user.FirstName, &user.LastName, &user.Role, &user.ConfirmPhone, &user.PerekBool, &user.PerekCard, &user.Age, &user.Sex, &user.ConfirmEmail)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (p *MySQL) GetWinners() ([]types.Winner, error) {

	winners := make([]types.Winner, 0)

	rows, err := p.db.Query("SELECT user_id, win_time, prize_type FROM winners WHERE prize_type = ? OR prize_type = ? ORDER BY win_time ASC",
		types.PrizeTypeSert, types.PrizeTypeGlav)
	if err != nil {
		return nil, errors.Wrap(err, "err with Query in GetWinners")
	}

	defer rows.Close()
	winner := types.Winner{}

	for rows.Next() {
		if err = rows.Scan(&winner.ID, &winner.Date, &winner.Prize); err != nil {
			return nil, errors.Wrap(err, "err with Scan in GetWinners")
		}
		winners = append(winners, winner)
	}

	return winners, nil
}

func (p *MySQL) GetMyBills(userID int) ([]types.MyBill, error) {

	bills := make([]types.MyBill, 0)

	rows, err := p.db.Query("SELECT download_time, status_for_user FROM bill_good_response WHERE user_id = ? ORDER BY download_time ASC", userID)
	if err != nil {
		return nil, errors.Wrap(err, "err with Query in GetMyBills")
	}

	defer rows.Close()
	bill := types.MyBill{}

	for rows.Next() {
		if err = rows.Scan(&bill.Date, &bill.Status); err != nil {
			return nil, errors.Wrap(err, "err with Scan in GetMyBills")
		}
		bills = append(bills, bill)
	}

	return bills, nil
}

func (p *MySQL) GetMyPrizes(userID int) ([]types.MyPrize, error) {

	prizes := make([]types.MyPrize, 0)

	rows, err := p.db.Query("SELECT download_time, prize_type, prize_status FROM winners "+
		"LEFT JOIN bill_good_response USING (bill_id) "+
		"WHERE winners.user_id = ? ORDER BY winners.win_time ASC", userID)
	if err != nil {
		return nil, errors.Wrap(err, "err with Query in GetMyPrizes")
	}

	defer rows.Close()
	prize := types.MyPrize{}

	for rows.Next() {
		if err = rows.Scan(&prize.Date, &prize.Prize, &prize.Status); err != nil {
			return nil, errors.Wrap(err, "err with Scan in GetMyPrizes")
		}
		prizes = append(prizes, prize)
	}

	return prizes, nil
}

func (p *MySQL) GetUserIDByBillID(billID int) (int, error) {

	var userID int
	err := p.db.QueryRow("SELECT user_id FROM bill_good_response WHERE bill_id = ?", billID).Scan(&userID)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

func (p *MySQL) SaveRequestLog(body, route, ip string) error {

	_, err := p.db.Exec("INSERT INTO request_log (dump_request, ip, route) VALUES (?, ?, ?)", body, ip, route)
	if err != nil {
		return err
	}

	return nil
}

func (p *MySQL) GetArrayConfirmEmail() ([]string, error) {

	users := make([]string, 0)

	rows, err := p.db.Query("SELECT email FROM forsend")
	if err != nil {
		return nil, errors.Wrap(err, "err with Query in GetMyPrizes")
	}

	defer rows.Close()
	var user string

	for rows.Next() {
		if err = rows.Scan(&user); err != nil {
			return nil, errors.Wrap(err, "err with Scan in GetMyPrizes")
		}
		users = append(users, user)
	}

	return users, nil
}
