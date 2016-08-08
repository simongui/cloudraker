package mysql

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/jpillora/backoff"
)

// GetMasterStatusResult Represents the response from a 'show master status' MySQL query.
type GetMasterStatusResult struct {
	BinlogFile      string
	BinlogPos       int
	BinlogDoDB      string
	BinlogIgnoreDB  string
	ExecutedGtidSet string
}

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

// GetQueryResponse Returns the rows from the specified SQL query.
func GetQueryResponse(host string, port int, username string, password string, timeout string, sql string) (*sql.Rows, error) {
	backoff := &backoff.Backoff{
		Min:    100 * time.Millisecond,
		Max:    2 * time.Second,
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

// SetReplicationGrants Sets the proper grants for replicating from the specified MySQL instance.
func SetReplicationGrants(host string, port int, username string, password string, timeout string, replicationUser string, replicationPassword string) error {
	sql := fmt.Sprintf("GRANT REPLICATION SLAVE, REPLICATION CLIENT ON *.* TO '%s'@'%%' IDENTIFIED BY '%s';",
		replicationUser,
		replicationPassword)

	_, err := GetQueryResponse(host, port, username, password, timeout, sql)
	if err != nil {
		return err
	}
	return nil
}

// GetReadOnly Returns the read_only status of the specified MySQL instance.
// SELECT @@read_only
func GetReadOnly(host string, port int, username string, password string, timeout string) (bool, error) {
	var value int

	rows, err := GetQueryResponse(host, port, username, password, timeout, "SELECT @@read_only;")
	if err != nil {
		return false, err
	}
	defer rows.Close()

	rows.Next()
	err = rows.Scan(&value)
	if err != nil {
		return false, err
	}

	err = rows.Err()
	if err != nil {
		return false, err
	}
	return value != 0, nil
}

// SetReadOnly Sets the specified MySQL instance as read_only.
// show master status
func SetReadOnly(host string, port int, username string, password string, timeout string, readOnly bool) error {
	readOnlyValue := 0
	if readOnly == true {
		readOnlyValue = 1
	}
	sql := fmt.Sprintf("SET GLOBAL read_only=%d;",
		readOnlyValue)

	_, err := GetQueryResponse(host, port, username, password, timeout, sql)
	if err != nil {
		return err
	}
	return nil
}

// GetMasterStatus Returns the binlog and binlog position from the specified host.
// show master status
func GetMasterStatus(host string, port int, username string, password string, timeout string) (*GetMasterStatusResult, error) {
	var result = &GetMasterStatusResult{}

	rows, err := GetQueryResponse(host, port, username, password, timeout, "show master status;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rows.Next()

	err = rows.Scan(&result.BinlogFile, &result.BinlogPos, &result.BinlogDoDB, &result.BinlogIgnoreDB, &result.ExecutedGtidSet)

	if err != nil {
		return nil, err
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return result, nil

}
