package cddd

import (
	"errors"
)

type Cddd interface {
	DatabaseExists(database string) (bool, error)
	CreateDatabase(database string, idempotent bool) (err error)
	DropDatabase(database string, idempotent bool) (err error)
}

func New(driver, username, password, hostname, port string) (Cddd, error) {
	if driver == "" {
		return nil, errors.New("driver is null")
	}

	if username == "" {
		return nil, errors.New("username is null")
	}

	switch driver {
	case "mysql":
		return newMySQL(driver, username, password, hostname, port)
	case "postgres":
		return newPostgreSQL(driver, username, password, hostname, port)
	default:
		return nil, errors.New("this driver is not supported: " + driver)
	}

}
