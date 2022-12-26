package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"os"

	"github.com/iamstarter/expenseapi/expense"
	_ "github.com/lib/pq"
)

func main() {

	expense.InitDB()
	fmt.Println("create table success")

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.KeyAuth(func(key string, c echo.Context) (bool, error) {
		return key == "27-Dec-2022", nil
	}))

	e.GET("/expenses", expense.GetExpensesHandler)
	e.GET("/expenses/:id", expense.GetExpenseHandler)
	e.PUT("/expenses/:id", expense.UpdateExpenseHandler)
	e.POST("/expenses", expense.CreateExpensesHandler)

	port := os.Getenv("PORT")
	log.Println("Server started at: %t", port)

	log.Fatal(e.Start(port))

	log.Println("bye bye")
}
