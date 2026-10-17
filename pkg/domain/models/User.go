package models

import (
	"fmt"
	"time"
)

type User struct {
	Id        int                    `json:"id"`
	User_id   string                 `json:"user_id"`
	Name      string                 `json:"name"`
	Email     string                 `json:"email"`
	Address   map[string]interface{} `json:"address"`
	Password  string                 `json:"password"`
	Create_at time.Time              `json:"created_at"`
	Delete_at time.Time              `json:"deleted_at"`
}

type UserDetails struct {
	Name     string                 `json:"name"`
	Email    string                 `json:"email"`
	Address  map[string]interface{} `json:"address"`
	Password string                 `json:"password"`
}

type UserUpdate struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UrlShorter struct {
	User_id     string    `json:"user_id"`
	LongestUrl  string    `json:"longest"`
	ShortestUrl string    `json:"short"`
	Created_at  time.Time `json:"created_at"`
	Deleted_at  time.Time `json:"deleted_at"`
}

func (obj *UserDetails) Validate() error {
	if obj == nil {
		return nil
	}
	if obj.Name == "" {
		return fmt.Errorf("it should have name")
	} else if obj.Email == "" {
		return fmt.Errorf("Mail must be given")
	} else {
		return fmt.Errorf("failed was missing, check with json")
	}
}
