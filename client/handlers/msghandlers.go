package handlers

import (
	"bufio"
	"fmt"
	"os"
	"simpleMessenger/client/user"
	msgr "simpleMessenger/pkg"
)

func PostMsg(u *user.User, msg string) error {
	var success bool
	err := u.C.Call("Msgr.PostMsg", msgr.Msg{
		Msg:  []byte(msg),
		Who:  u.Status.UId,
		Chat: u.Status.ChatId,
	}, &success)
	if err != nil {
		return err
	}
	if !success {
		return fmt.Errorf("failed to post message^ %s\n", msg)
	}
	return nil

}

func GetUnread(u *user.User) ([]msgr.Msg, error) {
	var resp msgr.UHist
	err := u.C.Call("Msgr.FetchUnread", *u.Status, &resp)
	if err != nil {
		return nil, err
	}
	u.Status = &resp.User
	return resp.History, nil
}

func GetHistory(u *user.User) ([]msgr.Msg, error) {
	var resp msgr.UHist
	err := u.C.Call("Msgr.FetchAll", *u.Status, &resp)
	if err != nil {
		return nil, err
	}
	u.Status = &resp.User
	return resp.History, nil
}

func SaveHist(hist []msgr.Msg, f string, append bool) error {
	flags := os.O_RDWR | os.O_CREATE
	if append {
		flags |= os.O_APPEND
	}

	file, err := os.OpenFile(f, flags, 0666)
	if err != nil {
		return err
	}
	defer func() {
		file.Close()
	}()

	w := bufio.NewWriter(file)
	for _, msg := range hist {
		w.Write(unmarshalMsg(msg))
		w.Write(msg.Msg)
	}
	w.Flush()

	return nil
}

func unmarshalMsg(m msgr.Msg) []byte {
	return []byte(fmt.Sprintf("%s\t%s\t", m.Who, m.Time.String()))
}
