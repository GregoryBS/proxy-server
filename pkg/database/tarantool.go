package database

import (
	"log"
	"net/http"
	"sync"

	"github.com/tarantool/go-tarantool"
)

type DB struct {
	tConn *tarantool.Connection
	sync.Mutex
	counter uint
}

func GetTarantool() *DB {
	addr := "tarantool:5555"
	opts := tarantool.Opts{User: "proxy", Pass: "proxy_pass"}
	conn, err := tarantool.Connect(addr, opts)
	if err != nil {
		log.Fatalln("Cannot connect to tarantool")
		return nil
	}
	return &DB{tConn: conn}
}

func CloseTarantool(db *DB) {
	db.tConn.Close()
}

func (db *DB) Insert(resp *http.Response, req *http.Request) error {
	_, err := db.tConn.Insert("requests", []interface{}{})
	if err != nil {
		return err
	}

	_, err = db.tConn.Insert("responses", []interface{}{})
	if err != nil {
		return err
	}
	return nil
}
