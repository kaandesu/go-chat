package main

import "net"

type User struct {
	conn            net.Conn
	connectedServer *Server
	mainServer      *Server
	address         string
	username        string
}

func NewUser(address, username string, conn net.Conn, mainServer *Server) *User {
	return &User{
		address:    address,
		username:   username,
		mainServer: mainServer,
		conn:       conn,
	}
}
