// Package driver holds the driver interface.
package driver

import (
	"errors"
	"fmt"
	"log"
	neturl "net/url"
	"strings"

	"github.com/promoboxx/migrate/driver/bash"
	"github.com/promoboxx/migrate/driver/cassandra"
	"github.com/promoboxx/migrate/driver/mysql" // alias to allow `url string` func signature in New
	"github.com/promoboxx/migrate/driver/postgres"
	"github.com/promoboxx/migrate/file"
)

type TxnType int

const (
	TxnNone TxnType = iota
	TxnPerFile
	TxnSingle
)

// Driver is the interface type that needs to implemented by all drivers.
type Driver interface {

	// Initialize is the first function to be called.
	// Check the url string and open and verify any connection
	// that has to be made.
	Initialize(url string) error

	// Close is the last function to be called.
	// Close any open connection here.
	Close() error

	// FilenameExtension returns the extension of the migration files.
	// The returned string must not begin with a dot.
	FilenameExtension() string

	// Migrate is the heart of the driver.
	// It will receive a file which the driver should apply
	// to its backend or whatever. The migration function should use
	// the pipe channel to return any errors or other useful information.
	Migrate(file file.File, pipe chan interface{})

	// Version returns the current migration version.
	Version() (uint64, error)
}

// New returns Driver and calls Initialize on it
func New(url string, txnType TxnType) (Driver, error) {
	u, err := neturl.Parse(url)
	if err != nil {
		return nil, err
	}

	switch u.Scheme {
	case "postgres":
		// For postgres we support multiple transaction strategies
		var d Driver
		switch txnType {
		case TxnNone:
			d = &postgres.NoTxnDriver{}
			log.Println("Migration scripts will be executed with no explicit transactions")
		case TxnPerFile:
			d = &postgres.PerFileTxnDriver{}
			log.Println("Each migration script will be executed in its own transaction")
		case TxnSingle:
			d = &postgres.SingleTxnDriver{}
			log.Println("All migration scripts will be executed in a single transaction")
		}

		verifyFilenameExtension("postgres", d)
		if err := d.Initialize(url); err != nil {
			return nil, err
		}
		return d, nil

	case "mysql":
		d := &mysql.Driver{}
		verifyFilenameExtension("mysql", d)
		if err := d.Initialize(url); err != nil {
			return nil, err
		}
		return d, nil

	case "bash":
		d := &bash.Driver{}
		verifyFilenameExtension("bash", d)
		if err := d.Initialize(url); err != nil {
			return nil, err
		}
		return d, nil

	case "cassandra":
		d := &cassandra.Driver{}
		verifyFilenameExtension("cassanda", d)
		if err := d.Initialize(url); err != nil {
			return nil, err
		}
		return d, nil
	default:
		return nil, errors.New(fmt.Sprintf("Driver '%s' not found.", u.Scheme))
	}
}

// verifyFilenameExtension panics if the drivers filename extension
// is not correct or empty.
func verifyFilenameExtension(driverName string, d Driver) {
	f := d.FilenameExtension()
	if f == "" {
		panic(fmt.Sprintf("%s.FilenameExtension() returns empty string.", driverName))
	}
	if f[0:1] == "." {
		panic(fmt.Sprintf("%s.FilenameExtension() returned string must not start with a dot.", driverName))
	}
}

// getTxnType returns the transaction behavior specified, or an error for unknown
// or undefined behaviors
func GetTxnType(txnType string) (TxnType, error) {
	switch strings.ToLower(txnType) {
	case "none":
		return TxnNone, nil

	case "single":
		return TxnSingle, nil

	case "perfile":
		return TxnPerFile, nil
	}

	return TxnNone, fmt.Errorf("Unknown transaction type requested: '%s'", txnType)
}
