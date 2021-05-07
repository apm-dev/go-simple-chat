package chats

import (
	"errors"
	"sync"
)

type UserChats map[string]*Room

type Room struct {
	msgs  *[]Message
	conns map[*connection]bool
}

type Subscription struct {
	conn *connection
	from string
	to   string
}

type Message struct {
	From string      `json:"sender"`
	To   string      `json:"-"`
	Data string      `json:"data"`
	Time int64       `json:"time"`
	Conn *connection `json:"-"`
}

var (
	ErrNoChatFound = errors.New("there is no chat")

	clock  sync.RWMutex
	uclock sync.RWMutex
)

func GetChatsOfUser(uname string) ([]string, error) {
	clock.RLock()
	uc, ok := h.chats[uname]
	clock.RUnlock()

	if !ok || len(uc) == 0 {
		return nil, ErrNoChatFound
	}

	uclock.RLock()
	defer uclock.RUnlock()
	cs := make([]string, 0, len(uc))
	if len(uc) == 0 {
		return nil, ErrNoChatFound
	}
	for c := range uc {
		cs = append(cs, c)
	}
	return cs, nil
}

func GetUserChatMessages(sender string, receiver string) ([]Message, error) {
	clock.RLock()
	uc, ok := h.chats[sender]
	clock.RUnlock()

	if !ok || len(uc) == 0 {
		return nil, ErrNoChatFound
	}

	uclock.RLock()
	room, ok := uc[receiver]
	uclock.RUnlock()

	if !ok {
		return nil, ErrNoChatFound
	}

	return *room.msgs, nil
}
