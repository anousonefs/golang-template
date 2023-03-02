//go:build unit

package user

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestListUsers(t *testing.T) {
	// Arrange
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/v1/users", strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	newsMockRows := sqlmock.NewRows([]string{"name", "age"}).
		AddRow("anousone", 23)

	db, mock, err := sqlmock.New()
	mock.ExpectQuery("SELECT name, age FROM users").WillReturnRows(newsMockRows)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	usecase := NewUserUsecase(db)
	h := UserHandler{
		UserUC: usecase,
	}

	c := e.NewContext(req, rec)
	expected := "[{\"name\":\"anousone\",\"age\":23,\"phone\":\"\"}]"

	// Act
	err = h.ListUsers(c)

	// Assertions
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, expected, strings.TrimSpace(rec.Body.String()))
	}
}
