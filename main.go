package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"html/template"
	"net/http"
)

type Article struct {
	Id                     uint16
	Title, Anons, FullText string
}

func index(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("template/index.html", "template/header.html")

	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/golang")
	if err != nil {
		panic(err)
	}

	defer db.Close()

	res, err := db.Query("SELECT * FROM `articles`")
	if err != nil {
		panic(err)
	}

	for res.Next() {
		var post Article
		err = res.Scan(&post.Title, &post.Anons, &post.FullText)
		if err != nil {
			panic(err)
		}

		//fmt.Println(fmt.Sprintf("User: %s %s with age %d", user.Name, user.LastName, user.Age))
	}

	t.ExecuteTemplate(w, "index", nil)
}

func create(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("template/create.html", "template/header.html")

	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	t.ExecuteTemplate(w, "create", nil)
}

func saveArticle(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	anons := r.FormValue("anons")
	fullText := r.FormValue("full_text")

	if title == "" || anons == "" || fullText == "" {
		fmt.Fprintf(w, "Не все данные заполнены")
	} else {
		db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/golang")
		if err != nil {
			panic(err)
		}

		defer db.Close()

		insert, err := db.Query(fmt.Sprintf("INSERT INTO `articles` (`title`, `anons`, `full_text`) VALUES('%s', '%s', '%s')",
			title, anons, fullText))
		if err != nil {
			panic(err)
		}

		defer insert.Close()

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func handleFunc() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	http.HandleFunc("/", index)
	http.HandleFunc("/create", create)
	http.HandleFunc("/save_article", saveArticle)
	http.ListenAndServe(":8080", nil)
}

func main() {
	handleFunc()
}
