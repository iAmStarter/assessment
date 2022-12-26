package expense

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"strconv"
	"strings"
	"testing"
)

func TestCreateExpense(t *testing.T) {

	body := bytes.NewBufferString(`{
		"title": "strawberry smoothie",
		"amount": 79,
		"note": "night market promotion discount 10 bath", 
		"tags": ["food", "beverage"]
	}`)

	var e Expense
	res := request(http.MethodPost, uri("expenses"), body)
	err := res.Decode(&e)

	assert.Nil(t, err)
	assert.EqualValues(t, http.StatusCreated, res.StatusCode)
	assert.NotEqual(t, 0, e.ID)
	assert.Equal(t, "strawberry smoothie", e.Title)
	assert.Equal(t, "night market promotion discount 10 bath", e.Note)
	assert.Equal(t, float32(79), e.Amount)
	assert.Equal(t, []string{"food", "beverage"}, e.Tags)
}

func TestUpdateExpense(t *testing.T) {
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
	err := res.Decode(&update)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, e.ID, update.ID)
	assert.Equal(t, "apple smoothie", update.Title)
	assert.Equal(t, "no discount", update.Note)
	assert.Equal(t, float32(89), update.Amount)
	assert.Equal(t, []string{"beverage"}, update.Tags)
}

func TestGetAllExpenses(t *testing.T) {

	seedExpense(t)
	var es []Expense
	res := request(http.MethodGet, uri("expenses"), nil)
	err := res.Decode(&es)

	assert.Nil(t, err)
	assert.EqualValues(t, http.StatusOK, res.StatusCode)
	assert.Greater(t, len(es), 0)
}

func TestGetExpense(t *testing.T) {

	e := seedExpense(t)
	var latestExpense Expense
	res := request(http.MethodGet, uri("expenses", strconv.Itoa((e.ID))), nil)
	err := res.Decode(&latestExpense)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, e.ID, latestExpense.ID)
	assert.NotEmpty(t, latestExpense.Title)
	assert.NotEmpty(t, latestExpense.Note)
	assert.NotEmpty(t, latestExpense.Tags)
}

func seedExpense(t *testing.T) Expense {
	var e Expense
	body := bytes.NewBufferString(`{
		"title": "Soju",
		"amount": 9000,
		"note": "Need to drink Soju 10 units on new year festival", 
		"tags": ["drinking"]
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
	auth := "Bearer gopher2022"
	req.Header.Add("Authorization", auth)
	req.Header.Add("Content-Type", "application/json")
	client := http.Client{}
	res, err := client.Do(req)
	return &Response{res, err}
}
