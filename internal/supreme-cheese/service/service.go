package service

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/lilekov-studio/supreme-cheese/internal/clients/mysql"
	"github.com/lilekov-studio/supreme-cheese/internal/supreme-cheese/types"
	"github.com/lilekov-studio/supreme-cheese/internal/supreme-cheese/types/config"
	"github.com/lilekov-studio/supreme-cheese/pkg/infrastruct"
	"github.com/lilekov-studio/supreme-cheese/pkg/logger"
	"github.com/pkg/errors"
	"net/http"
	"strings"
)

type Service struct {
	p                   *mysql.MySQL
	secretKey           string
	email               *config.ConfigForSendEmail
	pathForBill         string
	pathForSert         string
	secretGiftery       string
	secretProverkaCheka string
	smsLogin            string
	smsPassword         string
}

func NewService(pg *mysql.MySQL, cnf *config.Config) (*Service, error) {

	return &Service{
		p:                   pg,
		secretKey:           cnf.SecretKeyJWT,
		email:               cnf.Email,
		pathForBill:         cnf.PathForBill,
		pathForSert:         cnf.PathForSert,
		secretGiftery:       cnf.GifterySecret,
		secretProverkaCheka: cnf.ProverkaChekaSecret,
		smsLogin:            cnf.SMSLogin,
		smsPassword:         cnf.SMSPassword,
	}, nil
}

func replaceSpace(user *types.UserRegister) {
	user.Phone = strings.ReplaceAll(user.Phone, " ", "")
	user.Email = strings.ReplaceAll(user.Email, " ", "")
}

func (s *Service) GetUserByID(id int) (*types.User, error) {

	user, err := s.p.GetUserByID(id)
	if err != nil {
		if err != sql.ErrNoRows {
			logger.LogError(errors.Wrap(err, "err with s.p.GetUserByID"))
			return nil, infrastruct.ErrorInternalServerError
		}
		return nil, infrastruct.ErrorJWTIsBroken
	}

	return user, nil
}

func (s *Service) GetWinners() (*types.Winners, error) {

	winners, err := s.p.GetWinners()
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with s.p.GetWinners in GetWinners"))
		return nil, infrastruct.ErrorInternalServerError
	}

	for i := range winners {
		user, err := s.GetUserByID(winners[i].ID)
		if err != nil {
			logger.LogError(errors.Wrap(err, "err with s.GetUserByID in GetWinners"))
			return nil, infrastruct.ErrorInternalServerError
		}

		winners[i].FirstName = user.FirstName
		winners[i].Email = s.replaceEmail(user.Email)
	}

	winnersGlav := types.Winners{Glav: winners}

	return &winnersGlav, nil
}

func (s *Service) GetMyBills(userID int) (*types.LcBill, error) {

	bills, err := s.p.GetMyBills(userID)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with s.p.GetWinners in GetWinners"))
		return nil, infrastruct.ErrorInternalServerError
	}

	myBills := types.LcBill{LcBill: bills}

	return &myBills, nil
}

func (s *Service) GetMyPrizes(userID int) (*types.LsPrize, error) {

	prizes, err := s.p.GetMyPrizes(userID)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with s.p.GetWinners in GetWinners"))
		return nil, infrastruct.ErrorInternalServerError
	}

	lsPrize := types.LsPrize{LsPrize: prizes}

	return &lsPrize, nil
}

func (s *Service) replaceEmail(email string) string {
	var newEmail string
	var ns int
	var chars = []byte(email)

	if ns = strings.Index(email, "@"); ns < 0 {
		newEmail = "********"
	} else {
		newEmail = string(chars[0]) + "*****" + string(chars[ns:len(email)])
	}

	return newEmail
}

func (s *Service) OS(os *types.OS) error {

	user := &types.User{}
	var err error

	if os.UserID != 0 {
		user, err = s.p.GetUserByID(os.UserID)
		if err != nil {
			logger.LogError(errors.Wrap(err, "err with s.p.GetUserByID in OS"))
			user = &types.User{}
		}
	}

	id, err := s.p.LogerForOS(user.ID, os)
	if err != nil {
		logger.LogError(errors.Wrap(err, "err with  s.p.LogerForOS"))
	}
	os.MessageID = int(id)

	if err := s.SendEmailOS(os, user); err != nil {
		return err
	}

	return nil
}

func (s *Service) findCityByURL(ip string) (string, error) {

	var chars = []byte(ip)
	ns := 0
	if ns = strings.Index(ip, ":"); ns > 0 {
		ip = string(chars[0:ns])
	}

	fmt.Println(ip)
	req, err := http.NewRequest("GET", fmt.Sprintf(types.URLForIP+ip), nil)

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return "", errors.Wrap(err, "err with client.Do")
	}
	defer resp.Body.Close()

	city := &types.IP{}
	if err = json.NewDecoder(resp.Body).Decode(&city); err != nil {
		return "default", err
	}

	return city.Name, nil
}

func (s *Service) SaveRequestLog(body []byte, route string, ip string) error {
	return s.p.SaveRequestLog(string(body), route, ip)
}
