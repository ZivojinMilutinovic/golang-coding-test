package store_api

import (
	"time"
)

type ValueType int

const (
	StringType = iota
	ListType
)

type Value struct {
	Type      ValueType
	Data      interface{}
	ExpiresAt time.Time // zero value if no expiry
}

type CommandType int

const (
	CmdGet CommandType = iota
	CmdSet
	CmdUpdate
	CmdRemove
	CmdPush
	CmdPop
)

type Command struct {
	Type    CommandType
	Key     string
	Value   interface{}
	TTL     time.Duration //used for Set
	ReplyCh chan interface{}
}

type Store struct {
	data map[string]*Value
	cmds chan Command
}

func NewStore() *Store {
	s := &Store{
		data: make(map[string]*Value),
		cmds: make(chan Command, 1000),
	}
	go s.run()
	return s
}

func (s *Store) run() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case cmd := <-s.cmds:
			s.handleCommand(cmd)
		case <-ticker.C:
			s.cleanupExpired()
		}
	}

}

func (s *Store) handleCommand(cmd Command) {
	switch cmd.Type {
	case CmdGet:
		s.handleGet(cmd)
	case CmdSet:
		s.handleSet(cmd)
	case CmdUpdate:
		s.handleUpdate(cmd)
	case CmdRemove:
		s.handleRemove(cmd)
	case CmdPush:
		s.handlePush(cmd)
	case CmdPop:
		s.handlePop(cmd)
	}
}

func (s *Store) handleGet(cmd Command) {
	v, ok := s.data[cmd.Key]
	if ok && !s.isExpired(v) {
		cmd.ReplyCh <- v.Data
	} else {
		cmd.ReplyCh <- nil
	}
}

func (s *Store) handleSet(cmd Command) {
	expiresAt := time.Time{}
	if cmd.TTL > 0 {
		expiresAt = time.Now().Add(cmd.TTL)
	}
	s.data[cmd.Key] = &Value{
		Type:      detectType(cmd.Value),
		Data:      cmd.Value,
		ExpiresAt: expiresAt,
	}
	cmd.ReplyCh <- true
}

func (s *Store) handleUpdate(cmd Command) {
	v, ok := s.data[cmd.Key]
	if ok && !s.isExpired(v) {
		v.Data = cmd.Value
		cmd.ReplyCh <- true
	} else {
		cmd.ReplyCh <- false
	}
}

func (s *Store) handleRemove(cmd Command) {
	delete(s.data, cmd.Key)
	cmd.ReplyCh <- true
}

func (s *Store) handlePush(cmd Command) {
	v, ok := s.data[cmd.Key]
	if !ok || v.Type != ListType {
		v = &Value{Type: ListType, Data: []string{}, ExpiresAt: time.Time{}}
		s.data[cmd.Key] = v
	}
	list := v.Data.([]string)
	item := cmd.Value.(string)
	list = append(list, item)
	v.Data = list
	cmd.ReplyCh <- true
}

func (s *Store) handlePop(cmd Command) {
	v, ok := s.data[cmd.Key]
	if ok && v.Type == ListType {
		list := v.Data.([]string)
		if len(list) == 0 {
			cmd.ReplyCh <- nil
		} else {
			item := list[0]
			v.Data = list[1:]
			cmd.ReplyCh <- item
		}
	} else {
		cmd.ReplyCh <- nil
	}
}

func (s *Store) isExpired(v *Value) bool {
	if v.ExpiresAt.IsZero() {
		return false
	}
	return time.Now().After(v.ExpiresAt) // is curren time after specific time
}

func (s *Store) cleanupExpired() {
	now := time.Now()
	for k, v := range s.data {
		if !v.ExpiresAt.IsZero() && now.After(v.ExpiresAt) {
			delete(s.data, k)
		}
	}
}

func detectType(val interface{}) ValueType {
	switch val.(type) {
	case []string:
		return ListType
	default:
		return StringType
	}
}

func (s *Store) Get(key string) interface{} {
	reply := make(chan interface{})
	s.cmds <- Command{Type: CmdGet, Key: key, ReplyCh: reply}
	return <-reply
}

func (s *Store) Set(key string, value interface{}, ttl time.Duration) {
	reply := make(chan interface{})
	s.cmds <- Command{Type: CmdSet, Key: key, Value: value, TTL: ttl, ReplyCh: reply}
	<-reply
}

func (s *Store) Update(key string, value interface{}) bool {
	reply := make(chan interface{})
	s.cmds <- Command{Type: CmdUpdate, Key: key, Value: value, ReplyCh: reply}
	return (<-reply).(bool)
}

func (s *Store) Remove(key string) {
	reply := make(chan interface{})
	s.cmds <- Command{Type: CmdRemove, Key: key, ReplyCh: reply}
	<-reply
}

func (s *Store) Push(key string, value string) {
	reply := make(chan interface{})
	s.cmds <- Command{Type: CmdPush, Key: key, Value: value, ReplyCh: reply}
	<-reply
}

func (s *Store) Pop(key string) string {
	reply := make(chan interface{})
	s.cmds <- Command{Type: CmdPop, Key: key, ReplyCh: reply}
	res := <-reply
	if res == nil {
		return ""
	}
	return res.(string)
}
