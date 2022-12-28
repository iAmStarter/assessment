//go:build unit

package expense

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

var ex = Expense{
	ID:     1,
	Title:  "strawberry smoothie",
	Amount: 100,
	Note:   "night market promotion discount 10 bath",
	Tags:   []string{"strawberry", "smoothie"},
}

func TestGetAllExpenses(t *testing.T) {

	mJson, _ := json.Marshal(&ex)

	// Arrange
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/expenses", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	mockRows := sqlmock.NewRows([]string{"id", "title", "amount", "note", "tags"}).
		AddRow(&ex.ID, &ex.Title, &ex.Amount, &ex.Note, pq.Array(&ex.Tags))

	db, mock, err := sqlmock.New()
	defer db.Close()
	mock.ExpectQuery("SELECT (.+) FROM expenses").WillReturnRows(mockRows)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	h := handler{db}
	c := e.NewContext(req, rec)
	expected := fmt.Sprintf(`[%s]`, mJson)

	// Act
	err = h.GetExpenses(c)

	// Assertions
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, expected, strings.TrimSpace(rec.Body.String()))
	}
}