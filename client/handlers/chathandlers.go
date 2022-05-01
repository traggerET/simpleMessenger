package handlers

import (
	"simpleMessenger/client/user"
)

func LeaveChat(u *user.User) error {
	return u.C.Call("Msgr.Leave", *u.Status, u.Status)
}

func CreateChat(u *user.User) error {
	return u.C.Call("Msgr.NewChat", *u.Status, u.Status)
}
func JoinChat(u *user.User, chat string) error {
	tmpStat := u.Status
	tmpStat.ChatId = chat
	err := u.C.Call("Msgr.Join", tmpStat, u.Status)
	return err
}
