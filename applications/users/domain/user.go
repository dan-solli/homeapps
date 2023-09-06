package user

import "github.com/google/uuid"

type User struct {
	id   uuid.UUID
	Name string
}

type UserList struct {
	users []User
}

func (s *UserList) Add(name string) {
	s.users = append(s.users, User{
		id:   uuid.New(),
		Name: name,
	})
}

func (s *UserList) Users() []User {
	return s.users
}

func (s *UserList) Search(search string) []User {
	var users []User
	for _, user := range s.users {
		if user.Name == search {
			users = append(users, user)
		}
	}
	return users
}
