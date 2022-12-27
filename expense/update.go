package expense

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

func UpdateExpenseHandler(c echo.Context) error {
	id := c.Param("id")
	var e Expense
	err := c.Bind(&e)
	_, err = db.Exec("UPDATE expenses SET title=$2, amount=$3, note=$4, tags=$5 where id=$1", id, e.Title, e.Amount, e.Note, pq.Array(e.Tags))
	if err != nil {
		log.Fatal("can't prepare query one row statement", err)
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}

	return GetExpenseHandler(c)
}
