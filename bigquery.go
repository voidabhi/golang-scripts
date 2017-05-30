
import (
	bigquery "code.google.com/p/google-api-go-client/bigquery/v2"
	"code.google.com/p/goauth2/oauth/jwt"
	"encoding/json"
	"fmt"
)

func main() {
	iss := "yourname@developer.gserviceaccount.com"
	scope := bigquery.BigqueryScope
	pem := `-----BEGIN RSA PRIVATE KEY-----
...
-----END RSA PRIVATE KEY-----`
	token := jwt.NewToken(iss, scope, []byte(pem))
	transport, _ := jwt.NewTransport(token)
	client := transport.Client()
	bq, _ := bigquery.New(client)
	call := bq.Tabledata.List("projectid", "dataset", "table")
	call.MaxResults(10)
	list, _ := call.Do()
	buf, _ := json.Marshal(list)
	fmt.Println(string(buf))
}
