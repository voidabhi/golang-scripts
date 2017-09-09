
package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// An uninteresting service.
type Service struct {
	ch        chan bool
	waitGroup *sync.WaitGroup
}

// Make a new Service.
func NewService() *Service {
	s := &Service{
		ch:        make(chan bool),
		waitGroup: &sync.WaitGroup{},
	}
	s.waitGroup.Add(1)
	return s
}

// Accept connections and spawn a goroutine to serve each one.  Stop listening
// if anything is received on the service's channel.
func (s *Service) Serve(listener *net.TCPListener) {
	defer s.waitGroup.Done()
	for {
		select {
		case <-s.ch:
			log.Println("stopping listening on", listener.Addr())
			listener.Close()
			return
		default:
		}
		listener.SetDeadline(time.Now().Add(1e9))
		conn, err := listener.AcceptTCP()
		if nil != err {
			if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
				continue
			}
			log.Println(err)
		}
		log.Println(conn.RemoteAddr(), "connected")
		s.waitGroup.Add(1)
		go s.serve(conn)
	}
}

// Stop the service by closing the service's channel.  Block until the service
// is really stopped.
func (s *Service) Stop() {
	close(s.ch)
	s.waitGroup.Wait()
}

// Serve a connection by reading and writing what was read.  That's right, this
// is an echo service.  Stop reading and writing if anything is received on the
// service's channel but only after writing what was read.
func (s *Service) serve(conn *net.TCPConn) {
	defer conn.Close()
	defer s.waitGroup.Done()
	for {
		select {
		case <-s.ch:
			log.Println("disconnecting", conn.RemoteAddr())
			return
		default:
		}
		conn.SetDeadline(time.Now().Add(1e9))
		buf := make([]byte, 4096)
		if _, err := conn.Read(buf); nil != err {
			if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
				continue
			}
			log.Println(err)
			return
		}
		if _, err := conn.Write(buf); nil != err {
			log.Println(err)
			return
		}
	}
}

func main() {

	// Listen on 127.0.0.1:48879.  That's my favorite port number because in
	// hex 48879 is 0xBEEF.
	laddr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:48879")
	if nil != err {
		log.Fatalln(err)
	}
	listener, err := net.ListenTCP("tcp", laddr)
	if nil != err {
		log.Fatalln(err)
	}
	log.Println("listening on", listener.Addr())

	// Make a new service and send it into the background.
	service := NewService()
	go service.Serve(listener)

	// Handle SIGINT and SIGTERM.
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)

	// Stop the service gracefully.
	service.Stop()

}
