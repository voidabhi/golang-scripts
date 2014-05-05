
package main

import(
	"fmt"
	"io/ioutil"
	"encoding/json"
)

const SAMPLE_JSON = "sample.json"

type Person struct {
	Name string
	Age int64
}

func main(){
	b,err:= ioutil.ReadFile(SAMPLE_JSON)
	if err!=nil{
		fmt.Printf("Error!")
	}
	var m Person
	json.Unmarshal(b,&m)
	fmt.Printf(m.Name)
	fmt.Printf(string(m.Age))
}