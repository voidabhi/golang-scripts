
package main

import (
	"encoding/json"
	"log"
	"time"

	"code.google.com/p/goauth2/oauth/jwt"
	bigquery "github.com/google/google-api-go-client/bigquery/v2"
)

const (
	App     = "XXX"
	Dataset = "XXX"
	Table   = "XXX"
	Email   = "XXXXXXX@developer.gserviceaccount.com"
	Scope   = bigquery.BigqueryScope

	PEM = `
-----BEGIN RSA PRIVATE KEY-----
...
-----END RSA PRIVATE KEY-----`
)

func main() {

	token := jwt.NewToken(Email, Scope, []byte(PEM))
	transport, err := jwt.NewTransport(token)
	if err != nil {
		log.Fatal(err)
	}

	client := transport.Client()
	bq, err := bigquery.New(client)
	if err != nil {
		log.Fatal(err)
	}

	rows := make([]*bigquery.TableDataInsertAllRequestRows, 0)

	row := &bigquery.TableDataInsertAllRequestRows{
		Json: make(map[string]bigquery.JsonValue, 0),
	}
	row.Json["url"] = "https://github.com"
	row.Json["source"] = "example"
	row.Json["t"] = time.Now().Unix()
	row.Json["http_status"] = 200

	rows = append(rows, row)

	req := &bigquery.TableDataInsertAllRequest{
		Rows: rows,
	}

	call := bq.Tabledata.InsertAll(App, Dataset, Table, req)
	resp, err := call.Do()
	if err != nil {
		log.Fatal(err)
	}

	buf, _ := json.Marshal(resp)
	log.Print(string(buf))
}
