package models

import (
	"time"
)

type House struct {
	Id                int        `json:"id"`
	Address           string     `json:"address"`
	BuildYear         int        `json:"buildYear"`
	Developer         *string    `json:"developer"`
	DateOfCreation    time.Time  `json:"dareOfCreation"`
	Lastflataddeddate *time.Time `json:"lastFlatAddDate"`
}

type Flat struct {
	Id               int `json:"id"`
	HouseId          int `json:"houseId"`
	Price            int `json:"price"`
	RoomsNumber      int `json:"roomsNumber"`
	FlatNumber       int `json:"flatNumber"`
	ModerationStatus int `json:"moderationSatus"`
}

type User struct {
	Id           int    `json:"id"`
	Login        string `json:"login"`
	PasswordHash string `json:"passwordHash"`
	Password     string `json:"password"`
	UserRole     string `json:"userRole"`
}

type Session struct {
	SessionId string    `json:"sessionId"`
	UserId    int       `json:"userId"`
	UserRole  string    `json:"userRole"`
	ExpireAt  time.Time `json:"expireAt"`
}

type Tokens struct {
	NewAccessToken  bool
	AccessToken     string
	NewRefreshToken bool
	RefreshToken    string
}
