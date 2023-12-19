package main

import (
	"database/sql"
	"fmt"

	"github.com/baseba/got/handler"
	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
)

type Album struct {
	ID     int64
	Title  string
	Artist string
	Price  float32
}

var db *sql.DB

func main() {
	// Capture connection properties.
	// Get a database handle.
	// var err error
	// db, err = sql.Open("mysql", "(https://pokeslots-baseba.turso.io)/pokeslots?token=eyJhbGciOiJFZERTQSIsInR5cCI6IkpXVCJ9.eyJpYXQiOiIyMDIzLTEyLTE5VDA3OjMyOjUzLjU1NjE3NzA5OVoiLCJpZCI6ImZkNzlmM2YyLTllM2QtMTFlZS1iNTk2LTEyYWIwZGY3MGIxZiJ9.FZAaaz6hXlsamztLiQZ19XEbdJZOLE9xcf1HWIJIzpeulJxTZPDunXrNami39PKYx3jVDmP-DP0BZGMhsRK6Aw")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// pingErr := db.Ping()
	// if pingErr != nil {
	// 	log.Fatal(pingErr)
	// }
	// fmt.Println("Connected!")

	app := echo.New()

	userHandler := handler.UserHandler{}
	app.GET("/user", userHandler.HandleUserShow)
	slotsHandler := handler.SlotsHandler{}
	app.GET("/slots/:money", slotsHandler.HandleSlotsShow)

	app.POST("/win/:amount", func(c echo.Context) error {
		fmt.Println("you win " + c.Param("amount") + " current money = " + c.Path())
		return nil
	})

	app.POST("/lose/:amount", func(c echo.Context) error {
		fmt.Println("you lost " + c.Param("amount"))
		return nil
	})

	app.Start(":3000") //envs?
	fmt.Println("im alive!")

}
