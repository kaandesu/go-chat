package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
)

type Message struct {
	from    *User
	payload []byte
}

type Server struct {
	msgch        chan Message
	quitch       chan struct{}
	listener     net.Listener
	logger       *log.Logger
	users        map[string]*User
	childServers map[string]*Server
	serverch     chan *Server
	address      string
}

func NewServer(address string) *Server {
	return &Server{
		address:      address,
		msgch:        make(chan Message),
		quitch:       make(chan struct{}),
		users:        make(map[string]*User),
		childServers: make(map[string]*Server),
		serverch:     make(chan *Server, 5), // NOTE: 5 concurrent servers, i think, update later
		logger:       NewLogger("./chat.log"),
	}
}

func returnFirstAvailablePort(portNo int, tryUntil int) (string, error) {
	var err error
	var ln net.Listener
	for i := range tryUntil {
		host := ":" + strconv.Itoa(portNo+i)
		ln, err = net.Listen("tcp", host)
		if err == nil {
			ln.Close()
			return host, nil
		}
	}
	ln.Close()
	return "", errors.New("no port available try a larger range")
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.address)
	fmt.Printf("Starting server on [%s] \n", s.address)
	s.logger.Printf("Starting server on [%s] \n", s.address)
	if err != nil {
		return err
	}

	defer ln.Close()
	s.listener = ln

	go s.acceptLoop()
	go s.handleMessages()

	go s.handleChildServers()
	<-s.quitch
	close(s.msgch)
	return nil
}

func (s *Server) handleChildServers() {
	for server := range s.serverch {
		fmt.Printf("Server address: %s", server.address)
		if err := server.Start(); err != nil {
			fmt.Printf("Error Child Server: %s", err)
		}
	}
}

func (s *Server) acceptLoop() {
	for {
		con, err := s.listener.Accept()
		if err != nil {
			break
		}

		usrAddr := con.RemoteAddr().String()

		if _, found := s.users[usrAddr]; !found {
			con.Write([]byte("Enter username: "))
			s.users[usrAddr] = NewUser(usrAddr, "", con)
		}

		go s.handleConection(con)

	}
}

func (s *Server) handleConection(con net.Conn) {
	buf := make([]byte, 2048)

	defer con.Close()
	for {
		n, err := con.Read(buf)
		if err != nil {
			formatted := fmt.Sprintf("%s disconnected! (%d online) \n", s.users[con.RemoteAddr().String()].username, len(s.users)-1)
			s.logger.Print(formatted)
			fmt.Print(formatted)
			delete(s.users, con.RemoteAddr().String())
			break
		}

		msg := buf[:n]
		usrAddr := con.RemoteAddr().String()

		usr, found := s.users[usrAddr]

		isMessage := true
		if found && usr.username == "" {
			s.users[usrAddr].username = strings.ReplaceAll(string(msg), "\n", "")
			con.Write([]byte("Welcome " + string(msg) + "\n"))

			f := fmt.Sprintf("%s connected! (%d online) \n", s.users[con.RemoteAddr().String()].username, len(s.users))
			s.logger.Print(f)
			fmt.Print(f)
			isMessage = false
			if usr.connectedServer == nil {
				con.Write([]byte("Enter room name: "))
			}
		}
		if found && usr.connectedServer == nil {
			port, err := returnFirstAvailablePort(3000, 20)
			if err != nil {
				log.Fatalln(err)
			}
			server := NewServer(port)
			s.serverch <- server
			usr.connectedServer = server
			// TODO: transfer user connection to that server?????
			isMessage = false
		}

		if isMessage {
			s.msgch <- Message{
				from:    s.users[usrAddr],
				payload: msg,
			}
		}

	}
}

func (s *Server) handleMessages() {
	for msg := range s.msgch {
		formatted := fmt.Sprintf("> %s: %s \n", msg.from.username, string(msg.payload))

		for _, user := range s.users {
			if user.address != msg.from.address {
				user.conn.Write([]byte(formatted))
			}
		}

		fmt.Print(formatted)
		s.logger.Print(formatted)
	}
}
