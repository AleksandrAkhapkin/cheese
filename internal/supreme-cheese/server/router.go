package server

import (
	"github.com/gorilla/mux"
	"github.com/lilekov-studio/supreme-cheese/internal/supreme-cheese/server/handlers"
	"net/http"
)

func NewRouter(h *handlers.Handlers) *mux.Router {

	router := mux.NewRouter().StrictSlash(true)
	router.Use(h.RecoverPanic)
	router.Use(h.RequestLog)
	router.Use(h.AddCORS)

	authRouter := router.PathPrefix("").Subrouter()
	authRouter.Use(h.CheckUserInDBUsers)

	adminRouter := router.PathPrefix("").Subrouter()
	adminRouter.Use(h.CheckUserInDBUsers)
	adminRouter.Use(h.CheckRoleAdmin)

	winRouter := router.PathPrefix("").Subrouter()
	winRouter.Use(h.CheckUserInDBUsers)
	winRouter.Use(h.CheckCanGetSertificate)

	router.Methods(http.MethodOptions, http.MethodGet).Path("/mail/unsubscribe/{email}").HandlerFunc(h.UnsubscribeMail)

	router.Methods(http.MethodOptions, http.MethodGet).Path("/ping").HandlerFunc(h.Ping)

	router.Methods(http.MethodOptions, http.MethodPost).Path("/register").HandlerFunc(h.RegisterUser)
	router.Methods(http.MethodOptions, http.MethodPost).Path("/auth").HandlerFunc(h.Auth)
	router.Methods(http.MethodOptions, http.MethodPost).Path("/auth/token").HandlerFunc(h.AuthByToken)

	router.Methods(http.MethodOptions, http.MethodPost).Path("/recoverpassword").HandlerFunc(h.RecoverPassword)

	authRouter.Methods(http.MethodOptions, http.MethodGet).Path("/cabinet").HandlerFunc(h.GetCabinetByClaims)
	authRouter.Methods(http.MethodOptions, http.MethodPut).Path("/cabinet").HandlerFunc(h.PutCabinetByClaims)

	authRouter.Methods(http.MethodOptions, http.MethodGet).Path("/check/verifyphone").HandlerFunc(h.VerifyPhone)
	authRouter.Methods(http.MethodOptions, http.MethodPost).Path("/register/phone").HandlerFunc(h.RegisterPhone)
	authRouter.Methods(http.MethodOptions, http.MethodPost).Path("/register/phone/confirmation").HandlerFunc(h.ConfirmationPhone)

	authRouter.Methods(http.MethodOptions, http.MethodPost).Path("/check/upload").HandlerFunc(h.UploadBill)

	adminRouter.Methods(http.MethodOptions, http.MethodPost).Path("/check/a/write").HandlerFunc(h.WriteBill)
	adminRouter.Methods(http.MethodOptions, http.MethodPost).Path("/check/a/addname").HandlerFunc(h.AddCheeseName)
	adminRouter.Methods(http.MethodOptions, http.MethodPost).Path("/check/a/invalid").HandlerFunc(h.InvalidBill)
	adminRouter.Methods(http.MethodOptions, http.MethodPost).Path("/check/a/cash").HandlerFunc(h.SendCashback)

	router.Methods(http.MethodOptions, http.MethodPost).Path("/os").HandlerFunc(h.OS)

	router.Methods(http.MethodOptions, http.MethodGet).Path("/winners").HandlerFunc(h.GetWinners)
	router.Methods(http.MethodOptions, http.MethodGet).Path("/cabinet/mybill").HandlerFunc(h.GetMyBills)
	router.Methods(http.MethodOptions, http.MethodGet).Path("/cabinet/myprize").HandlerFunc(h.GetMyPrize)

	winRouter.Methods(http.MethodOptions, http.MethodPost).Path("/certificate").HandlerFunc(h.ChoiceCertificate)

	return router
}
