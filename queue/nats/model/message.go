// Copyright (c) 2016 - Sarjono Mukti Aji <me@simukti.net>
// Unless otherwise noted, this source code license is MIT-License

package model

var (
	UserEmailSubject = "user_email"
)

type UserEmail struct {
	UserId uint32 `json:"user_id"`
	Type   string `json:"type"`
}
