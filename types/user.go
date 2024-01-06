package types

import "strings"

type User struct {
	Name       string
	Email      string
	Validation map[string]string
}

func NewUser(name, email string) *User {
	return &User{
		Name:       name,
		Email:      email,
		Validation: make(map[string]string, 0),
	}
}

func (u *User) IsValid() bool {
	return u.Validation["Email"] == "" && u.Validation["Name"] == ""
}

func (u *User) Validate(storage *map[string]User) {
	if _, exist := (*storage)[u.Name]; exist == true {
		u.Validation["Name"] = "Username is already taken."
	}

	if strings.Contains(u.Email, "@") == false {
		u.Validation["Email"] = "Not valid email"
	}
}
