package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

type User struct {
	address  string
	username string
}

type Message struct {
	from    *User
	payload []byte
}

type Server struct {
	msgch    chan Message
	quitch   chan struct{}
	listener net.Listener
	users    map[string]*User
	address  string
}

func NewServer(address string) *Server {
	return &Server{
		address: address,
		msgch:   make(chan Message),
		quitch:  make(chan struct{}),
		users:   make(map[string]*User),
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.address)
	if err != nil {
		return err
	}

	defer ln.Close()
	s.listener = ln

	go s.acceptLoop()

	<-s.quitch
	close(s.msgch)
	return nil
}

func (s *Server) acceptLoop() {
	for {
		con, err := s.listener.Accept()
		if err != nil {
			break
		}

		usrAddr := con.RemoteAddr().String()

		_, found := s.users[usrAddr]

		if !found {
			con.Write([]byte("Enter username: "))
			s.users[usrAddr] = NewUser(usrAddr, "")
		}
		go s.handleConection(con)

	}
}

func NewUser(address, username string) *User {
	return &User{
		address:  address,
		username: username,
	}
}

func (s *Server) handleConection(con net.Conn) {
	buf := make([]byte, 2048)
	defer con.Close()
	for {
		n, err := con.Read(buf)
		if err != nil {
			fmt.Println("Client disconnected.")
			break
		}

		msg := buf[:n]
		usrAddr := con.RemoteAddr().String()

		usr, found := s.users[usrAddr]

		if found && usr.username == "" {
			s.users[usrAddr].username = strings.ReplaceAll(string(msg), "\n", "")
			con.Write([]byte("Welcome " + string(msg) + "\n"))
		} else {
			s.msgch <- Message{
				from:    s.users[usrAddr],
				payload: msg,
			}
		}

	}
}

func (s *Server) handleMessages() {
	logger := NewLogger("./chat.log")
	for msg := range s.msgch {
		formatted := fmt.Sprintf("> %s: %s", msg.from.username, string(msg.payload))
		fmt.Println(formatted)
		logger.Println(formatted)
	}
}

func NewLogger(filename string) *log.Logger {
	logFile, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}

	return log.New(logFile, "[silme]", log.Ldate|log.Ltime)
}

func main() {
	server := NewServer(":3000")

	go server.handleMessages()

	log.Fatal(server.Start())
}
