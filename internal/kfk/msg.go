package kfk

import (
	"github.com/segmentio/kafka-go"
	"time"
)

type Msg struct {
	Time time.Time
	Chat string
	Msg  []byte
	Who  string
}

func MarshalKafka(m Msg) kafka.Message {
	return kafka.Message{
		Key:   []byte(m.Who),
		Value: m.Msg,
		Time:  m.Time,
	}
}

func UnmarshalKafka(m kafka.Message) Msg {
	return Msg{
		Time: m.Time,
		Msg:  m.Value,
		Who:  string(m.Key),
		Chat: m.Topic,
	}
}
