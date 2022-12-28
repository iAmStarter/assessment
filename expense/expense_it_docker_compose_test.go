//go:build itdockercompose

package expense

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/stretchr/testify/assert"
	"io"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestCreateExpense(t *testing.T) {
	eh, err := setup()

	body := bytes.NewBufferString(`{
		"title": "strawberry smoothie",
		"amount": 79,
		"note": "night market promotion discount 10 bath", 
		"tags": ["food", "beverage"]
	}`)

	var e Expense
	res := request(http.MethodPost, uri("expenses"), body)
	err = res.Decode(&e)

	if assert.NoError(t, err) {
		assert.EqualValues(t, http.StatusCreated, res.StatusCode)
		assert.NotEqual(t, 0, e.ID)
		assert.Equal(t, "strawberry smoothie", e.Title)
		assert.Equal(t, "night market promotion discount 10 bath", e.Note)
		assert.Equal(t, float32(79), e.Amount)
		assert.Equal(t, []string{"food", "beverage"}, e.Tags)
	}

	err = shutdown(eh)
	assert.NoError(t, err)
}

func TestUpdateExpense(t *testing.T) {
	eh, err := setup()

	e := seedExpense(t)

	body := bytes.NewBufferString(`{
		"title": "apple smoothie",
		"amount": 89,
		"note": "no discount",
		"tags": ["beverage"]
	}`)
	res := request(http.MethodPut, uri("expenses", strconv.Itoa((e.ID))), body)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	var update Expense
	res = request(http.MethodGet, uri("expenses", strconv.Itoa((e.ID))), nil)
	err = res.Decode(&update)

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, e.ID, update.ID)
		assert.Equal(t, "apple smoothie", update.Title)
		assert.Equal(t, "no discount", update.Note)
		assert.Equal(t, float32(89), update.Amount)
		assert.Equal(t, []string{"beverage"}, update.Tags)
	}

	err = shutdown(eh)
	assert.NoError(t, err)
}

func TestGetAllExpenses(t *testing.T) {
	eh, err := setup()
	seedExpense(t)
	var es []Expense
	res := request(http.MethodGet, uri("expenses"), nil)
	err = res.Decode(&es)

	if assert.NoError(t, err) {
		assert.EqualValues(t, http.StatusOK, res.StatusCode)
		assert.Greater(t, len(es), 0)
	}

	err = shutdown(eh)
	assert.NoError(t, err)
}

func TestGetExpense(t *testing.T) {
	eh, err := setup()
	e := seedExpense(t)
	var latestExpense Expense
	res := request(http.MethodGet, uri("expenses", strconv.Itoa((e.ID))), nil)
	err = res.Decode(&latestExpense)

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, e.ID, latestExpense.ID)
		assert.NotEmpty(t, latestExpense.Title)
		assert.NotEmpty(t, latestExpense.Note)
		assert.NotEmpty(t, latestExpense.Tags)
	}
	err = shutdown(eh)
	assert.NoError(t, err)
}

func seedExpense(t *testing.T) Expense {
	var e Expense
	body := bytes.NewBufferString(`{
		"title": "strawberry smoothie",
		"amount": 79,
		"note": "night market promotion discount 10 bath", 
		"tags": ["food", "beverage"]
	}`)

	err := request(http.MethodPost, uri("expenses"), body).Decode(&e)
	if err != nil {
		t.Fatal("can't create expense:", err)
	}

	return e
}

func uri(paths ...string) string {
	host := "http://localhost:2565"
	if paths == nil {
		return host
	}

	url := append([]string{host}, paths...)
	return strings.Join(url, "/")
}

type Response struct {
	*http.Response
	err error
}

func (r *Response) Decode(v interface{}) error {
	if r.err != nil {
		return r.err
	}

	return json.NewDecoder(r.Body).Decode(v)
}

func request(method, url string, body io.Reader) *Response {
	req, _ := http.NewRequest(method, url, body)
	auth := "Bearer api-key-naja"
	req.Header.Add("Authorization", auth)
	req.Header.Add("Content-Type", "application/json")
	client := http.Client{}
	res, err := client.Do(req)
	return &Response{res, err}
}

func setup() (*echo.Echo, error) {
	port := ":2565"
	db, err := sql.Open("postgres", "postgresql://root:root@db/go-example-db?sslmode=disable")

	if err != nil {
		log.Fatal("Connect to database error", err)
	}

	eh := echo.New()
	go func(e *echo.Echo) {
		h := InitDB(db)
		e.Use(middleware.KeyAuth(func(key string, c echo.Context) (bool, error) {
			return key == "api-key-naja", nil
		}))

		e.GET("/expenses", h.GetExpenses)
		e.GET("/expenses/:id", h.GetExpense)
		e.PUT("/expenses/:id", h.UpdateExpense)
		e.POST("/expenses", h.CreateExpenses)

		e.Start(port)
	}(eh)
	for {
		conn, _ := net.DialTimeout("tcp", fmt.Sprintf("localhost%s", port), 30*time.Second)
		if conn != nil {
			conn.Close()
			break
		}
	}

	return eh, err
}

func shutdown(eh *echo.Echo) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return eh.Shutdown(ctx)
}
