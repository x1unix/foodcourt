package database

import (
	"fmt"
	"../config"
	"../logger"
	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

const SQL_CON_ERR = "Failed to create SQL connection ('%s'): %s"
const SQL_CON_SUCCESS = "SQL connection created ('%s')"

// Returns SQL connect params
func getSqlConnectParams() string {
	host, port, dbase, user, pass :=
		config.Get(config.DB_HOST, "localhost"),
		config.Get(config.DB_PORT, "3306"),
		config.Get(config.DB_NAME, "foodcourt"),
		config.Get(config.DB_USER, "foodcourt"),
		config.Get(config.DB_PASS, "")

	// Create DSN builder
	dsnConfig := &mysql.Config{
		User: user,
		Passwd: pass,
		Net: "tcp",
		Addr: host + ":" + port,
		DBName: dbase,
	}

	// Format DSN to string
	dsn := dsnConfig.FormatDSN()

	return dsn
}

// Creates a new SQL connection instance
func GetInstance() *sqlx.DB {

	log := logger.GetLogger()

	sqlDsn := getSqlConnectParams()

	db, err := sqlx.Connect("mysql", sqlDsn)

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
	sqlDsn := getSqlConnectParams()

	connection, err := sqlx.Connect("mysql", sqlDsn)

	defer connection.Close()

	if err != nil {
		return false, err
	} else {
		return true, nil
	}
}
