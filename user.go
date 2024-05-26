package main

type User struct {
	address  string
	username string
}

func NewUser(address, username string) *User {
	return &User{
		address:  address,
		username: username,
	}
}
