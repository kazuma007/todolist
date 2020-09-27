package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	user := os.Getenv("USER")
	pass := os.Getenv("PASSWORD")
	path := os.Getenv("DB_PATH")
	db, err := sql.Open("mysql", user+":"+pass+"@"+"tcp(us-cdbr-east-02.cleardb.com)/"+path)
	// db, err := sql.Open("mysql", "b7aaf02befbd30:84fcc7f0@tcp(us-cdbr-east-02.cleardb.com)/heroku_8fe516dcd1295f5")
	// db, err := sql.Open("mysql", "db_user:password@tcp(localhost:3306)/my_db")

	if err != nil {
		log.Fatal(err)
	}

	todoitem := NewTodoItem(db)
	if err := todoitem.CreateTable(); err != nil {
		log.Fatal(err)
	}

	handlers := NewHandlers(todoitem)
	http.HandleFunc("/show", handlers.GetHandler)
	http.HandleFunc("/", handlers.TopPageHandler)
	http.HandleFunc("/save", handlers.SaveHandler)
	http.HandleFunc("/top", handlers.FormPageHandler)

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}
	fmt.Println("http:localhost/" + port + "で起動中...")
	http.ListenAndServe(":"+port, nil)
}
