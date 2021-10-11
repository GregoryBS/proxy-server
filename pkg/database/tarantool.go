package database

import (
	"bytes"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
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
	db.Lock()
	db.counter += 1
	db.Unlock()

	cookies := req.Cookies()
	cMap := make(map[string]string, 0)
	for _, c := range cookies {
		cMap[c.Name] = c.Value
	}
	readBody := make([]byte, req.ContentLength)
	var parsedBody url.Values
	if req.ContentLength > 0 {
		req.Body.Read(readBody)
		parsedBody["form"] = []string{string(readBody)}
	}
	if req.ParseForm() == nil {
		parsedBody = req.PostForm
	}

	_, err := db.tConn.Insert("requests", []interface{}{db.counter, req.Method, req.URL.Path, req.URL.Query(), req.Header, cMap, parsedBody})
	if err != nil {
		return err
	}

	readBody = make([]byte, resp.ContentLength)
	buf := bytes.NewBuffer(readBody)
	io.Copy(buf, resp.Body)
	_, err = db.tConn.Insert("responses", []interface{}{db.counter, resp.Status, resp.StatusCode, buf.String(), resp.Header})
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) Select(table string, id int) ([]interface{}, error) {
	result, err := db.tConn.Select(table, "primary", 0, 1, tarantool.IterEq, []interface{}{id})
	if err != nil {
		log.Printf("Cannot find in table: %s id: %d - error: %v", table, id, err)
		return nil, err
	}
	if len(result.Data) == 0 {
		log.Printf("Cannot find in table: %s id: %d - error: empty data", table, id)
		return nil, errors.New("No data with such id")
	}
	return result.Data[0].([]interface{}), nil
}
