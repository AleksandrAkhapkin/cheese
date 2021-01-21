package mysql

import (
	"github.com/lilekov-studio/supreme-cheese/internal/supreme-cheese/types"
	"github.com/pkg/errors"
)

func (p *MySQL) GetBoolPerekPrizeType(billID int) (bool, error) {

	var perek bool

	err := p.db.QueryRow("SELECT perek_by_user FROM bill_good_response WHERE bill_id = ?", billID).Scan(&perek)
	if err != nil {
		return false, err
	}

	return perek, nil
}

func (p *MySQL) WriteWinner(userID, billID int, prizeType, prizeStatus string) error {

	_, err := p.db.Exec("INSERT INTO winners (user_id, bill_id, prize_type, prize_status) VALUES (?, ?, ?, ?)",
		userID, billID, prizeType, prizeStatus)
	if err != nil {
		return err
	}

	return nil
}

func (p *MySQL) CountCashBackInProject(userID int) (int, error) {

	var count int
	err := p.db.QueryRow("SELECT COUNT(*) FROM winners WHERE "+
		"(user_id = ? AND prize_type = ?) OR "+
		"(user_id = ? AND prize_type = ?) OR "+
		"(user_id = ? AND prize_type = ?) OR "+
		"(user_id = ? AND prize_type = ?)",
		userID, types.PrizePhone,
		userID, types.PrizePerek,
		userID, types.PrizePhoneBIG,
		userID, types.PrizePerekBIG).Scan(&count)

	if err != nil {
		return 0, err
	}

	return count, nil
}
func (p *MySQL) GetTimeThreeCashAgo(userID int) (string, error) {

	var time string
	err := p.db.QueryRow("SELECT win_time FROM (SELECT win_time FROM winners WHERE "+
		"(user_id = ? AND prize_type = ? AND prize_status = ?) OR "+
		"(user_id = ? AND prize_type = ? AND prize_status = ?) OR "+
		"(user_id = ? AND prize_type = ? AND prize_status = ?) OR "+
		"(user_id = ? AND prize_type = ? AND prize_status = ?) OR "+
		"(user_id = ? AND prize_type = ? AND prize_status = ?) OR "+
		"(user_id = ? AND prize_type = ? AND prize_status = ?) "+
		"ORDER BY win_time DESC LIMIT 4) AS _t ORDER BY win_time ASC LIMIT 1",
		userID, types.PrizePerek, types.PrizeStatusPerekWait,
		userID, types.PrizePerekBIG, types.PrizeStatusPerekWait,
		userID, types.PrizePhone, types.PrizeStatusPhone,
		userID, types.PrizePerek, types.PrizeStatusPerek,
		userID, types.PrizePhoneBIG, types.PrizeStatusPhone,
		userID, types.PrizePerekBIG, types.PrizeStatusPerek).Scan(&time)
	if err != nil {
		return "", err
	}
	return time, nil
}

func (p *MySQL) GetBigCashUsers() ([]string, error) {

	var phone string
	phones := make([]string, 0)

	rows, err := p.db.Query("SELECT phone FROM big_cash")
	if err != nil {
		return nil, errors.Wrap(err, "err with Query in GetBigCashUsers")
	}

	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(&phone); err != nil {
			return nil, errors.Wrap(err, "err with Scan in GetBigCashUsers")
		}
		phones = append(phones, phone)
	}

	return phones, nil
}

func (p *MySQL) DelBigWinUser(phone string) error {

	_, err := p.db.Exec("DELETE FROM big_cash WHERE phone = ?", phone)
	if err != nil {
		return err
	}

	return nil
}
