package rds

import (
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"strings"
)

const (
	defaultStorageSocket = "localhost:6379"
	defaultChatIdLen     = 10
)

type keysStorage struct {
	clients *redis.Client
	chats   *redis.Client
}

var ks *keysStorage

func init() {
	ks = &keysStorage{}
}

func Start() error {

	ks.chats = redis.NewClient(&redis.Options{
		Addr: defaultStorageSocket,
		DB:   0,
	})
	ks.clients = redis.NewClient(&redis.Options{
		Addr: defaultStorageSocket,
		DB:   1,
	})

	if ks.chats == nil || ks.clients == nil {
		return errors.New("failed to establish keystorage Redis connection")
	}
	return nil
}

func IsUniqueId(uId string) bool {
	_, err := ks.clients.Get(uId).Int()
	if err != redis.Nil {
		return false
	}
	return true
}

func NewUser(uId string) error {

	if !IsUniqueId(uId) {
		return errors.New("this name is already in use by smone else... Choose another one")
	}
	err := ks.clients.Set(uId, "0", 0).Err()
	if err != nil {
		return err
	}

	err = ks.chats.Set("0", uId+",", 0).Err()
	if err != nil {
		return err
	}

	return nil
}

func Connect(uId, chatId string) error {
	currChat, err := ks.clients.Get(uId).Result()
	if err == redis.Nil {
		return errors.New("EnterRoom: user does not exists")
	}

	chatUsers, err := ks.chats.Get(chatId).Result()
	if err == redis.Nil {
		return errors.New("EnterRoom: room does not exists")
	}

	err = ks.clients.Set(uId, chatId, 0).Err()
	if err != nil {
		return err
	}

	chatUsers = chatUsers + fmt.Sprintf("%s,", uId)
	err = ks.chats.Set(chatId, chatUsers, 0).Err()
	if err != nil {
		return err
	}

	oldUsers, err := eraseFromList(uId, currChat)
	if err != nil {
		return err
	}

	if currChat != "0" && oldUsers == "" {
		return ks.chats.Del(currChat).Err()
	}

	return nil
}

func Disconnect(uId string) error {
	chatId, err := ks.clients.Get(uId).Result()
	if err == redis.Nil {
		return errors.New("LeaveRoom: user does not exists")
	}

	if chatId == "0" {
		_, err := eraseFromList(uId, "0")
		if err != nil {
			return err
		}
		return ks.clients.Del(uId).Err()
	}
	return Connect(uId, "0")

}

func eraseFromList(uId, chatId string) (string, error) {
	oldChatUsers, err := ks.chats.Get(chatId).Result()
	if err != nil {
		return "", err
	}
	oldChatUsers = strings.ReplaceAll(oldChatUsers, uId+",", "")
	return oldChatUsers, nil
}

func NewChat(firstClient string) (string, error) {
	key := make([]byte, defaultChatIdLen)
	rand.Read(key)
	chatId := string(key)

	err := ks.chats.Set(chatId, "", 0).Err()
	if err != nil {
		return "0", err
	}

	err = Connect(firstClient, chatId)
	if err != nil {
		return "0", err
	}
	return chatId, nil
}

func FlushKeys() {
	if ks.clients != nil {
		ks.clients.FlushAll()
	}
	if ks.chats != nil {
		ks.chats.FlushAll()
	}
}

func ShutDown() error {
	FlushKeys()

	err := ks.clients.Close()
	if err != nil {
		return err
	}

	err = ks.chats.Close()
	if err != nil {
		return err
	}

	ks.clients, ks.chats = nil, nil
	return nil
}
