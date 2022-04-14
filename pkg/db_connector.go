package pkg

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 26257
	dbname   = "iotwin"
	user     = "root"
	password = ""
)

//CockroachDB
func SetupDB() *sql.DB {

	fmt.Println(os.Getenv("MQTT_CLIENT"))

	dbInfo := fmt.Sprintf(
		"postgres://%s@%s:%d/%s?sslmode=disable",
		user, host, port, dbname)

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
