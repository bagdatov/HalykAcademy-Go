package server

import (
	"errors"
	"sync"
)

type userBase struct {
	mu    sync.Mutex
	store map[string]user
}

type user struct {
	Username string
	Email    string
}

var (
	errNotFound     = errors.New("user not found")
	errAlreadyExist = errors.New("user already exist")
)

var base = &userBase{
	store: map[string]user{
		"alibi": {Username: "alibi", Email: "test@mail.ru"},
		"jojo":  {Username: "jojo", Email: "jojo@gmail.com"},
	},
}

func (ub *userBase) findUser(username string) (user, error) {
	ub.mu.Lock()
	defer ub.mu.Unlock()

	u, ok := ub.store[username]
	if !ok {
		return u, errNotFound
	}

	return u, nil
}

func (ub *userBase) showAll() []string {
	ub.mu.Lock()
	defer ub.mu.Unlock()

	users := make([]string, 0, len(ub.store))

	for k := range ub.store {
		users = append(users, k)
	}

	return users
}

func (ub *userBase) saveUser(u user) error {
	ub.mu.Lock()
	defer ub.mu.Unlock()

	_, ok := ub.store[u.Username]
	if ok {
		return errAlreadyExist
	}

	ub.store[u.Username] = u
	return nil
}
