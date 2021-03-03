package goose

import (
	"database/sql"
	"fmt"
	"strconv"
	"sync"
)

const VERSION = "v2.7.0-rc3"

var (
	duplicateCheckOnce sync.Once
	minVersion         = int64(0)
	maxVersion         = int64((1 << 63) - 1)
	timestampFormat    = "20060102150405"
	verbose            = false
)

// SetVerbose set the goose verbosity mode
func SetVerbose(v bool) {
	verbose = v
}

// Run runs a goose command.
func Run(command string, db *sql.DB, dir string, args ...string) error {
	switch command {
	case "up":
		if err := UpUnapplied(db, dir); err != nil {
			return err
		}

		// NOTE: Business logic change from Forked rep.
		// We want to apply any previously unmigrated version.
		// The initial package will only apply new versions with dates AFTER
		// the last migration. In our fast pace / multi dev environment this leaves
		// migrations un applied once a hotfix surpasses it.
		// case "up":
		// 	if err := Up(db, dir); err != nil {
		// 		return err
		// 	}
	case "up-by-one":
		if err := UpByOne(db, dir); err != nil {
			return err
		}
	case "up-to":
		if len(args) == 0 {
			return fmt.Errorf("up-to must be of form: goose [OPTIONS] DRIVER DBSTRING up-to VERSION")
		}

		version, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("version must be a number (got '%s')", args[0])
		}
		if err := UpTo(db, dir, version); err != nil {
			return err
		}
	case "apply":
		if len(args) == 0 {
			return fmt.Errorf("apply must be of form: goose [OPTIONS] DRIVER DBSTRING apply VERSION")
		}
		version, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("version must be a number (got '%s')", args[0])
		}
		if err := Apply(db, dir, version); err != nil {
			return err
		}
	case "revert":
		if len(args) == 0 {
			return fmt.Errorf("revert must be of form: goose [OPTIONS] DRIVER DBSTRING revert VERSION")
		}
		version, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("version must be a number (got '%s')", args[0])
		}
		if err := Revert(db, dir, version); err != nil {
			return err
		}
	case "create":
		if len(args) == 0 {
			return fmt.Errorf("create must be of form: goose [OPTIONS] DRIVER DBSTRING create NAME [go|sql]")
		}

		migrationType := "go"
		if len(args) == 2 {
			migrationType = args[1]
		}
		if err := Create(db, dir, args[0], migrationType); err != nil {
			return err
		}
	case "down":
		if err := Down(db, dir); err != nil {
			return err
		}
	case "down-to":
		if len(args) == 0 {
			return fmt.Errorf("down-to must be of form: goose [OPTIONS] DRIVER DBSTRING down-to VERSION")
		}

		version, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("version must be a number (got '%s')", args[0])
		}
		if err := DownTo(db, dir, version); err != nil {
			return err
		}
	case "fix":
		if err := Fix(dir); err != nil {
			return err
		}
	case "redo":
		if err := Redo(db, dir); err != nil {
			return err
		}
	case "reset":
		if err := Reset(db, dir); err != nil {
			return err
		}
	case "status":
		if err := Status(db, dir); err != nil {
			return err
		}
	case "version":
		if err := Version(db, dir); err != nil {
			return err
		}
	default:
		return fmt.Errorf("%q: no such command", command)
	}
	return nil
}
