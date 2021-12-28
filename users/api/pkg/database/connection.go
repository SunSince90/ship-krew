package database

import (
	"fmt"
	"net/url"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	queryCharset       string = "charset"
	queryParseTime     string = "parseTime"
	queryTrue          string = "True"
	queryLocation      string = "loc"
	queryLocationLocal string = "Local"
	queryReadTimeout   string = "readTimeout"
	queryWriteTimeout  string = "writeTimeout"
	queryTimeout       string = "timeout"
)

func NewDatabaseConnection(settings *Settings) (*gorm.DB, error) {
	if settings == nil {
		return nil, fmt.Errorf("no database settings provided")
	}

	dsn, err := buildDSN(settings)
	if err != nil {
		return nil, fmt.Errorf("cannot successfully create database address: %w", err)
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("cannot establish connection to database: %w", err)
	}

	return db, nil
}

func buildDSN(settings *Settings) (string, error) {
	host := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		settings.User,
		settings.Password,
		settings.Address,
		settings.Port,
		settings.Name)

	dsn, err := url.Parse(host)
	if err != nil {
		return "", fmt.Errorf("cannot create query: %w", err)
	}

	queryValues := url.Values{}
	queryValues.Add(queryCharset, settings.Charset)
	queryValues.Add(queryParseTime, queryTrue)
	queryValues.Add(queryLocation, queryLocationLocal)
	queryValues.Add(queryReadTimeout, settings.ReadTimeout.String())
	queryValues.Add(queryWriteTimeout, settings.WriteTimeout.String())
	queryValues.Add(queryTimeout, time.Minute.String())

	if _, err := url.ParseQuery(queryValues.Encode()); err != nil {
		return "", fmt.Errorf("cannot build database address: %w", err)
	}

	dsn.RawQuery = queryValues.Encode()
	return dsn.String(), nil
}
