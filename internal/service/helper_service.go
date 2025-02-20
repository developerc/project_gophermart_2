package service

import (
	"bytes"
	"context"
	"encoding/json"

	"net/http"

	dbstorage "github.com/developerc/project_gophermart_2/internal/db_storage"
	"github.com/developerc/project_gophermart_2/internal/general"
	"github.com/theplant/luhn"
)

type User struct {
	Name string
}

type OrderSum struct {
	Order string
	Sum   float64
}

func (s *Service) SetUserCookie(usr string) (*http.Cookie, error) {
	var cookie *http.Cookie
	u := &User{
		Name: usr,
	}

	u.Name = usr
	if encoded, err := s.secure.Encode("user", u); err == nil {
		cookie = &http.Cookie{
			Name:  "user",
			Value: encoded,
		}
		return cookie, nil
	} else {
		return nil, err
	}
}

func (s *Service) GetUserFromCookie(cookieValue string) (string, error) {
	var usr string
	u := &User{
		Name: usr,
	}
	if err := s.secure.Decode("user", cookieValue, u); err != nil {
		return "", err
	}

	return u.Name, nil
}

func (s *Service) PostUserOrders(ctx context.Context, usr string, buf bytes.Buffer) error {
	err := checkLuhna(buf.String())
	if err != nil {
		return err
	}

	if err := dbstorage.UploadOrder(ctx, s.serverSettings.GetServerSettings().DB, usr, buf.String()); err != nil {
		return err
	}
	s.chSignal.ChStart <- struct{}{}
	return nil
}

func (s *Service) GetUserOrders(ctx context.Context, usr string) ([]byte, error) {
	arrUploadedOrder, err := dbstorage.GetUserOrders(ctx, s.serverSettings.GetServerSettings().DB, usr)
	if err != nil {
		return nil, err
	}
	jsonBytes, err := json.Marshal(arrUploadedOrder)
	if err != nil {
		return nil, err
	}
	if len(arrUploadedOrder) == 0 {
		return nil, &general.ErrorNoContent{}
	}

	return jsonBytes, nil
}

func (s Service) GetUserBalance(ctx context.Context, usr string) ([]byte, error) {
	userBalance, err := dbstorage.GetUserBalance(ctx, s.serverSettings.GetServerSettings().DB, usr)
	if err != nil {
		return nil, err
	}
	jsonBytes, err := json.Marshal(userBalance)
	if err != nil {
		return nil, err
	}
	return jsonBytes, nil
}

func (s *Service) PostBalanceWithdraw(ctx context.Context, usr string, buf bytes.Buffer) error {
	var err error
	orderSum := OrderSum{}
	if err = json.Unmarshal(buf.Bytes(), &orderSum); err != nil {
		return err
	}
	err = checkLuhna(orderSum.Order)
	if err != nil {
		return err
	}
	err = dbstorage.BalanceWithdraw(ctx, s.serverSettings.GetServerSettings().DB, usr, orderSum.Order, orderSum.Sum)
	if err != nil {
		return err
	}
	return nil
}

func checkLuhna(order string) error {
	intNum := 0
	for _, runeValue := range order {
		if runeValue < 48 || runeValue > 57 {
			return &general.ErrorNumOrder{}
		}
		intNum = intNum*10 + int(runeValue-48)
	}
	if !luhn.Valid(intNum) {
		return &general.ErrorNumOrder{}
	}
	return nil
}

func (s *Service) GetUserWithdrawals(ctx context.Context, usr string) ([]byte, error) {
	arrWithdrawOrder, err := dbstorage.GetUserWithdrawals(ctx, s.serverSettings.GetServerSettings().DB, usr)
	if err != nil {
		return nil, err
	}
	jsonBytes, err := json.Marshal(arrWithdrawOrder)
	if err != nil {
		return nil, err
	}
	if len(arrWithdrawOrder) == 0 {
		return nil, &general.ErrorNoContent{}
	}

	return jsonBytes, nil
}
