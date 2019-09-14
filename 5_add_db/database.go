package main

import (
	"database/sql"
	"time"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
)

// ErrNoRowsAffected is returned if no rows were affected by the sql query
var ErrNoRowsAffected = errors.New("no rows affected")

// Database : Connection information struct.
type Database struct {
	connection        *sql.DB
	Driver            string `json:"driver"`
	Protocol          string `json:"protocol"`
	Host              string `json:"host"`
	Port              string `json:"port"`
	Name              string `json:"name"`
	User              string `json:"user"`
	Password          string `json:"password"`
	PingLast          time.Time
	PingCacheDuration time.Duration
}

// Open : Open a database connection.
func (d *Database) Open() error {
	dataSrc := d.User + ":" + d.Password + "@" + d.Protocol + "(" + d.Host + ":" + d.Port + ")/" + d.Name
	var err error
	if d.connection, err = sql.Open(d.Driver, dataSrc); nil != err {
		fmt.Println(err)
		return err
	}
	d.PingLast = time.Unix(0, 0)
	return nil
}

// ReOpen : ReOpen a database connection if it has been closed.
func (d *Database) ReOpen() error {
	if err := d.Ping(); nil != err {
		if err := d.Open(); nil != err {
			return err
		}
	}
	return nil
}

// Close : Close a database connection.
func (d *Database) Close() error {
	return d.connection.Close()
}

// Exec : Executes the command on an already open database connection.
func (d *Database) Exec(query string, args ...interface{}) (sql.Result, error) {
	var result sql.Result
	var err error
	result, err = d.connection.Exec(query, args...)
	if nil != err {
		return nil, err
	}
	return result, nil
}

// SetMaxConnections : Set the maximum number of connections that can be opened
func (d *Database) SetMaxConnections(n uint64) {
	d.connection.SetMaxOpenConns(int(n))
}

// SetMaxIdleConnections : Set the maximum number of idle connections to leave open
func (d *Database) SetMaxIdleConnections(n uint64) {
	d.connection.SetMaxIdleConns(int(n))
}

// SetConnMaxLifetime : Set the duration in seconds that connections should remain open for.
func (d *Database) SetConnMaxLifetime(seconds uint64) {
	d.PingCacheDuration = time.Second * time.Duration(seconds)
	d.connection.SetConnMaxLifetime(d.PingCacheDuration)
}

// QueryRow : Query the database and return the rows ready to be scanned.
func (d *Database) QueryRow(query string, args ...interface{}) *sql.Row {
	return d.connection.QueryRow(query, args...)
}

// Query : Query the database and return the rows ready to be scanned.
func (d *Database) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return d.connection.Query(query, args...)
}

// Ping : Test the database connectivity.
func (d *Database) Ping() error {
	if time.Now().Sub(d.PingLast) >= d.PingCacheDuration {
		d.PingLast = time.Now()
		return d.connection.Ping()
	}
	return nil
}

// Version : Return the SQL database engine version.
func (d *Database) Version() (string, error) {
	dbRow := d.connection.QueryRow("SELECT VERSION() AS v")
	var sqlVersion string
	if err := dbRow.Scan(&sqlVersion); nil != err {
		return "", err
	}
	return sqlVersion, nil
}
