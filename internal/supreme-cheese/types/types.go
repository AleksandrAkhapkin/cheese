package types

import (
	"io"
)

const (
	RoleUser  = "user"
	RoleAdmin = "admin"

	SMSURL       = "https://api.iqsms.ru/messages/v2/send.json"
	CheckBillURL = "https://proverkacheka.com/api/v1/check/get"
	HomePage     = "http://supreme-cheese.ru"

	StatusBillWait                      = "wait"
	StatusBillValid                     = "valid"
	StatusBillInvalid                   = "invalid"
	StatusBillDouble                    = "double"
	InnerStatusBillNotCanScan           = "not can scan"
	InnerStatusBillErrorInHandWrite     = "error hand write"
	InnerStatusBillNotFindFACKINGCheese = "not find cheese"
	InnerStatusBillScanButInvalid       = "scan is good, but not valid"
	InnerStatusBillInvalidTimeBill      = "error time bill"

	PrizeTypeSert           = "Сертификат"
	PrizeTypeGlav           = "Главный приз"
	PrizePhone              = "30р."
	PrizePhoneBIG           = "60р."
	PrizePerek              = "300 баллов"
	PrizePerekBIG           = "600 баллов"
	PrizeStatusPhone        = "Зачислен"
	PrizeStatusWaitSend     = "Ожидает отправки"
	PrizeStatusPerek        = "Ожидает выплаты"
	PrizeStatusPerekWait    = "В обработке"
	PrizeStatusSertChoice   = "Получить приз"
	PrizeStatusAccept       = "Ожидает отправки на email"
	PrizeStatusSendEmail    = "Отправлен на email"
	PrizeStatusLimitInDay   = "Вы превысили кол-во выплат за день"
	PrizeStatusLimitInPromo = "Вы превысили кол-во выплат за акцию"
	TokenForAuthByURL       = "?token="
	URLForIP                = "http://api.sypexgeo.net/json/"
	LimitForSendSMS         = 20
	DaysLimitCashBack       = 3
	ProjectLimitCashBack    = 30

	CashBackSum    = 30
	CashBackSumBIG = 60

	URLProverkaCheka = "https://proverkacheka.com/check/"

	GifteryNominal                  = 3000
	GifteryTypeMethodMakeOrder      = "Make order"
	GifteryTypeMethodGetCertificate = "Get Certificate"

	WinDayStatusYes    = "yes"
	WinDayStatusNo     = "no"
	WinDayStatusDouble = "double"
	WinDayStatusWait   = "wait"
)

////////////////////////////////// users ///////////////////////////////////

type UserRegister struct {
	ID           int    `json:"id"`
	GeneratePass string `json:"generate_pass"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Phone        string `json:"phone"`
	Email        string `json:"email"`
	Role         string `json:"role"`
	City         string `json:"city"`
	Sex          string `json:"sex"`
	Age          int    `json:"age"`
}

type User struct {
	ID           int    `json:"id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Phone        string `json:"phone"`
	Email        string `json:"email"`
	Role         string `json:"role"`
	ConfirmPhone bool   `json:"confirm_phone"`
	PerekBool    bool   `json:"perek_bool"`
	PerekCard    string `json:"perek_card"`
	City         string `json:"city"`
	Age          int    `json:"age"`
	Sex          string `json:"sex"`
	ConfirmEmail bool   `json:"confirm_email"`
}

type UserAndPass struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	Pass      string `json:"pass"`
}

////////////////////////////////// service ///////////////////////////////////

type Token struct {
	Token string `json:"token"`
}

type RecoverPass struct {
	GeneratePass string `json:"generate_pass"`
	Email        string `json:"email"`
}

