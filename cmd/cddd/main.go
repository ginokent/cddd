package main

import (
	"github.com/djeeno/cddd/pkg/cddd"
	"gopkg.in/urfave/cli.v1"
	"log"
	"os"
)

func main() {
	var (
		driver     string
		username   string
		password   string
		hostname   string
		port       string
		database   string
		idempotent bool
	)

	app := cli.NewApp()
	app.Name = "cddd"
	app.Version = "0.0.1"
	app.Usage = "Just execute query 'CREATE DATABASE' or 'DROP DATABASE'"
	app.EnableBashCompletion = true

	//
	// flags
	//
	driverFlag := cli.StringFlag{
		Name:        "driver, D",
		Usage:       "Set database driver (--driver (mysql|postgre))",
		EnvVar:      "CORD_DRIVER",
		Value:       "",
		Destination: &driver,
	}

	usernameFlag := cli.StringFlag{
		Name:        "username, u",
		Usage:       "Set username",
		EnvVar:      "CORD_USERNAME",
		Value:       "",
		Destination: &username,
	}

	passwordFlag := cli.StringFlag{
		Name:        "password, p",
		Usage:       "Set database user password",
		EnvVar:      "CORD_PASSWORD",
		Value:       "",
		Destination: &password,
	}

	hostnameFlag := cli.StringFlag{
		Name:        "hostname, H",
		Usage:       "Set database hostname or ipaddress",
		EnvVar:      "CORD_HOSTNAME",
		Value:       "localhost",
		Destination: &hostname,
	}

	portFlag := cli.StringFlag{
		Name:        "port, P",
		Usage:       "Set database host's port number",
		EnvVar:      "CORD_PORT",
		Value:       "",
		Destination: &port,
	}

	databaseFlag := cli.StringFlag{
		Name:        "database, d",
		Usage:       "Set database name",
		EnvVar:      "CORD_DATABASE",
		Value:       "",
		Destination: &database,
	}

	idempotentFlag := cli.BoolFlag{
		Name:        "idempotent, i",
		Usage:       "Execute query 'CREATE DATABASE' or 'DROP DATABASE' as idempotent (true) or not (false)",
		EnvVar:      "CORD_IDEMPOTENT",
		Destination: &idempotent,
	}

	app.Flags = []cli.Flag{
		driverFlag,
		usernameFlag,
		passwordFlag,
		hostnameFlag,
		portFlag,
		databaseFlag,
		idempotentFlag,
	}

	//
	// create
	//
	createAction := func(c *cli.Context) error {
		cd, err := cddd.New(driver, username, password, hostname, port)
		if err != nil {
			return err
		}
		if err := cd.CreateDatabase(database, idempotent); err != nil {
			return err
		}
		return nil
	}

	create := cli.Command{
		Name:   "create",
		Usage:  "Execute query 'CREATE DATABASE'",
		Action: createAction,
	}

	//
	// drop
	//
	dropAction := func(c *cli.Context) error {
		cd, err := cddd.New(driver, username, password, hostname, port)
		if err != nil {
			return err
		}
		if err := cd.DropDatabase(database, idempotent); err != nil {
			return err
		}
		return nil
	}

	drop := cli.Command{
		Name:   "drop",
		Usage:  "Execute query 'DROP DATABASE'",
		Action: dropAction,
	}

	//
	// sub commands
	//
	app.Commands = []cli.Command{
		create,
		drop,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}

}
