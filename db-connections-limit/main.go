package main

import (
	"database/sql"
	"fmt"
	"os"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// Simulate simultaneous MySQL Connections & test its limitation.

// Launch several goroutines and INSERT into DB for each goroutine concurrently.
// Make these goroutine sleep for a second, to keep the connection open,
// since the close() will be deferred.
func main() {
	fmt.Println("Go - MySQL Connection Limitation test")

	// Open up our database connection.
	// NOTE: I Manually created a database called test.
	// db, err := sql.Open("mysql", "root:example@tcp(127.0.0.1:3306)/test")
	var sqlURL = RequireEnv("GOFUN_MYSQL_URL")
	fmt.Println("sqlURL : ", sqlURL)
	db, err := sql.Open("mysql", sqlURL)

	// db.SetMaxOpenConns(10)

	if err != nil {
		panic(err.Error())
	}

	// defer the close till after the main function has finished executing
	defer db.Close()

	// Set the maximum number of go routines that we wish to run concurrently.
	var MAXGOROUTINES = 20

	// To wait for multiple goroutines to finish, we can use a wait group.
	var wg sync.WaitGroup

	for i := 0; i < MAXGOROUTINES; i++ {

		wg.Add(1)

		// Launch several goroutines and INSERT into DB for each goroutine concurrently.
		go insertValueInDB(i, db, &wg)
	}

	fmt.Println("WaitGroup is waiting for the goroutines to finish")
	wg.Wait()
	fmt.Println("-- END - WaitGroup has finished blocking")
}

// insertValueInDB performs a db.Query insert.
// It opens up a new DB connection to do so. http://go-database-sql.org/connection-pool.html
func insertValueInDB(i int, db *sql.DB, wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Println("--- iteration count : ", i)

	insert, err := db.Query("INSERT INTO test_table VALUES ( ?, 'TEST' )", i)

	if err != nil {
		panic(err.Error())
	}

	defer insert.Close()
	fmt.Println("insert : ", insert)

	// Make this goroutine sleep for a second, to keep the connection open,
	// since the close() will be deferred.
	time.Sleep(2 * time.Second)
}

// RequireEnv looks up an environment variable and panics if
// it's not present or is an empty string.
func RequireEnv(env string) string {
	v, ok := os.LookupEnv(env)
	if !ok {
		panic("missing required environment variable: " + env)
	}
	if v == "" {
		panic("empty required environment variable: " + env)
	}
	return v
}
