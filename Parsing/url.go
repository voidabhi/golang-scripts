
package main

import (
	"fmt"
	"net/url"
	"strings"
)

type UrlData struct {
	Scheme string
	User string
	Password string
	Host string
	Port string
	Path string
	Query string
	Fragment string
}

func main(){

	s := "postgres://user:pass@example.com:5443/path?k=v#f"
	
	var data UrlData
	
	// Parses url
	u,err:=url.Parse(s)
	if err!=nil{
		panic(err)
	}
	
	data.Scheme = u.Scheme
	data.User = u.User.Username()
	data.Password,_ = u.User.Password()
	h := strings.Split(u.Host,":")
	data.Host = h[0]
	data.Port = h[1]
	
	data.Path = u.Path
	data.Fragment = u.Fragment
	
	data.Query = u.RawQuery
	// Parses raw query , m stores queries in key value pairs
	m,_:= url.ParseQuery(u.RawQuery)
}