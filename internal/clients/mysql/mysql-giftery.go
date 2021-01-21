package mysql

import "github.com/lilekov-studio/supreme-cheese/internal/supreme-cheese/types"

func (p *MySQL) GetWinnerCertificateByUserID(userID int) (int, error) {

	var billID int

	err := p.db.QueryRow("SELECT bill_id FROM winners WHERE user_id = ? AND prize_type = ?",
		userID, types.PrizeTypeSert).Scan(&billID)
	if err != nil {
		return 0, err
	}

	return billID, nil
}
