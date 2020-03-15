package main

import (
	"context"
	"crypto/md5"
	"errors"
	"fmt"
	"strings"
	"time"
)

type User struct {
	tableName       string  `sql:"user" json:"-"` // nolint
	Id              int     `sql:",pk" json:"id"`
	Login           string  `sql:",unique" json:"login"`
	Password        *string `json:"password,omitempty"`
	IsGuest         bool    `sql:",use_zero" json:"is_guest"`
	Email           string  `json:"email"`
	CurrentPassword *string `json:"current_password,omitempty" sql:"-" pg:"-"`
	Token           *string `json:"token,omitempty" sql:"-" pg:"-"` // токен для сброса пароля
	Timestamps
}

const companyWhenspeak = 1

func (user *User) Get() (err error) {
	err = DB().Model(user).WherePK().Select()
	return
}

func (user *User) Create() (err error) {
	*user.Password = fmt.Sprintf("%x", md5.Sum([]byte(*user.Password)))
	err = DB().Insert(user)

	return
}

func (user *User) Update() (err error) {
	if user.Password != nil {
		*(user.CurrentPassword) = fmt.Sprintf("%x", md5.Sum([]byte(*(user.CurrentPassword))))
		currentUser := new(User)
		currentUser.Id = user.Id
		err = currentUser.Get()
		if err != nil {
			return
		}
		if *currentUser.Password != *user.CurrentPassword {
			return errors.New("incorrect current password")
		}
		*user.Password = fmt.Sprintf("%x", md5.Sum([]byte(*(user.Password))))
	}
	_, err = DB().Model(user).WherePK().UpdateNotZero()
	return
}

func (user *User) BeforeUpdate(ctx context.Context) (context.Context, error) {
	user.UpdatedAt = time.Now()
	return ctx, nil
}

func (user *User) GetUserByLogin(login string) (err error) {
	err = DB().Model(user).
		Where("LOWER(login) = ?", strings.ToLower(login)).
		Limit(1).
		Select()

	return
}

func (user *User) GetUserByEmail(email string) (err error) {
	err = DB().Model(user).
		Where("LOWER(email) = ?", strings.ToLower(email)).
		Select()

	return
}
