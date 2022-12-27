package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/iamstarter/expenseapi/expense"
	_ "github.com/lib/pq"
)

func main() {

	expense.InitDB()

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

	go func() {
		if err := e.Start(port); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("Shutting down the server")
		}
	}()

	log.Println("Server started at: %", port)

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	<-shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}

	log.Println("Server stopped")
}
