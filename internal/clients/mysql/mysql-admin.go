package mysql

import (
	"github.com/lilekov-studio/supreme-cheese/internal/supreme-cheese/types"
	"github.com/pkg/errors"
)

func (p *MySQL) AddCheeseName(name string) error {

	var err error
	_, err = p.db.Exec("INSERT INTO product_name (name) VALUES (?)", name)
	if err != nil {
		return err
	}

	return nil
}

func (p *MySQL) GetArrayNotFindCheeseBills() ([]int, error) {

	billsID := make([]int, 0)
	rows, err := p.db.Query("SELECT bill_id FROM bill_good_response WHERE status = ?", types.InnerStatusBillNotFindFACKINGCheese)
	if err != nil {
		return nil, errors.Wrap(err, "err with Query in GetArrayNotFindCheeseBills")
	}

	defer rows.Close()

	var billID int

	for rows.Next() {
		if err = rows.Scan(&billID); err != nil {
			return nil, errors.Wrap(err, "err with Scan in GetArrayNotFindCheeseBills")
		}
		billsID = append(billsID, billID)
	}

	return billsID, nil
}

func (p *MySQL) GetPositionByBillID(billID int) ([]types.Items, error) {

	positions := make([]types.Items, 0)
	rows, err := p.db.Query("SELECT name, price, count, sum FROM position_in_bill WHERE bill_id = ?", billID)
	if err != nil {
		return nil, errors.Wrap(err, "err with Query in GetPositionByBillID")
	}

	defer rows.Close()

	position := types.Items{}

	for rows.Next() {
		if err = rows.Scan(&position.Name, &position.Price, &position.Count, &position.Sum); err != nil {
			return nil, errors.Wrap(err, "err with Scan in GetPositionByBillID")
		}
		positions = append(positions, position)
	}

	return positions, nil
}

//for scripts
//////////////////////////////////////////////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////////////////////////////////////////////
func (p *MySQL) GetWinnersCertificate() ([]types.WaitCertificate, error) {

	waitCertificates := make([]types.WaitCertificate, 0)
	rows, err := p.db.Query("SELECT user_id, prize_status FROM winners WHERE prize_type = ? ORDER BY win_time ASC", types.PrizeTypeSert)
	if err != nil {
		return nil, errors.Wrap(err, "err with Query in GetPositionByBillID")
	}

	defer rows.Close()

	waitCertificate := types.WaitCertificate{}

	for rows.Next() {
		if err = rows.Scan(&waitCertificate.UserID, &waitCertificate.NameCertificate); err != nil {
			return nil, errors.Wrap(err, "err with Scan in GetPositionByBillID")
		}
		waitCertificates = append(waitCertificates, waitCertificate)
	}

	return waitCertificates, nil

}

func (p *MySQL) GetArrayVALIDForScript() ([]int, error) {

	billsID := make([]int, 0)
	rows, err := p.db.Query("SELECT bill_id FROM bill_good_response WHERE status = ?", types.StatusBillValid)
	if err != nil {
		return nil, errors.Wrap(err, "err with Query in GetArrayNotFindCheeseBills")
	}

	defer rows.Close()

	var billID int

	for rows.Next() {
		if err = rows.Scan(&billID); err != nil {
			return nil, errors.Wrap(err, "err with Scan in GetArrayNotFindCheeseBills")
		}
		billsID = append(billsID, billID)
	}

	return billsID, nil
}
