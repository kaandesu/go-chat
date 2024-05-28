package main

import "net"

type User struct {
	conn            net.Conn
	connectedServer *Server
	address         string
	username        string
}

func NewUser(address, username string, conn net.Conn) *User {
	return &User{
		address:  address,
		username: username,
		conn:     conn,
	}
}
