package user

import (
	"net/rpc"
	msgr "simpleMessenger/pkg"
)

type User struct {
	OnIP   string
	Status *msgr.UStat
	C      *rpc.Client
}
