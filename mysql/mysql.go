package mysql

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/jpillora/backoff"
)

type nullLogger struct{}

// Print Prints the log output.
func (logger nullLogger) Print(v ...interface{}) {}

// OpenConnection Returns a sql.DB instance connected to the MySQL shard.
func OpenConnection(host string, port int, username string, password string, timeout string) (*sql.DB, error) {
	mysql.SetLogger(nullLogger{})

	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%d)/?timeout=%s", username, password, host, port, timeout)
	connection, err := sql.Open("mysql", connectionString)

	if err != nil {
		return nil, err
	}
	return connection, nil
}

// GetServerID Returns the Id of the MySQL server instance.
// SELECT @@server_id
func GetServerID(host string, port int, username string, password string, timeout string) (string, error) {
	var value string

	rows, err := GetQueryResponse(host, port, username, password, timeout, "SELECT @@server_id;")
	if err != nil {
		return "", err
	}
	defer rows.Close()

	rows.Next()
	err = rows.Scan(&value)
	fmt.Println(value)
	if err != nil {
		return "", err
	}

	err = rows.Err()
	if err != nil {
		return "", err
	}
	return value, nil
}

// GetHostname Returns the hostname of the MySQL server instance.
// SELECT @@hostname
func GetHostname(host string, port int, username string, password string, timeout string) (string, error) {
	var value string

	rows, err := GetQueryResponse(host, port, username, password, timeout, "SELECT @@hostname;")
	if err != nil {
		return "", err
	}
	defer rows.Close()

	rows.Next()
	err = rows.Scan(&value)
	if err != nil {
		return "", err
	}

	err = rows.Err()
	if err != nil {
		return "", err
	}
	return value, nil
}

// GetQueryResponse Returns the rows from the specified SQL query.
func GetQueryResponse(host string, port int, username string, password string, timeout string, sql string) (*sql.Rows, error) {
	backoff := &backoff.Backoff{
		Min:    100 * time.Millisecond,
		Max:    10 * time.Second,
		Factor: 2,
		Jitter: true,
	}

	for {
		duration := backoff.Duration()

		db, err := OpenConnection(host, port, username, password, timeout)
		if err != nil {
			//fmt.Printf("%s, retrying in %s\n", err, duration)
			time.Sleep(duration)
			continue
		}
		defer db.Close()

		rows, err := db.Query(sql)
		if err != nil {
			//fmt.Printf("%s, retrying in %s\n", err, duration)
			time.Sleep(duration)
			continue
		}
		return rows, nil
	}
}
