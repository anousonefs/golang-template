//go:build integration

package user

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"

	config "github.com/anousoneFS/clean-architecture/config"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

const serverPort = 80

func TestITListUsers(t *testing.T) {
	// Setup server
	eh := echo.New()
	go func(e *echo.Echo) {
		cfg, err := config.NewConfig()
		if err != nil {
			log.Fatal(err)
		}
		db, err := sql.Open(cfg.DBDriver, cfg.DSNInfo())
		if err != nil {
			log.Fatal(err)
		}
		usecase := NewUserUsecase(db)
		h := UserHandler{
			UserUC: usecase,
		}
		e.GET("/v1/users", h.ListUsers)
		e.Start(fmt.Sprintf(":%d", serverPort))
	}(eh)

	for {
		conn, err := net.DialTimeout("tcp", fmt.Sprintf("localhost:%d", serverPort), 30*time.Second)
		if err != nil {
			log.Println(err)
		}
		if conn != nil {
			conn.Close()
			break
		}
	}

	// Arrange
	reqBody := ``
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:%d/v1/users", serverPort), strings.NewReader(reqBody))
	assert.NoError(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	client := http.Client{}

	// Act
	resp, err := client.Do(req)
	assert.NoError(t, err)

	byteBody, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	resp.Body.Close()

	// Assertions
	expected := "[{\"name\":\"anousone\",\"age\":23,\"phone\":\"\"}]"

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, expected, strings.TrimSpace(string(byteBody)))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = eh.Shutdown(ctx)
	assert.NoError(t, err)
}