type PutCabinet struct {
	ID        int    `json:"id"`
	NewPhone  string `json:"phone"`
	OldPass   string `json:"old_pass"`
	NewPass   string `json:"new_pass"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type ConfirmationEmail struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}
type ConfirmationPhone struct {
	Phone string `json:"phone"`
	Code  string `json:"code"`
}

type CheckConfirmationPhone struct {
	Phone        string `json:"phone"`
	ConfirmPhone bool   `json:"confirm"`
}

type CheckCode struct {
	Code bool `json:"code"`
}

type Authorize struct {
	Email string `json:"email"`
	Pass  string `json:"pass"`
}

type RegisterPhone struct {
	Phone string `json:"phone"`
}

type ConfirmUsers struct {
	Email string
}

type UploadFile struct {
	UserID          int
	Body            io.Reader
	PerecrestokBool bool   `json:"perek_bool"`
	PerecrestokCard string `json:"perek_card"`
}

////////////////////////////////// bill ///////////////////////////////////

type AddName struct {
	Name string `json:"name"`
}
type InvalidBill struct {
	BillID int    `json:"bill_id"`
	Text   string `json:"text"`
}

type HandWriteBill struct {
	BillID int    `json:"bill_id"`
	FN     string `json:"fn"`
	FD     string `json:"fd"`
	FP     string `json:"fp"`
	Date   string `json:"time"`
	Sum    string `json:"sum"`
}

type CheckBillReq struct {
	Code int          `json:"code"`
	Data HuetaForBill `json:"data"`
}

type HuetaForBill struct {
	JSON CheckBillDataReq `json:"json"`
}

type CheckBillDataReq struct {
	CodeOperation           int     `json:"code"`
	Shop                    string  `json:"user"`
	Items                   []Items `json:"items"`
	FnsUrl                  string  `json:"fnsUrl"`
	SellerAddress           string  `json:"sellerAddress"`
	KktRegId                string  `json:"kktRegId"`
	RetailPlace             string  `json:"retail_place"`
	RetailPlaceAddress      string  `json:"retailPlaceAddress"`
	ShopInn                 string  `json:"userInn"`
	Date                    string  `json:"dateTime"`
	CheckNumber             int     `json:"requestNumber"`
	TotalSum                int64   `json:"totalSum"`
	ShiftNumber             int     `json:"shiftNumber"`
	OperationType           int     `json:"operationType"`
	FiscalDriveNumber       string  `json:"fiscalDriveNumber"`
	FiscalDocumentNumber    int64   `json:"fiscalDocumentNumber"`
	FiscalSign              int64   `json:"fiscalSign"`
	FiscalDocumentFormatVer int     `json:"fiscal_document_format_ver"`
}

type Items struct {
	Sum   float64 `json:"sum"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
	Count float64 `json:"quantity"`
}

////////////////////////////////// API ///////////////////////////////////

type Winners struct {
	Glav []Winner `json:"glav"`
}

type Winner struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	Prize     string `json:"prize"`
	Email     string `json:"email"`
	Date      string `json:"date"`
}

type LcBill struct {
	LcBill []MyBill `json:"lc_bill"`
}

type LsPrize struct {
	LsPrize []MyPrize `json:"ls_prize"`
}

type MyBill struct {
	Date   string `json:"date"`
	Status string `json:"status"`
}

type MyPrize struct {
	Date   string `json:"date"`
	Prize  string `json:"prize"`
	Status string `json:"status"`
}

type OS struct {
	UserID    int    `json:"user_id"`
	MessageID int    `json:"message_id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Text      string `json:"text"`
}

////////////////////////////////// other service ///////////////////////////////////

type IP struct {
	NameCity `json:"city"`
}

type NameCity struct {
	Name string `json:"name_ru"`
}

type SendCashback struct {
	Phone  string `json:"phone"`
	Amount int    `json:"amount"`
}

type ChoiceCertificate struct {
	ProductID int `json:"product_id"`
}

type Giftery struct {
	Status string `json:"status"`
	GifteryError
	GifteryOK
}

type GifteryOK struct {
	Data GifteryOKData `json:"data"`
}
type GifteryOKData struct {
	ID          int    `json:"id"`
	Certificate string `json:"certificate"`
}

type GifteryError struct {
	Error GifteryErrorInsert `json:"error"`
}
type GifteryErrorInsert struct {
	Code      int    `json:"code"`
	TextError string `json:"text"`
}

type WaitCertificate struct {
	UserID          int
	NameCertificate string
}

type TimeBill struct {
	TimeBuy    string
	TimeUpload string
	BillID     int
	UserID     int
}

type Sms struct {
	Messages []SMSMessages `json:"messages"`
	Login    string        `json:"login"`
	Password string        `json:"password"`
}

type SMSMessages struct {
	Phone    string `json:"phone"`
	ClientID int    `json:"clientId"`
	Text     string `json:"text"`
}
