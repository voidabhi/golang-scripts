package main

import "fmt"
import yaml2 "gopkg.in/yaml.v2"
import "github.com/ghodss/yaml"

func marshall() {
	x := struct {
		Name string
		Age  int
	}{
		"jeno",
		33,
	}
	fmt.Println("marshall some yaml ...")
	out, err := yaml2.Marshal(x)
	if err != nil {
		panic(err)
	}

	fmt.Printf("--> original:\n%#v\n--> marshalled:\n%s", x, out)
}

func str2yaml() {
	y := []byte(`name: jeno
age: 33
food:
  - pacal
  - rum
`)
	j2, err := yaml.YAMLToJSON(y)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	fmt.Println(string(j2))
}
func main() {
	//marshall()
	str2yaml()
}
