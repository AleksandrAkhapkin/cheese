package mysql

import "github.com/lilekov-studio/supreme-cheese/internal/supreme-cheese/types"

func (p *MySQL) CountSMSByUser(userID int) (int, error) {

	var count int
	err := p.db.QueryRow("SELECT COUNT(*) FROM log_sms WHERE user_id = ?", userID).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (p *MySQL) UpdateWinnerCertStatus(billID int, status string) error {

	if _, err := p.db.Exec("UPDATE winners SET prize_status = ? WHERE bill_id = ? AND prize_type = ?", status, billID, types.PrizeTypeSert); err != nil {
		return err
	}

	return nil
}

func (p *MySQL) UpdateWinnerCashBackStatus(billID int, statusPrize string) error {

	if _, err := p.db.Exec("UPDATE winners SET prize_status = ? WHERE "+
		"(bill_id = ? AND prize_type = ?) OR "+
		"(bill_id = ? AND prize_type = ?) OR "+
		"(bill_id = ? AND prize_type = ?) OR "+
		"(bill_id = ? AND prize_type = ?)",
		statusPrize,
		billID, types.PrizePhone,
		billID, types.PrizePerek,
		billID, types.PrizePhoneBIG,
		billID, types.PrizePerekBIG); err != nil {
		return err
	}

	return nil
}

func (p *MySQL) FindCertificateWinnerHaveChoiceCertificate(userID int) (int, error) {

	var bill_id int
	err := p.db.QueryRow("SELECT bill_id FROM winners WHERE user_id = ? AND prize_type = ? AND prize_status = ?",
		userID, types.PrizeTypeSert, types.PrizeStatusSertChoice).Scan(&bill_id)
	if err != nil {
		return 0, err
	}

	return bill_id, nil
}

func (p *MySQL) HaveWinCertificate(userID int) (bool, error) {

	var trash string
	err := p.db.QueryRow("SELECT prize_status FROM winners WHERE user_id = ? AND prize_type = ? ",
		userID, types.PrizeTypeSert).Scan(&trash)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (p *MySQL) PUTWinDayStatus(billID int, statusPrize string) error {

	if _, err := p.db.Exec("UPDATE bill_good_response SET day_win = ? WHERE bill_id = ?",
		statusPrize, billID); err != nil {
		return err
	}

	return nil
}

func (p *MySQL) CashBackByBillID(billID int) error {

	var trash int
	err := p.db.QueryRow("SELECT user_id FROM winners WHERE bill_id = ? ", billID).Scan(&trash)
	if err != nil {
		return err
	}

	return nil
}
