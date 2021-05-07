package chats

import (
	"encoding/json"
	"fmt"
)

type Hub interface {
	Run()
}

type hub struct {
	chats     map[string]UserChats
	broadcast chan Message
	join      chan Subscription
	leave     chan Subscription
}

var h = &hub{
	chats:     make(map[string]UserChats),
	broadcast: make(chan Message),
	join:      make(chan Subscription),
	leave:     make(chan Subscription),
}

func GetHub() Hub {
	return h
}

func (h *hub) Run() {
	for {
		select {
		case s := <-h.join:
			fmt.Println("Joined:", s.from)

			// check if sender has no chats create one
			// or if has a chat with receiver just add connection to it
			if _, ok := h.chats[s.from]; !ok {
				h.chats[s.from] = make(UserChats)
			} else if _, ok := h.chats[s.from][s.to]; ok {
				h.chats[s.from][s.to].conns[s.conn] = true
				continue
			}

			// check if receiver has a chat with sender,
			// use that room otherwise create new room
			_, rcvHasChats := h.chats[s.to]
			_, rcvHasChatWithMe := h.chats[s.to][s.from]

			if rcvHasChats && rcvHasChatWithMe {
				room := h.chats[s.to][s.from]
				room.conns[s.conn] = true
				h.chats[s.from][s.to] = room
			} else {
				h.chats[s.from][s.to] = &Room{
					msgs: &[]Message{},
					conns: map[*connection]bool{
						s.conn: true,
					},
				}
			}

		case s := <-h.leave:
			connections := h.chats[s.from][s.to].conns
			if connections != nil {
				if _, ok := connections[s.conn]; ok {
					close(s.conn.send)
					delete(connections, s.conn)
				}
			}
			fmt.Println("Leaved:", s.from)

		case m := <-h.broadcast:
			msg, err := json.Marshal(m)
			if err != nil {
				fmt.Println(err)
				msg = []byte(m.Data)
			}
			fmt.Println("Message:", m.From, "->", m.To, "at", m.Time, m.Data)
			// add message to room messsages
			*h.chats[m.From][m.To].msgs = append(*h.chats[m.From][m.To].msgs, m)
			// broadcast messages to other connections
			connections := h.chats[m.From][m.To].conns
			for c := range connections {
				if c == m.Conn {
					continue
				}
				select {
				case c.send <- msg:
				default:
					close(c.send)
					delete(connections, c)
				}
			}
		}
	}
}
