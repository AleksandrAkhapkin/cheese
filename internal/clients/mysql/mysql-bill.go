package mysql

import (
	"fmt"
	"github.com/lilekov-studio/supreme-cheese/internal/supreme-cheese/types"
	"github.com/pkg/errors"
)

func (p *MySQL) WriteBillByUser(userID int, perekByUser bool, perekCard string) (int64, error) {

	var billID int64

	tx, err := p.db.Begin()
	if err != nil {
		return 0, errors.Wrap(err, "err with Begin")
	}

	res, err := tx.Exec("INSERT INTO bill_good_response (user_id, status, status_for_user, perek_by_user, perek_card) "+
		"VALUES (?, ?, ?, ?, ?)", userID, types.StatusBillWait, types.StatusBillWait, perekByUser, perekCard)
	if err != nil {
		tx.Rollback()
		return 0, errors.Wrap(err, "err with MySQL Exec in WriteBillByUser")
	}

	billID, err = res.LastInsertId()
	if err != nil {
		tx.Rollback()
		return 0, errors.Wrap(err, "err with MySQL LastInsertId in WriteBillByUser")
	}

	if err = tx.Commit(); err != nil {
		tx.Rollback()
		return 0, errors.Wrap(err, "err with MySQL Commit in WriteBillByUser")
	}

	return billID, nil
}

func (p *MySQL) WriteGoodResponseBill(userID, billID int, bill *types.CheckBillReq) error {

	var url string
	url = fmt.Sprintf("%s%s-%d-%d", types.URLProverkaCheka, bill.Data.JSON.FiscalDriveNumber, bill.Data.JSON.FiscalDocumentNumber, bill.Data.JSON.FiscalSign)
	_, err := p.db.Exec("UPDATE bill_good_response SET code_operation = ?, shop = ?, fns_url = ?, "+
		"seller_address = ?, kkt_reg_id = ?, retail_place = ?, retail_place_address = ?, shop_inn = ?, "+
		"date_check = ?, check_number = ?, total_sum_cop = ?, shift_number = ?, operation_type = ?, drive_num = ?, "+
		"doc_num = ?, fiscal_sign = ?, fiscal_doc_format = ?, url = ? "+
		"WHERE user_id = ? AND bill_id = ?",
		bill.Data.JSON.CodeOperation,
		bill.Data.JSON.Shop,
		bill.Data.JSON.FnsUrl,
		bill.Data.JSON.SellerAddress,
		bill.Data.JSON.KktRegId,
		bill.Data.JSON.RetailPlace,
		bill.Data.JSON.RetailPlaceAddress,
		bill.Data.JSON.ShopInn,
		bill.Data.JSON.Date,
		bill.Data.JSON.CheckNumber,
		bill.Data.JSON.TotalSum,
		bill.Data.JSON.ShiftNumber,
		bill.Data.JSON.OperationType,
		bill.Data.JSON.FiscalDriveNumber,
		bill.Data.JSON.FiscalDocumentNumber,
		bill.Data.JSON.FiscalSign,
		bill.Data.JSON.FiscalDocumentFormatVer,
		url,

		userID,
		billID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (p *MySQL) UpdateInnerStatusBill(billID int, status string) error {

	if _, err := p.db.Exec("UPDATE bill_good_response SET status = ? WHERE bill_id = ?", status, billID); err != nil {
		return err
	}

	return nil
}

func (p *MySQL) UpdateStatusBillForUser(billID int, status string) error {

	if _, err := p.db.Exec("UPDATE bill_good_response SET status_for_user = ?, check_time = NOW() WHERE bill_id = ?", status, billID); err != nil {
		return err
	}

	return nil
}

func (p *MySQL) WriteBillPosition(userID, billID int, arrayPos []types.Items) error {

	for i, _ := range arrayPos {
		if _, err := p.db.Exec("INSERT INTO position_in_bill (bill_id, user_id, name, price, count, sum) "+
			"VALUES (?, ?, ?, ?, ?, ?)",
			billID, userID, arrayPos[i].Name, arrayPos[i].Price, arrayPos[i].Count, arrayPos[i].Sum); err != nil {
			return err
		}
	}
	return nil
}

func (p *MySQL) FindBillByFdFpFnInGoodResponse(fd, fp int64, fn string) (*types.CheckBillDataReq, error) {

	bill := types.CheckBillDataReq{}
	err := p.db.QueryRow("SELECT code_operation, shop, fns_url, seller_address, kkt_reg_id, "+
		"retail_place, retail_place_address, shop_inn, date_check, check_number, "+
		"total_sum_cop, shift_number, operation_type, drive_num, doc_num, fiscal_sign, fiscal_doc_format FROM bill_good_response "+
		"WHERE doc_num = ? AND fiscal_sign = ? AND  drive_num = ?", fd, fp, fn).Scan(
		&bill.CodeOperation,
		&bill.Shop,
		&bill.FnsUrl,
		&bill.SellerAddress,
		&bill.KktRegId,
		&bill.RetailPlace,
		&bill.RetailPlaceAddress,
		&bill.ShopInn,
		&bill.Date,
		&bill.CheckNumber,
		&bill.TotalSum,
		&bill.ShiftNumber,
		&bill.OperationType,
		&bill.FiscalDriveNumber,
		&bill.FiscalDocumentNumber,
		&bill.FiscalSign,
		&bill.FiscalDocumentFormatVer,
	)

	if err != nil {
		return nil, err
	}
	return &bill, nil
}

func (p *MySQL) CountBillByFdFpFnForAdmin(fd, fp int64, fn string) (int, error) {

	var count int
	err := p.db.QueryRow("SELECT COUNT(*) FROM bill_good_response WHERE doc_num = ? AND fiscal_sign = ? AND  drive_num = ?", fd, fp, fn).Scan(&count)

	if err != nil {
		return 0, err
	}
	return count, nil
}

func (p *MySQL) GetValidTime() ([]types.TimeBill, error) {

	bills := make([]types.TimeBill, 0)

	rows, err := p.db.Query("SELECT bill_id, user_id, date_check, download_time FROM bill_good_response WHERE status = 'valid'")
	if err != nil {
		return nil, errors.Wrap(err, "err with Query in GetMyPrizes")
	}

	defer rows.Close()
	var bill types.TimeBill

	for rows.Next() {
		if err = rows.Scan(&bill.BillID, &bill.UserID, &bill.TimeBuy, &bill.TimeUpload); err != nil {
			return nil, errors.Wrap(err, "err with Scan in GetMyPrizes")
		}
		bills = append(bills, bill)
	}
	return bills, nil
}
