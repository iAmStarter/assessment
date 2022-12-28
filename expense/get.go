package expense

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

func (h *handler) GetExpenses(c echo.Context) error {
	rows, err := h.database.Query("SELECT id, title, amount, note, tags FROM expenses")
	if err != nil {
		log.Fatal("can't prepare query all expenses statement", err)
	}

	if err != nil {
		log.Fatal("can't query all expenses ", err)
	}

	var expenses = []Expense{}

	for rows.Next() {
		e := Expense{}
		err := rows.Scan(&e.ID, &e.Title, &e.Amount, &e.Note, pq.Array(&e.Tags))

		if err != nil {
			log.Fatal("can't Scan row into variable ", err)
		}
		expenses = append(expenses, e)
	}

	return c.JSON(http.StatusOK, expenses)
}

func (h *handler) GetExpense(c echo.Context) error {
	id := c.Param("id")
	stmt, err := h.database.Prepare("SELECT id, title, amount, note, tags FROM expenses where id=$1")
	if err != nil {
		log.Fatal("can't prepare query one row statement ", err)
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}

	row := stmt.QueryRow(id)
	e := Expense{}
	err = row.Scan(&e.ID, &e.Title, &e.Amount, &e.Note, pq.Array(&e.Tags))
	if err != nil {
		log.Fatal("can't Scan row into variables ", err)
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, e)
}
