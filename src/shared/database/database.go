package database

import (
	"database/sql"
	"fmt"
	"../config"
	"../logger"
)

const SQL_CON_ERR = "Failed to create SQL connection ('%s'): %s"
const SQL_CON_SUCCESS = "SQL connection created ('%s')"

func getSqlConnectParams() string {
	host, port, dbase, user, pass :=
		config.Get(config.DB_HOST, "localhost"),
		config.Get(config.DB_PORT, "3306"),
		config.Get(config.DB_NAME, "foodcourt"),
		config.Get(config.DB_USER, "foodcourt"),
		config.Get(config.DB_PASS, "")

	return user + ":" + pass + "@" + host + ":" + port + "/" + dbase
}

// Creates a new SQL connection instance
func GetInstance() *sql.DB {

	log := logger.GetLogger()

	sqlDsn := getSqlConnectParams()

	db, err := sql.Open("mysql", sqlDsn)

	// Throw error if occurred
	if (err != nil) {
		errMsg := fmt.Sprintf(SQL_CON_ERR, sqlDsn, err.Error())
		log.Error(errMsg)
		panic(errMsg)
	}

	return db
}

// Test SQL connection
func TestConnection() (bool, error) {
	connection := GetInstance()
	defer connection.Close()

	err := connection.Ping()

	if err != nil {
		return false, err
	} else {
		return true, nil
	}
}
