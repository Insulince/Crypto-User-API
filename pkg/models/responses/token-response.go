package responses

import "crypto-users/pkg/models"

type Token struct {
	Message string `json:"message"`
	Token models.Token `json:"token"`
	QS    string       `json:"qs"`
}
