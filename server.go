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
	name         string
	address      string
}

func NewServer(address, name string) *Server {
	return &Server{
		address:      address,
		name:         name,
		msgch:        make(chan Message),
		quitch:       make(chan struct{}),
		users:        make(map[string]*User),
		childServers: make(map[string]*Server),
		serverch:     make(chan *Server),
		logger:       NewLogger("./chat.log"),
	}
}

func availablePort(portNo int, tryUntil int) (string, error) {
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
	s.logAndPrint(fmt.Sprintf("Starting server on [%s] \n", s.address))
	if err != nil {
		return err
	}

	defer ln.Close()
	s.listener = ln

	go s.acceptLoop()
	go s.distributeMessages()

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
			s.logAndPrint(formatted)
			delete(s.users, con.RemoteAddr().String())
			break
		}

		msg := strings.ReplaceAll(string(buf[:n]), "\n", "")
		fmt.Print(msg)
		usrAddr := con.RemoteAddr().String()

		usr, found := s.users[usrAddr]

		isMessage := true
		if found && usr.username == "" {
			s.users[usrAddr].username = msg
			con.Write([]byte(fmt.Sprintf("Welcome %s! \n", msg)))

			f := fmt.Sprintf("%s connected! (%d online) \n", s.users[con.RemoteAddr().String()].username, len(s.users))
			s.logAndPrint(f)
			isMessage = false
			if usr.connectedServer == nil {
				con.Write([]byte("Enter room name: "))
				continue
			}
		}

		if found && usr.connectedServer == nil {
			childServer, serverfound := s.childServers[msg]
			if serverfound {
				usr.connectedServer = childServer
				childServer.users[usr.address] = usr
				go childServer.distributeMessages()
			} else {

				port, err := availablePort(3000, 20)
				if err != nil {
					s.logger.Fatalln(err)
					log.Fatalln(err)
				}
				user := s.users[usrAddr]
				server := NewServer(port, msg)
				server.users[user.address] = user
				user.connectedServer = server
				s.childServers[msg] = server
				s.serverch <- server
			}
			con.Write([]byte("\n(type /help for all the commands)\n\n"))
			isMessage = false
		}

		if found && usr.connectedServer != nil && isMessage {
			handleMessages(usr, buf[:n])
		}

	}
}

func (s *Server) logAndPrint(text string) {
	fmt.Print(text)
	s.logger.Print(text)
}

func (s *Server) broadcastMessage(msg Message) {
	formatted := fmt.Sprintf("[%s:%s] >> %s: %s \n", msg.from.connectedServer.name, msg.from.connectedServer.address, msg.from.username, string(msg.payload))
	for _, user := range s.users {
		if user.address != msg.from.address {
			user.conn.Write([]byte(formatted))
		}
	}
	s.logAndPrint(formatted)
}

func handleMessages(from *User, payload []byte) {
	message := strings.ReplaceAll(string(payload), "\n", "")
	// TODO: instead of checking the whole message check its a prefixes
	switch message {
	case "/exit":
		// TODO: when connection is closed delete the user address from the main and room server
		from.conn.Close()
	case "/join":
		from.conn.Write([]byte("Command not available. \n\n"))
	case "/rename":
		from.conn.Write([]byte("Command not available. \n\n"))
	case "/list":
		from.conn.Write([]byte("Command not available. \n\n"))
	case "/help":
		// TODO: i hate this long line below, will chanage it
		from.conn.Write([]byte("\n/exit : disconnect from the server \n/list : list all available rooms \n/join <room_name> : join another room \n/rename <name> : change your username (will notify all the users in your current room) \n/help lists all commands\n\n"))
	default:
		if strings.HasPrefix(message, "/") {
			from.conn.Write([]byte("Invalid commmand: " + message + "\n\n"))
		} else {
			from.connectedServer.msgch <- Message{
				from:    from.connectedServer.users[from.address],
				payload: payload,
			}
		}
	}
}

func (s *Server) distributeMessages() {
	for msg := range s.msgch {
		msg.from.connectedServer.broadcastMessage(msg)
	}
}
