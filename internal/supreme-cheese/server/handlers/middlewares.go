package handlers

import (
	"fmt"
	"github.com/lilekov-studio/supreme-cheese/internal/supreme-cheese/types"
	"github.com/lilekov-studio/supreme-cheese/pkg/infrastruct"
	"github.com/lilekov-studio/supreme-cheese/pkg/logger"
	"github.com/pkg/errors"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func (h *Handlers) RecoverPanic(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {

			r := recover()
			if r == nil {
				return
			}

			err := errors.New(url.QueryEscape(fmt.Sprintf("PANIC:'%v'\nRecovered in: %s", r, infrastruct.IdentifyPanic())))
			logger.LogError(err)
			apiErrorEncode(w, infrastruct.ErrorInternalServerError)
		}()

		handler.ServeHTTP(w, r)

	})
}

func (h *Handlers) AddCORS(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		(w).Header().Set("Access-Control-Allow-Origin", "*")
		(w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, x-api-token, X-api-token")

		if r.Method == "OPTIONS" {
			return
		}

		handler.ServeHTTP(w, r)
	})
}

func (h *Handlers) RequestLog(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		func() {
			req, err := httputil.DumpRequest(r, !strings.Contains(r.RequestURI, "/check/upload"))
			if err != nil {
				logger.LogError(errors.Wrap(err, "err with DumpRequest"))
				return
			}

			if err = h.srv.SaveRequestLog(req, r.RequestURI, r.RemoteAddr); err != nil {
				logger.LogInfo(url.QueryEscape(r.RequestURI))
				logger.LogInfo(url.QueryEscape(string(req)))
				logger.LogError(errors.Wrap(err, "err with SaveRequestLog"))
				return
			}
		}()

		handler.ServeHTTP(w, r)
	})
}

//
//func (h *Handlers) CheckRoleTeacher(handler http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		claims, err := infrastruct.GetClaimsByRequest(r, h.secretKey)
//		if err != nil {
//			apiErrorEncode(w, err)
//			return
//		}
//
//		if claims.Role != types.RoleTeacher {
//			apiErrorEncode(w, infrastruct.ErrorPermissionDenied)
//			return
//		}
//
//		handler.ServeHTTP(w, r)
//	})
//}
//
func (h *Handlers) CheckRoleAdmin(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, err := infrastruct.GetClaimsByRequest(r, h.secretKey)
		if err != nil {
			apiErrorEncode(w, err)
			return
		}

		if claims.Role != types.RoleAdmin {
			apiErrorEncode(w, infrastruct.ErrorPermissionDenied)
			return
		}

		handler.ServeHTTP(w, r)
	})
}

//
//func (h *Handlers) CheckRoleStudent(handler http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		claims, err := infrastruct.GetClaimsByRequest(r, h.secretKey)
//		if err != nil {
//			apiErrorEncode(w, err)
//			return
//		}
//
//		if claims.Role != types.RoleStudent {
//			apiErrorEncode(w, infrastruct.ErrorPermissionDenied)
//			return
//		}
//
//		handler.ServeHTTP(w, r)
//	})
//}

//func (h *Handlers) CheckRoleAdminAndTeacher(handler http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		claims, err := infrastruct.GetClaimsByRequest(r, h.secretKey)
//		if err != nil {
//			apiErrorEncode(w, err)
//			return
//		}
//
//		if claims.Role != types.RoleAdmin && claims.Role != types.RoleTeacher {
//			apiErrorEncode(w, infrastruct.ErrorPermissionDenied)
//			return
//		}
//
//		handler.ServeHTTP(w, r)
//	})
//}

//func (h *Handlers) RecordRequest(handler http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		claims, err := infrastruct.GetClaimsByRequest(r, h.secretKey)
//		if err != nil {
//			handler.ServeHTTP(w, r)
//			return
//		}
//
//		if err := h.srv.RecordTime(&types.RecordTime{
//			UserID:     claims.UserID,
//			RequestURL: r.RequestURI,
//		}); err != nil {
//			logger.LogError(err)
//		}
//
//		handler.ServeHTTP(w, r)
//	})
//}

func (h *Handlers) CheckUserInDBUsers(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, err := infrastruct.GetClaimsByRequest(r, h.secretKey)
		if err != nil {
			apiErrorEncode(w, err)
			return
		}
		if _, err = h.srv.GetUserByID(claims.UserID); err != nil {
			apiErrorEncode(w, err)
			return
		}

		handler.ServeHTTP(w, r)
	})
}

func (h *Handlers) CheckCanGetSertificate(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, err := infrastruct.GetClaimsByRequest(r, h.secretKey)
		if err != nil {
			apiErrorEncode(w, err)
			return
		}

		err = h.srv.CheckCanGetSertificate(claims.UserID)
		if err != nil {
			apiErrorEncode(w, err)
			return
		}

		handler.ServeHTTP(w, r)
	})
}
