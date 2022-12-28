package expense

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

func (h *handler) CreateExpenses(c echo.Context) error {
	var e Expense
	err := c.Bind(&e)
	fmt.Println("expenseObj", e.Tags)

	if err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	row := h.database.QueryRow("INSERT INTO expenses (title, amount, note, tags) values ($1, $2, $3, $4) RETURNING id, title, amount, note, tags", e.Title, e.Amount, e.Note, pq.Array(&e.Tags))

	err = row.Scan(&e.ID, &e.Title, &e.Amount, &e.Note, pq.Array(&e.Tags))

	if err != nil {
		fmt.Println("can't scan id", err)
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, e)
}
