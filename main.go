package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"text/template"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func setup() error {

	var err error
	db, err = sql.Open("sqlite3", "./urlshortended.sqlite3")

	if err != nil {
		log.Fatal(err)
	}

	return nil

}

func homePage(w http.ResponseWriter, r *http.Request) {

	// fmt.Println(r.Method)

	if r.Method == "GET" {
		t, _ := template.ParseFiles("index.html")
		t.Execute(w, nil)
	} else {

		r.ParseForm()
		url := r.Form.Get("url")
		shortcut := r.Form.Get("shortcut") + strconv.Itoa(rand.Intn(100))
		fmt.Println(url)
		fmt.Println(shortcut)

		stmt, _ := db.Prepare("INSERT INTO urls(real_url, shortened_url) VALUES(?,?)")

		stmt.Exec(url, shortcut)

		stmt.Close()

		// db.Close()

		fmt.Fprintf(w, "<a href=\"url?shortcut=%v\">Copy This link </a> ", shortcut)

	}

}

func serveURL(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	if val, err := r.Form["shortcut"]; !err {
		fmt.Fprintln(w, "We got Hacky Boy Here")
	} else {
		shortcut := val[0]
		var url string
		stmt, _ := db.Prepare("SELECT real_url FROM urls WHERE shortened_url = ?")
		stmt.QueryRow(shortcut).Scan(&url)

		fmt.Println(url)

		http.Redirect(w, r, url, 301)

	}

}

func main() {
	err1 := setup()

	if err1 != nil {
		log.Fatal(err1)
	}

	http.HandleFunc("/", homePage)
	http.HandleFunc("/url", serveURL)
	err := http.ListenAndServe(":9090", nil)

	if err != nil {
		fmt.Println("Big Error!")
	}
}
