package pkg

import (
	"database/sql"
	"fmt"
	"strconv"

	_ "github.com/lib/pq"
)

func SetupDB() *sql.DB {

	port, _ := strconv.Atoi(GetEnvVariable("DB_PORT"))

	dbInfo := fmt.Sprintf(
		"postgres://%s@%s:%d/%s?sslmode=disable",
		GetEnvVariable("DB_USER"), GetEnvVariable("DB_HOST"), port, GetEnvVariable("DB_NAME"))

	db, err := sql.Open("postgres", dbInfo)

	if err != nil {
		fmt.Printf(err.Error() + "\n")
	}
	// defer db.Close()

	if err != nil {
		errMsg := fmt.Sprintf("Error connecting to the database: %s", err)
		fmt.Printf(errMsg + "\n")
		//ErrorLogger.Println(errMsg)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	logMsg := "Successfully connected to the database"
	fmt.Printf(logMsg + "\n")
	//InfoLogger.Println(logMsg)

	return db
}
