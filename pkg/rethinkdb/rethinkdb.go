package pkg_rethinkdb

import (
	"os"

	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

// InitRethinkDB - initialize mongo
func InitRethinkDB() (*r.Session, error) {
	session, err := r.Connect(r.ConnectOpts{
		Address: os.Getenv("DB_URL"),
	})
	if err != nil {
		return nil, err
	}

	return session, nil
}
