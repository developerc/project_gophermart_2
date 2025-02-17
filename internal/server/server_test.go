package server

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/developerc/project_gophermart_2/internal/service"
	"github.com/stretchr/testify/require"
)

type LgnPsw struct {
	Lgn string `json:"login"`
	Psw string `json:"password"`
}

func TestApi(t *testing.T) {
	service, err := service.NewService()
	require.NoError(t, err)
	server, err := NewServer(service)
	require.NoError(t, err)
	tsrv := httptest.NewServer(server.SetupRoutes())
	defer tsrv.Close()
	var lgnPsv = LgnPsw{}
	lgnPsv.Lgn = "mylogin"
	lgnPsv.Psw = "mypassword"
	var buf bytes.Buffer
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	t.Run("#1_RegisterTest", func(t *testing.T) {
		jsonBytes, err := json.Marshal(lgnPsv)
		require.NoError(t, err)
		buf = *bytes.NewBuffer(jsonBytes)
		_, err = service.Register(ctx, buf)
		require.NoError(t, err)
	})
	t.Run("#2_LoginTest", func(t *testing.T) {
		_, err = service.UserLogin(buf)
		require.NoError(t, err)
	})
	t.Run("#3_PostUserOrdersTest", func(t *testing.T) {
		buf = *bytes.NewBuffer([]byte("12345678903"))
		err = service.PostUserOrders(ctx, lgnPsv.Lgn, buf)
		require.NoError(t, err)
	})
}
