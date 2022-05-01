package kfk

type UStat struct {
	Offset int64
	UId    string
	ChatId string
}

func NewUser(usId string) *UStat {
	return &UStat{
		Offset: 0,
		UId:    usId,
		ChatId: "0",
	}
}
