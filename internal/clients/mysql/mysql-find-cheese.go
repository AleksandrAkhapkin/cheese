package mysql

import (
	"github.com/pkg/errors"
)

func (p *MySQL) GetNamesOfProduct() ([]string, error) {

	names := make([]string, 0)

	rows, err := p.db.Query("SELECT name FROM product_name")
	if err != nil {
		return nil, errors.Wrap(err, "err with Query in GetNamesOfProduct")
	}

	defer rows.Close()
	var name string

	for rows.Next() {
		if err = rows.Scan(&name); err != nil {
			return nil, errors.Wrap(err, "err with Scan in GetMyPrizes")
		}
		names = append(names, name)
	}

	return names, nil
}

func (p *MySQL) MarkerCheesePosition(billID int, name string) error {
	_, err := p.db.Exec("UPDATE position_in_bill SET marker = ? WHERE bill_id = ? AND name = ?",
		"FIND", billID, name)
	if err != nil {
		return err
	}
	return nil
}
