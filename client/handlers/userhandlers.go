package handlers

import (
	"fmt"
	"net/rpc/jsonrpc"
	"simpleMessenger/client/user"
	msgr "simpleMessenger/pkg"
)

func Login(usr, ip string) (*user.User, error) {
	c, err := jsonrpc.Dial("tcp", ip+":38120")
	if err != nil {
		return nil, fmt.Errorf("\t[ERROR] Connection error\n\n")
	}

	response := &msgr.UStat{}
	err = c.Call("Msgr.NewUser", usr, &response)
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}

	return &user.User{OnIP: ip, Status: response, C: c}, nil
}
