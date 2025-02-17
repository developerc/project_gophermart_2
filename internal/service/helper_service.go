package service

import (
	"bytes"
	"context"
	"encoding/json"
	"time"

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
	//var usr string
	var cookie *http.Cookie
	//var err error
	u := &User{
		Name: usr,
	}

	//if len(cookieValue) == 0 {
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
	//}
	//return nil,  nil
}

func (s *Service) GetUserFromCookie(cookieValue string) (string, error) {
	var usr string
	u := &User{
		Name: usr,
	}
	if err := s.secure.Decode("user", cookieValue, u); err != nil {
		return "", err
	}
	//fmt.Println("u: ", u)

	return u.Name, nil
}

func (s *Service) PostUserOrders(ctx context.Context, usr string, buf bytes.Buffer) error {
	err := checkLuhna(buf.String())
	if err != nil {
		return err
	}

	/*ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()*/
	if err := dbstorage.UploadOrder(ctx, s.repo.GetServerSettings().DB, usr, buf.String()); err != nil {
		return err
	}
	return nil
}

func (s *Service) GetUserOrders(usr string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	arrUploadedOrder, err := dbstorage.GetUserOrders(ctx, s.repo.GetServerSettings().DB, usr)
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

func (s Service) GetUserBalance(usr string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()
	userBalance, err := dbstorage.GetUserBalance(ctx, s.repo.GetServerSettings().DB, usr)
	if err != nil {
		return nil, err
	}
	jsonBytes, err := json.Marshal(userBalance)
	if err != nil {
		return nil, err
	}
	return jsonBytes, nil
}

func (s *Service) PostBalanceWithdraw(usr string, buf bytes.Buffer) error {
	var err error
	orderSum := OrderSum{}
	if err = json.Unmarshal(buf.Bytes(), &orderSum); err != nil {
		return err
	}
	err = checkLuhna(orderSum.Order)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	err = dbstorage.BalanceWithdraw(ctx, s.repo.GetServerSettings().DB, usr, orderSum.Order, orderSum.Sum)
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

func (s *Service) GetUserWithdrawals(usr string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	arrWithdrawOrder, err := dbstorage.GetUserWithdrawals(ctx, s.repo.GetServerSettings().DB, usr)
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
