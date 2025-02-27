package domain

import "time"

type User struct {
	id        int
	username  string
	email     string
	password  string
	countryID int
	createdAt time.Time
	updatedAt time.Time
}

type UserData struct {
	ID        int
	Username  string
	Email     string
	Password  string
	CountryID int
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewUser(data UserData) (User, error) {
	return User{
		id:        data.ID,
		username:  data.Username,
		email:     data.Email,
		password:  data.Password,
		countryID: data.CountryID,
		createdAt: data.CreatedAt,
		updatedAt: data.UpdatedAt,
	}, nil
}

func (u User) ID() int {
	return u.id
}

func (u User) Username() string {
	return u.username
}

func (u User) Email() string {
	return u.email
}

func (u User) Password() string {
	return u.password
}

func (u User) CountryID() int {
	return u.countryID
}

func (u User) CreatedAt() time.Time {
	return u.createdAt
}

func (u User) UpdatedAt() time.Time {
	return u.updatedAt
}
