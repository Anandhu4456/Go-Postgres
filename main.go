package main

import (
	"database/sql"
	"fmt"
	"net/http"

	_ "github.com/lib/pq"
)

var db *sql.DB
var err error

type Book struct{
	isbn string
	title string
	author string
	price float32
}

func init() {
	db, err = sql.Open("postgres", "postgres://anandhu:password@localhost/bookstore?sslmode=disable")
	if err != nil {
		panic(err)
	}
	if err = db.Ping(); err != nil {
		panic(err)
	}
	fmt.Println("DATABASE IS CONNECTED ...")
}
func main() {
	http.HandleFunc("/", webHandler)

	http.ListenAndServe(":8080", nil)
}

func webHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		http.Error(w, "405 error", http.StatusMethodNotAllowed)
		return
	}
	rows, err := db.Query("SELECT * FROM books")
	if err !=nil{
		http.Error(w,"query error",http.StatusMethodNotAllowed)
		return
	}
	defer rows.Close()

	var books = make([]Book,0)
	for rows.Next(){
		bk:=Book{}
	err =rows.Scan(&bk.isbn, &bk.title, &bk.author, &bk.price)
	if err!=nil{
		http.Error(w,"internalserver error",http.StatusInternalServerError)
		return
	}
	books = append(books,bk)
	}
	for _,bk:=range books{
		fmt.Fprintf(w,"%s, %s, %s, $%.2f\n", bk.isbn,bk.title, bk.author, bk.price)
	}

}
