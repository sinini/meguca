/*
 Initialises and loads RethinkDB
*/

package db

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/bakape/meguca/config"
	"github.com/bakape/meguca/util"
	r "github.com/dancannon/gorethink"
)

const dbVersion = 6

var (
	// RSession exports the RethinkDB connection session. Used globally by the
	// entire server.
	RSession *r.Session

	// AllTables are all tables needed for meguca operation
	AllTables = [...]string{"main", "threads", "images", "accounts"}
)

// Document is a eneric RethinkDB Document. For DRY-ness.
type Document struct {
	ID string `gorethink:"id"`
}

// Central global information document
type infoDocument struct {
	Document
	DBVersion int `gorethink:"dbVersion"`

	// Is incremented on each new post. Ensures post number uniqueness
	PostCtr int64 `gorethink:"postCtr"`
}

// LoadDB establishes connections to RethinkDB and Redis and bootstraps both
// databases, if not yet done.
func LoadDB() (err error) {
	conf := config.Get().Rethinkdb

	if err := Connect(conf.Addr); err != nil {
		return err
	}

	var isCreated bool
	err = One(r.DBList().Contains(conf.Db), &isCreated)
	if err != nil {
		return util.WrapError("Error checking, if database exists", err)
	}
	if isCreated {
		RSession.Use(conf.Db)
		return verifyDBVersion()
	}

	if err := InitDB(conf.Db); err != nil {
		return err
	}

	go runCleanupTasks()
	return nil
}

// Connect establishes a connection to RethinkDB. Address passed separately for
// easier testing.
func Connect(addr string) (err error) {
	if addr == "" { // For easier use in tests
		addr = "localhost:28015"
	}
	RSession, err = r.Connect(r.ConnectOpts{Address: addr})
	if err != nil {
		err = util.WrapError("Error connecting to RethinkDB", err)
	}
	return
}

// Confirm database verion is compatible, if not refuse to start, so we don't
// mess up the DB irreversably.
func verifyDBVersion() error {
	var version int
	err := One(GetMain("info").Field("dbVersion"), &version)
	if err != nil {
		return util.WrapError("Error reading database version", err)
	}
	if version != dbVersion {
		return fmt.Errorf(
			"Incompatible RethinkDB database version: %d. "+
				"See docs/migration.md",
			version,
		)
	}
	return nil
}

// InitDB initialize a rethinkDB database
func InitDB(dbName string) error {
	log.Printf("Initialising database '%s'", dbName)
	if err := Write(r.DBCreate(dbName)); err != nil {
		return util.WrapError("Error creating database", err)
	}

	RSession.Use(dbName)

	if err := CreateTables(); err != nil {
		return err
	}

	main := [...]interface{}{
		infoDocument{Document{"info"}, dbVersion, 0},

		// History aka progress counters of boards, that get incremented on
		// post creation
		Document{"histCounts"},
	}
	if err := Write(r.Table("main").Insert(main)); err != nil {
		return util.WrapError("Error initializing database", err)
	}

	return CreateIndeces()
}

// CreateTables creates all tables needed for meguca operation
func CreateTables() error {
	fns := make([]func() error, 0, len(AllTables))

	for _, table := range AllTables {
		if table == "images" {
			continue
		}
		fns = append(fns, createTable(table))
	}

	fns = append(fns, func() error {
		return Write(r.TableCreate("images", r.TableCreateOpts{
			PrimaryKey: "SHA1",
		}))
	})

	return util.Waterfall(fns)
}

func createTable(name string) func() error {
	return func() error {
		return Write(r.TableCreate(name))
	}
}

// CreateIndeces create secondary indeces for faster table queries
func CreateIndeces() error {
	fns := []func() error{
		func() error {
			return Write(r.Table("threads").IndexCreate("board"))
		},
	}

	// Make sure all indeces are ready to avoid the race condition of and index
	// being accessed before its full creation.
	for _, table := range AllTables {
		fns = append(fns, waitForIndex(table))
	}

	return util.Waterfall(fns)
}

func waitForIndex(table string) func() error {
	return func() error {
		return Exec(r.Table(table).IndexWait())
	}
}

// UniqueDBName returns a unique datatabase name. Needed so multiple concurent
// `go test` don't clash in the same database.
func UniqueDBName() string {
	return "meguca_tests_" + strconv.FormatInt(time.Now().UnixNano(), 10)
}
