package msgr

import (
	"simpleMessenger/internal/kfk"
	"simpleMessenger/internal/rds"
)

func Start() error {
	err := rds.Start()
	if err != nil {
		return err
	}

	return kfk.Start()
}

func ShutDown() error {
	err := kfk.ShutDown()
	if err != nil {
		return err
	}

	return rds.ShutDown()
}

type UStat = kfk.UStat

func NewUser(uId string) (*UStat, error) {
	err := rds.NewUser(uId)
	if err != nil {
		return nil, err
	}
	return kfk.NewUser(uId), nil
}

func NewChat(firstU UStat) (*UStat, error) {
	chat, err := rds.NewChat(firstU.UId)
	if err != nil {
		return nil, err
	}

	err = kfk.Connect(&firstU, chat)
	if err != nil {
		return nil, err
	}

	return &firstU, nil
}

func Join(u UStat, chId string) (*UStat, error) {
	err := rds.Connect(u.UId, chId)
	if err != nil {
		return nil, err
	}

	err = kfk.Connect(&u, chId)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func Leave(u UStat) (*UStat, error) {
	err := rds.Disconnect(u.UId)
	if err != nil {
		return nil, err
	}
	kfk.Disconnect(&u)
	return &u, nil
}

type Msg = kfk.Msg

type UHist struct {
	User    UStat
	History []Msg
}

func SendMsg(msg Msg) error {
	return kfk.PostMsg(msg)
}

func ViewUnread(user UStat) (UStat, []Msg, error) {
	h, e := kfk.FetchUnread(&user)
	return user, h, e
}

func ViewHistory(user UStat) ([]Msg, error) {
	return kfk.FetchAll(&user)
}
