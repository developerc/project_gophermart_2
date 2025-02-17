package service

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"time"

	//"fmt"
	"log"
	"net/http"

	"github.com/developerc/project_gophermart_2/internal/config"
	dbstorage "github.com/developerc/project_gophermart_2/internal/db_storage"
	"github.com/developerc/project_gophermart_2/internal/loyalty"
	"github.com/gorilla/securecookie"
)

type repository interface {
	Register(ctx context.Context, buf bytes.Buffer) (*http.Cookie, error)
	UserLogin(ctx context.Context, buf bytes.Buffer) (*http.Cookie, error)
	GetUserFromCookie(cookieValue string) (string, error)
	GetServerSettings() *config.ServerSettings
	PostUserOrders(ctx context.Context, usr string, buf bytes.Buffer) error
	GetUserOrders(ctx context.Context, usr string) ([]byte, error)
	GetUserBalance(ctx context.Context, usr string) ([]byte, error)
	PostBalanceWithdraw(usr string, buf bytes.Buffer) error
	GetUserWithdrawals(usr string) ([]byte, error)
}

type Service struct {
	repo   repository
	secure *securecookie.SecureCookie
}

type LgnPsw struct {
	Lgn string `json:"login"`
	Psw string `json:"password"`
}

func (s *Service) Register(ctx context.Context, buf bytes.Buffer) (*http.Cookie, error) {
	var err error
	lgnPsw := LgnPsw{}
	if err = json.Unmarshal(buf.Bytes(), &lgnPsw); err != nil {
		return nil, err
	}
	log.Println("from Register:", lgnPsw)

	/*ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()*/
	if err = dbstorage.InsertUser(ctx, s.repo.GetServerSettings().DB, lgnPsw.Lgn, lgnPsw.Psw); err != nil {
		return nil, err
	}
	cookie, err := s.SetUserCookie(lgnPsw.Lgn)
	if err != nil {
		return nil, err
	}
	return cookie, nil
}

func (s *Service) UserLogin(ctx context.Context, buf bytes.Buffer) (*http.Cookie, error) {
	var err error
	lgnPsw := LgnPsw{}
	if err = json.Unmarshal(buf.Bytes(), &lgnPsw); err != nil {
		return nil, err
	}

	/*ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()*/
	if err = dbstorage.CheckLgnPsw(ctx, s.repo.GetServerSettings().DB, lgnPsw.Lgn, lgnPsw.Psw); err != nil {
		return nil, err
	}
	cookie, err := s.SetUserCookie(lgnPsw.Lgn)
	if err != nil {
		return nil, err
	}
	return cookie, nil
}

func (s *Service) GetAdresRun() string {
	return s.repo.GetServerSettings().AdresRun
}

func NewService() (*Service, error) {

	serverSettings, err := config.InitServerSettings()
	if err != nil {
		log.Println(err)
	}
	service := Service{repo: serverSettings}
	serverSettings.DB, err = sql.Open("pgx", serverSettings.AdresBase)
	if err != nil {
		return nil, err
	}

	const duration uint = 20
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(duration)*time.Second)
	defer cancel()
	if err := dbstorage.CreateTables(ctx, serverSettings.DB); err != nil {
		return nil, err
	}
	service.InitSecure()
	loyalty.RunLoyalty(serverSettings.DB, serverSettings.AdresAccrual)
	return &service, nil
}

func (s *Service) InitSecure() {
	var hashKey = []byte(s.repo.GetServerSettings().SecretCookies)
	var blockKey = []byte("a-lot-secret-qwe")
	s.secure = securecookie.New(hashKey, blockKey)
}
