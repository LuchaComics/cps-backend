package datastore

import (
	"time"

	user_s "github.com/LuchaComics/cps-backend/app/user/datastore"
)

const (
	UserActiveState   = 1
	UserInactiveState = 2
)

type LoginResult struct {
	User                   *user_s.User `json:"user"`
	AccessToken            string       `json:"access_token"`
	AccessTokenExpiryTime  time.Time    `json:"access_token_expiry_time"`
	RefreshToken           string       `json:"refresh_token"`
	RefreshTokenExpiryTime time.Time    `json:"refresh_token_expiry_time"`
}
