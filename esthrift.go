package main
/*
 * Example golang code to access Elasticsearch using Thrift
 * gen-go directory generated with:
 * ``thrift -gen go elasticsearch.thrift``
 */
import (
	"fmt"
	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/stuart-warren/es-thrift/gen-go/elasticsearch" // generated code
	"net"
)

const BUFFER_SIZE = 1024

var request = elasticsearch.RestRequest{
	Method: elasticsearch.Method_GET,
	Uri:    "/_cluster/health",
	Parameters: map[string]string{"pretty":"true"},
}

func Connect(host string, thriftPort string) (*elasticsearch.RestClient, error) {
	binaryProtocol := thrift.NewTBinaryProtocolFactoryDefault()
	socket, err := thrift.NewTSocket(net.JoinHostPort(host, thriftPort))
	if err != nil {
		return nil, err
	}
	bufferedTransport := thrift.NewTBufferedTransport(socket, BUFFER_SIZE)
	client := elasticsearch.NewRestClientFactory(bufferedTransport, binaryProtocol)
	if err := bufferedTransport.Open(); err != nil {
		return nil, err
	}
	return client, nil
}

func main() {
	client, err := Connect("localhost", "9500")
	if err != nil {
		// err
	}
	result, err2 := client.Execute(&request)
	if err2 != nil {
		// err2
	}
	fmt.Println(string(result.GetBody()))
}
