package main

import (
	"io"
	"log"
	"net/rpc"
)

type pipePair struct {
	reader *io.PipeReader
	writer *io.PipeWriter
}

func (this *pipePair) Read(p []byte) (int, error) {
	return this.reader.Read(p)
}

func (this *pipePair) Write(p []byte) (int, error) {
	return this.writer.Write(p)
}

func (this *pipePair) Close() error {
	this.writer.Close()
	return this.reader.Close()
}

type Serv struct {
	token string
}

func (this *Serv) SetToken(token string, dummy *int) error {
	this.token = token
	return nil
}

func (this *Serv) GetToken(dummy int, token *string) error {
	*token = this.token
	return nil
}

type Registrar struct {
}

func (this *Registrar) Register(name string, dummy *int) error {
	rpc.RegisterName(name, new(Serv))
	return nil	
}

func server(pipes pipePair) {
	rpc.Register(new(Registrar))
	rpc.ServeConn(&pipes)
}

func main() {
	var token string
	var in, out pipePair
	in.reader, out.writer = io.Pipe()
	out.reader, in.writer = io.Pipe()
	go server(out)
	client := rpc.NewClient(&in)
	// Register some objects
	client.Call("Registrar.Register", "First", nil)
	client.Call("Registrar.Register", "Second", nil)
	// Assign token values individually
	client.Call("First.SetToken", "abc", nil)
	client.Call("Second.SetToken", "def", nil)
	// Now try to read them
	client.Call("First.GetToken", 5, &token)
	log.Printf("first token is %v", token)
	client.Call("Second.GetToken", 5, &token)
	log.Printf("second token is %v", token)
}
