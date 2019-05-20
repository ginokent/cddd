package cddd

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
)

type postgres struct {
	driver           string
	dataSourcePrefix string
	dataSourceSuffix string
	db               *sql.DB
}

func newPostgreSQL(driver, username, password, hostname, port string) (Cddd, error) {
	if port == "" {
		port = "5432"
	}

	dataSourcePrefix := fmt.Sprintf("%s://%s:%s@%s:%s/",
		driver, username, password, hostname, port)

	dataSourceSuffix := "?sslmode=disable"

	conn, err := sql.Open(driver, dataSourcePrefix+dataSourceSuffix)
	if err != nil {
		return nil, fmt.Errorf("Open: %v", err)
	}

	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("Ping: %v", err)
	}

	return &postgres{
		driver:           driver,
		dataSourcePrefix: dataSourcePrefix,
		dataSourceSuffix: dataSourceSuffix,
		db:               conn,
	}, nil
}

func (cd *postgres) DatabaseExists(database string) (bl bool, err error) {
	check, _ := sql.Open(cd.driver, cd.dataSourcePrefix+database+cd.dataSourceSuffix)
	defer func() { // for catching deferred Close() error
		closeErr := check.Close()
		if closeErr == nil {
			return
		}
		err = fmt.Errorf("Close: failed to close: %v, original error: %v", closeErr, err)
	}()

	err = check.Ping()
	if err == nil {
		return true, nil
	}
	if err.Error() == fmt.Sprintf("pq: database \"%s\" does not exist", database) {
		return false, nil
	}
	return false, err
}

func (cd *postgres) CreateDatabase(database string, idempotent bool) (err error) {
	defer func() { // for catching deferred Close() error
		closeErr := cd.db.Close()
		if closeErr == nil {
			return
		}
		err = fmt.Errorf("Close: failed to close: %v, original error: %v", closeErr, err)
	}()

	if database == "" {
		return errors.New("database is null")
	}

	if idempotent {
		exists, err := cd.DatabaseExists(database)
		if err != nil {
			return err
		}
		if exists {
			return nil
		}
	}

	_, err = cd.db.Exec("CREATE DATABASE " + database)
	if err != nil {
		return fmt.Errorf("Exec: %v", err)
	}
	return err
}

func (cd *postgres) DropDatabase(database string, idempotent bool) (err error) {
	defer func() { // for catching deferred Close() error
		closeErr := cd.db.Close()
		if closeErr == nil {
			return
		}
		err = fmt.Errorf("Close: failed to close: %v, original error: %v", closeErr, err)
	}()

	if database == "" {
		return errors.New("database is null")
	}

	if idempotent {
		exists, err := cd.DatabaseExists(database)
		if err != nil {
			return err
		}
		if !exists {
			return nil
		}
	}

	_, err = cd.db.Exec("DROP DATABASE " + database)
	if err != nil {
		return fmt.Errorf("Exec: %v", err)
	}
	return err
}
