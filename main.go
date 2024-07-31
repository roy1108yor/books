package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

type Book struct {
	ID     int
	Title  string
	Author string
	Year   int
}

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("sqlite3", "./books.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	http.HandleFunc("/", listBooks)
	http.HandleFunc("/add", addBook)
	http.HandleFunc("/delete", deleteBook)
	http.HandleFunc("/edit", editBook)
	http.HandleFunc("/update", updateBook)

	log.Println("服务器在端口8080启动...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func listBooks(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, title, author, year FROM books")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var books []Book
	for rows.Next() {
		var book Book
		if err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Year); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		books = append(books, book)
	}

	tmpl, err := template.ParseFiles("books.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, books)
}

func addBook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "只允许POST请求", http.StatusMethodNotAllowed)
		return
	}

	title := r.FormValue("title")
	author := r.FormValue("author")
	year, err := strconv.Atoi(r.FormValue("year"))
	if err != nil {
		http.Error(w, "无效的出版年份", http.StatusBadRequest)
		return
	}

	_, err = db.Exec("INSERT INTO books (title, author, year) VALUES (?, ?, ?)", title, author, year)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func deleteBook(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	_, err := db.Exec("DELETE FROM books WHERE id = ?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func editBook(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	row := db.QueryRow("SELECT id, title, author, year FROM books WHERE id = ?", id)

	var book Book
	if err := row.Scan(&book.ID, &book.Title, &book.Author, &book.Year); err != nil {
		http.Error(w, "书籍未找到", http.StatusNotFound)
		return
	}

	tmpl, err := template.ParseFiles("edit.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, book)
}

func updateBook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "只允许POST请求", http.StatusMethodNotAllowed)
		return
	}

	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		http.Error(w, "无效的ID", http.StatusBadRequest)
		return
	}

	title := r.FormValue("title")
	author := r.FormValue("author")
	year, err := strconv.Atoi(r.FormValue("year"))
	if err != nil {
		http.Error(w, "无效的出版年份", http.StatusBadRequest)
		return
	}

	_, err = db.Exec("UPDATE books SET title = ?, author = ?, year = ? WHERE id = ?", title, author, year, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
