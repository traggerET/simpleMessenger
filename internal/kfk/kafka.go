package kfk

import (
	"context"
	"errors"
	"fmt"
	"github.com/segmentio/kafka-go"
	"io"
)

const (
	kafkaConnection  = "tcp"
	kafkaPort        = ":9092"
	kafkaDefaultChat = "hall"
)

type Chats map[string]*kafka.Conn

var chats Chats

func Start() error {
	if chats != nil {
		return errors.New("kafka initialized")
	}

	dr, err := kafka.DialLeader(context.Background(), kafkaConnection, kafkaPort, kafkaDefaultChat, 0)
	if err != nil {
		return err
	}

	chats = map[string]*kafka.Conn{
		kafkaDefaultChat: dr,
	}
	return nil
}

func Connect(stat *UStat, chatId string) error {
	if _, ok := chats[chatId]; !ok {
		rm, err := kafka.DialLeader(context.Background(), kafkaConnection, kafkaPort, chatId, 0)
		if err != nil {
			return err
		}
		chats[chatId] = rm
	}
	offset, err := chats[chatId].ReadLastOffset()
	if err != nil {
		return err
	}
	stat.Offset = offset
	stat.ChatId = chatId
	return nil
}

func PostMsg(msg Msg) error {
	if _, ok := chats[msg.Chat]; !ok {
		return fmt.Errorf("can not post msg:room %s does not exist", msg.Chat)
	}
	c, _ := chats[msg.Chat]

	marshalled := MarshalKafka(msg)
	_, err := c.WriteMessages(marshalled)

	return err
}

func Disconnect(usr *UStat) {
	usr.ChatId = kafkaDefaultChat
	usr.Offset = 0
}

func FetchAll(usr *UStat) ([]Msg, error) {
	usr.Offset = 0
	return FetchUnread(usr)
}

func FetchUnread(usr *UStat) ([]Msg, error) {
	conn, err := kafka.DialLeader(context.Background(), kafkaConnection, kafkaPort, usr.ChatId, 0)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	_, err = conn.Seek(usr.Offset, kafka.SeekStart)
	if err != nil {
		return nil, err
	}

	last, _ := chats[usr.ChatId].ReadLastOffset()

	msgs, err, newLast := readQueue(conn, last)
	usr.Offset = newLast
	return msgs, err
}

func readQueue(chat *kafka.Conn, last int64) ([]Msg, error, int64) {
	msgs := make([]Msg, 0)
	for {
		m, err := chat.ReadMessage(4096)
		curr := m.Offset

		if err == io.EOF {
			break
		}
		if err != nil {
			return msgs, err, m.Offset
		}

		msgs = append(msgs, UnmarshalKafka(m))
		if curr == last-1 {
			break
		}
	}
	return msgs, nil, last
}

func ShutDownOne(chatId string) error {
	_, ok := chats[chatId]
	if !ok {
		return fmt.Errorf("can not shut down room %s: no such room", chatId)
	}

	err := chats[chatId].DeleteTopics(chatId)
	if err != nil {
		return err
	}

	err = chats[chatId].Close()
	if err != nil {
		return err
	}

	return nil
}

func ShutDownAll() error {
	if chats == nil {
		return errors.New("kafka not initialised")
	}

	for i, _ := range chats {

		err := ShutDownOne(i)
		if err != nil {
			return err
		}

	}
	chats = nil
	return nil
}

func ShutDown() error {
	return ShutDownAll()
}
