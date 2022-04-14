package main

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
)

var schema = `
CREATE TABLE person (
	id SERIAL PRIMARY KEY,
    first_name text,
    last_name text,
    email text
);`

type Person struct {
	Id        int    `db:"id" json:"id"`
	FirstName string `db:"first_name" json:"first_name"`
	LastName  string `db:"last_name" json:"last_name"`
	Email     string `db:"email" json:"email"`
}

func main() {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Jakarta",
		"localhost", 5432, "fadli", "nill", "local")
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to database")

	if os.Getenv("DEV") == "true" {
		db.MustExec(schema)

		tx := db.MustBegin()
		tx.MustExec("INSERT INTO person (first_name, last_name, email) VALUES ($1, $2, $3)", "Jason", "Moiron", "jmoiron@jmoiron.net")
		tx.MustExec("INSERT INTO person (first_name, last_name, email) VALUES ($1, $2, $3)", "John", "Doe", "johndoeDNE@gmail.net")
		err = tx.Commit()
		if err != nil {
			log.Fatal(err)
		}
	}

	router := httprouter.New()
	router.GET("/person/:id", func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		var (
			id     = params.ByName("id")
			person Person
		)

		/*example:
		drop table person: SELECT name, email FROM users WHERE ID = '10';DROP TABLE person--*/

		//err := db.Get(&person, "SELECT * FROM person WHERE id = $1", id)
		err = db.Get(&person, fmt.Sprintf("SELECT * FROM person WHERE id = %s", id))
		if err != nil {
			ResponseError(writer, http.StatusInternalServerError, err)
			return
		}

		ResponseJSON(writer, http.StatusOK, person)
	})
	log.Print("server started at 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
