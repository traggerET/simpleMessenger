package main

import msgr "simpleMessenger/pkg"

type Msgr struct{}

func (*Msgr) NewUser(uId string, resp *msgr.UStat) error {
	resp, err := msgr.NewUser(uId)
	return err
}

func (*Msgr) NewChat(creator msgr.UStat, resp *msgr.UStat) error {
	resp, err := msgr.NewChat(creator)
	return err
}

func (*Msgr) Leave(u msgr.UStat, resp *msgr.UStat) error {
	resp, err := msgr.Leave(u)
	return err
}

func (*Msgr) Join(u msgr.UStat, resp *msgr.UStat) error {
	resp, err := msgr.Join(u, u.ChatId)
	if err != nil {
		resp.ChatId = "0"
	}
	return err
}

func (*Msgr) ViewHistory(u msgr.UStat, resp *msgr.UHist) error {
	hist, err := msgr.ViewHistory(u)
	resp.User, resp.History = u, hist
	return err
}

func (*Msgr) FetchUnread(u msgr.UStat, resp *msgr.UHist) error {
	var err error
	resp.User, resp.History, err = msgr.ViewUnread(u)
	return err
}

func (*Msgr) PostMsg(msg msgr.Msg, resp *bool) error {
	err := msgr.SendMsg(msg)
	if err == nil {
		*resp = true
	}
	return err
}
