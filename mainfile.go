package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq"
)


const (
	DbHost = "db"
	DbUser = "postgres-dev"
	DbPassword = "password"
	DbName = "dev"
	Migration = `CREATE TABLE IF NOT EXISTS cmp (
id user PRIMARY KEY,
role text NOT NULL,
task text NOT NULL,
created_at timestamp with time zone DEFAULT current_timestamp)`

)

type cmp struct {
	User string `json:"user" binding:"required"`
	Role string `json:"role" binding:"required"`
	CreatedAt time.Time `json:"created_at"`
}

var db *sql.DB

func GetData() ([]cmp, error){
	const q = `SELECT user, role, created_at FROM cmp ORDER BY created_at DESC LIMIT 100`

	rows, err := db.Query(q)
	if err != nil {
		return nil, err
	}

	results := make([]cmp, 0)

	for rows.Next() {
		var user string
		var role string
		var createAt time.Time
		// scanning the data from the returned rows
		err = rows.Scan(&user, &role, &createAt)
		if err != nil {
			return nil, err
		}
		// creating a new result
		results = append(results, cmp{user, role, createAt})
	}

	return results, nil
}

func AddData(data cmp) error {
	const q = `INSERT INTO data(user, role, created_at) VALUES ($1, $2, $3)`
	_, err := db.Exec(q, data.User, data.Role, data.CreatedAt)
	return err
}

func main() {
	var err error

	// create a router with a default configuration
	r := gin.Default()
	// endpoint to retrieve all posted bulletins
	r.GET("/board", func(context *gin.Context) {
		results, err := GetData()
		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"status": "internal error: " + err.Error()})
			return
		}
		context.JSON(http.StatusOK, results)
	})
	r.POST("/board", func(context *gin.Context) {
		var b cmp
		// reading the request's body & parsing the json
		if context.Bind(&b) == nil {
			b.CreatedAt = time.Now()
			if err := AddData(b); err != nil {
				context.JSON(http.StatusInternalServerError, gin.H{"status": "internal error: " + err.Error()})
				return
			}
			context.JSON(http.StatusOK, gin.H{"status": "ok"})
			return
		}
		// if binding was not successful, return an error
		context.JSON(http.StatusUnprocessableEntity, gin.H{"status": "invalid body"})
	})
	// open a connection to the database
	dbInfo := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", DbHost, DbUser, DbPassword, DbName)
	db, err = sql.Open("postgres", dbInfo)
	if err != nil {
		panic(err)
	}
	// do not forget to close the connection
	defer db.Close()
	// ensuring the table is created
	_, err = db.Query(Migration)
	if err != nil {
		log.Println("failed to run migrations", err.Error())
		return
	}
	// running the http server
	log.Println("running..")
	if err := r.Run(":8080"); err != nil {
		panic(err)
	}
}


